package common

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

var (
	RedisPool *redis.Pool
)

func InitRedisPool(network string, host string, password string, database int) (pool *redis.Pool) {
	pool = &redis.Pool{
		MaxIdle:     1000,
		MaxActive:   10000,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(network, host)
			if err != nil {
				return nil, err
			}
			if password != "" {

				if _, err := conn.Do("AUTH", password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			if _, err := conn.Do("SELECT", database); err != nil {
				panic(err)
				return nil, err
			}
			return conn, nil
		},
	}
	RedisPool = pool
	return
}

func RedisSet(ctx *gin.Context, redisConn redis.Conn, key string, value interface{}) (code int, err error) {
	parentCtx, ok := ctx.Get("ParentSpanContext")
	var redisSpan opentracing.Span
	if ok {
		if tracer := opentracing.GlobalTracer(); tracer != nil {
			redisSpan = tracer.StartSpan("RedisSpanSet", opentracing.ChildOf(parentCtx.(opentracing.SpanContext)))
			defer redisSpan.Finish()

			redisSpan.SetTag("redis_conn", redisConn)
			redisSpan.SetTag("key", key)
			redisSpan.SetTag("value", value)
		}
	}
	var valueJson []byte

	if valueJson, err = json.Marshal(value); err != nil {
		redisSpan.SetTag("error", true)
		return REDIS_SET_ERROR, err
	}

	if _, err := redisConn.Do("Set", key, valueJson); err != nil {
		redisSpan.SetTag("error", true)
		return REDIS_SET_ERROR, err
	}
	return
}

func RedisGet(ctx *gin.Context, redisConn redis.Conn, key string) (value string, err error) {
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		pctx := parent.Context()
		if tracer := opentracing.GlobalTracer(); tracer != nil {
			redisSpan := tracer.StartSpan("RedisSpanGet", opentracing.ChildOf(pctx))
			defer redisSpan.Finish()
			redisSpan.SetTag("redis_conn", redisConn)
			redisSpan.SetTag("key", key)
		}
	}
	if value, err = redis.String(redisConn.Do("Get", key)); err != nil {
		value = "0"
	}
	return
}

func RedisSetCommon(redisConn redis.Conn, key string, value interface{}) (code int, err error) {
	var valueJson []byte

	if valueJson, err = json.Marshal(value); err != nil {
		return REDIS_SET_ERROR, err
	}

	if _, err := redisConn.Do("Set", key, valueJson); err != nil {
		return REDIS_SET_ERROR, err
	}
	return
}

func RedisGetCommon(redisConn redis.Conn, key string) (value string, err error) {
	if value, err = redis.String(redisConn.Do("Get", key)); err != nil {
		value = "0"
	}
	return
}
