package mq

import (
    "errors"
    "fmt"
    "github.com/streadway/amqp"
    "time"
)

var MqConn *amqp.Connection

var MqChan *amqp.Channel

var MqQueue map[string]amqp.Queue

var MqUriStore string

func MQInit(rmqUri string)(*amqp.Connection) {
    var (
        err error
        conn *amqp.Connection
        trys = 10
    )
    MqUriStore = rmqUri
    for i:=1;i<trys;i++ {
        conn, err = amqp.Dial(rmqUri)
        if err != nil {
            fmt.Println(err)
        } else {
            break
        }
        time.Sleep(time.Second * 5)
    }
    MqConn = conn
    return conn
}

func MQChannelRefresh() {
    var err error
    MqChan, err = MqConn.Channel()
    if err != nil {
        panic(err)
    }
    MqQueue = make(map[string]amqp.Queue)
}

func GetMqChannel() *amqp.Channel {
    if MqConn.IsClosed() {
        MQInit(MqUriStore)
    }
    //MQChannelRefresh()
    MqChan1, err := MqConn.Channel()
    if err != nil {
        panic(err)
    }
    return MqChan1
}

func QueueAdd(key string, queue amqp.Queue) {
    MqQueue[key] = queue
}

func QueueGet(key string) (queue amqp.Queue, err error) {
    channel, result := MqQueue[key]
    if !result {
        err = errors.New("channel not found !")
        return
    }
    return channel, nil
}
