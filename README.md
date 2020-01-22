# gin_common
gin common for the project which using the  gin framework, like go_gin_sample

# 说明

使用 `gin_common` 来快速构建基于 `gin` 框架的程序。

也可以使用现有的例子 ： go_gin_sample

https://github.com/tsbxmw/go_gin_sample


## 内容

封装了一揽子常用的功能。

> common :

- app 构建 ：基于 github.com/urfave/cli 
- config ：基于 github.com/spf13/viper
- orm ：基于 github.com/jinzhu/gorm
- redis ：基于 github.com/garyburd/redigo/redis
- opentracing ： 基于 github.com/opentracing/opentracing-go
- rabbitmq ： 基于 github.com/streadway/amqp
- logrus ：基于 github.com/sirupsen/logrus
- consul ： 基于 github.com/hashicorp/consul

> middleware :

- apidoc ：基于 github.com/betacraft/yaag
- exception ： 修改于 go-gin 的 recovery
- logger-file : 基于 github.com/lestrrat/go-file-rotatelogs
- response ： common response
- jaeger : github.com/uber/jaeger-client-go/

> 重写 handler 方法

- NoRouter
- NoMethod

