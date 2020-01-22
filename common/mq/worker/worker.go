package worker

import (
    "github.com/tsbxmw/datasource/common"
    "github.com/tsbxmw/datasource/common/mq"
    "fmt"
    "github.com/smallnest/rpcx/log"
    "github.com/streadway/amqp"
    "sync"
    "time"
)

type Receiver interface {
    QueueName() string     // 获取接收者需要监听的队列
    RouterKey() string     // 这个队列绑定的路由
    OnError(error)         // 处理遇到的错误，当RabbitMQ对象发生了错误，他需要告诉接收者处理错误
    OnReceive([]byte) bool // 处理收到的消息, 这里需要告知RabbitMQ对象消息是否处理成功
}

type RabbitMQ struct {
    wg           sync.WaitGroup
    conn         *amqp.Connection
    channel      *amqp.Channel
    exchangeName string // exchange的名称
    exchangeType string // exchange的类型
    receivers    []Receiver
}

// New 创建一个新的操作RabbitMQ的对象
func New() *RabbitMQ {
    // 这里可以根据自己的需要去定义
    ExchangeType := "direct"
    ExchangeName := "data_1"
    return &RabbitMQ{
        exchangeName: ExchangeName,
        exchangeType: ExchangeType,
    }
}

func (rmq *RabbitMQ) prepareExchange() error {
    // 申明Exchange
    err := rmq.channel.ExchangeDeclare(
        rmq.exchangeName, // exchange
        rmq.exchangeType, // type
        true,             // durable
        false,            // autoDelete
        false,            // internal
        false,            // noWait
        nil,              // args
    )

    if nil != err {
        return err
    }
    return nil
}

func (rmq *RabbitMQ) Refresh() bool {
    //mq.MQInit(mq.MqUriStore)
    var err error
    rmq.conn, err = amqp.Dial(mq.MqUriStore)
    if err != nil {
        panic(err)
    }
    return true
}

// run 开始获取连接并初始化相关操作
func (rmq *RabbitMQ) run() {
    var err error

    if !rmq.Refresh() {
        log.Errorf("rabbit刷新连接失败，将要重连: %s", mq.MqUriStore)
        return
    }

    // 获取新的channel对象
    rmq.channel, err = rmq.conn.Channel()
    if err != nil {
        panic(err)
    }
    defer rmq.channel.Close()

    // 初始化Exchange
    rmq.prepareExchange()

    for _, receiver := range rmq.receivers {
        rmq.wg.Add(1)
        go rmq.listen(receiver) // 每个接收者单独启动一个goroutine用来初始化queue并接收消息
    }

    rmq.wg.Wait()

    log.Errorf("所有处理queue的任务都意外退出了")

    // 理论上rmq.run()在程序的执行过程中是不会结束的
    // 一旦结束就说明所有的接收者都退出了，那么意味着程序与RabbitMQ的连接断开
    // 那么则需要重新连接，这里尝试销毁当前连接

}

// Start 启动RabbitMQ的客户端
func (rmq *RabbitMQ) Start() {
    for {
        rmq.run()

        // 一旦连接断开，那么需要隔一段时间去重连
        // 这里最好有一个时间间隔
        time.Sleep(3 * time.Second)
    }
}

// RegisterReceiver 注册一个用于接收指定队列指定路由的数据接收者
func (rmq *RabbitMQ) RegisterReceiver(receiver Receiver) {
    rmq.receivers = append(rmq.receivers, receiver)
}

// Listen 监听指定路由发来的消息
// 这里需要针对每一个接收者启动一个goroutine来执行listen
// 该方法负责从每一个接收者监听的队列中获取数据，并负责重试
func (rmq *RabbitMQ) listen(receiver Receiver) {
    common.LogrusLogger.Info("start receiver!")
    defer rmq.wg.Done()

    // 这里获取每个接收者需要监听的队列和路由
    queueName := receiver.QueueName()
    routerKey := receiver.RouterKey()

    // 申明Queue
    _, err := rmq.channel.QueueDeclare(
        queueName, // name
        true,      // durable
        false,     // delete when usused
        false,     // exclusive(排他性队列)
        false,     // no-wait
        nil,       // arguments
    )
    if nil != err {
        // 当队列初始化失败的时候，需要告诉这个接收者相应的错误
        receiver.OnError(fmt.Errorf("初始化队列 %s 失败: %s", queueName, err.Error()))
    }

    // 将Queue绑定到Exchange上去
    err = rmq.channel.QueueBind(
        queueName,        // queue name
        routerKey,        // routing key
        rmq.exchangeName, // exchange
        false,            // no-wait
        nil,
    )
    if nil != err {
        receiver.OnError(fmt.Errorf("绑定队列 [%s - %s] 到交换机失败: %s", queueName, routerKey, err.Error()))
    }

    // 获取消费通道
    rmq.channel.Qos(1, 0, true) // 确保RabbitMQ会一个一个发消息
    msgs, err := rmq.channel.Consume(
        queueName, // queue
        "",        // consumer
        false,     // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
    if nil != err {
        receiver.OnError(fmt.Errorf("获取队列 %s 的消费通道失败: %s", queueName, err.Error()))
    }

    // 使用callback消费数据
    for msg := range msgs {
        // 当接收者消息处理失败的时候，
        // 比如网络问题导致的数据库连接失败，redis连接失败等等这种
        // 通过重试可以成功的操作，那么这个时候是需要重试的
        // 直到数据处理成功后再返回，然后才会回复RabbitMQ ack
        retryTimes := 0
        for !receiver.OnReceive(msg.Body) {
            log.Warnf("receiver 数据处理失败，将要重试")
            time.Sleep(1 * time.Second)
            if retryTimes > 10{
                panic(common.NewMySqlCreateError())
            }
            retryTimes ++
        }
        common.LogrusLogger.Info("Receive msg ok!")
        // 确认收到本条消息, multiple必须为false
        msg.Ack(false)
    }
}
