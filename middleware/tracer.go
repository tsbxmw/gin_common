package middleware

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/opentracing/opentracing-go"
    "github.com/uber/jaeger-client-go"
    jaegerConfig "github.com/uber/jaeger-client-go/config"
    _ "github.com/uber/jaeger-client-go/zipkin"
    "io"
)

var TracerCommon *opentracing.Tracer
var TracerCloser *io.Closer

func TracerInit(e *gin.Engine, jaegerAddr string, serviceName string) {
    NewTracer(jaegerAddr, serviceName)
    tracer := TracerMiddleware()
    e.Use(tracer)
}

func NewTracer(jaegerAddr string, serviceName string) {
    cfg := jaegerConfig.Configuration{
        Sampler: &jaegerConfig.SamplerConfig{
            Type:  "const",
            Param: 1,
        },
        Reporter: &jaegerConfig.ReporterConfig{
            LogSpans:           true,
            LocalAgentHostPort: jaegerAddr,
        },
        ServiceName: serviceName,
    }
    tracer, closer, err := cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
    if err != nil {
        panic("Init jaeger err : " + err.Error())
    }

    opentracing.SetGlobalTracer(tracer)
    TracerCommon = &tracer
    TracerCloser = &closer
}

func TracerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        TracerHandler(c.Request.URL.String(), "server", c, true, make(map[string]interface{}, 0))
        c.Next()
    }
}

func TracerHandler(component string, kind string, c *gin.Context, rootSpan bool, extension map[string]interface{}) {
    if TracerCommon == nil {
        return
    }
    var parentSpan opentracing.Span
    if parentSpanContext, ok := c.Get("ParentSpanContext"); ok {
        if tracer := opentracing.GlobalTracer(); tracer != nil {
            parentSpan = tracer.StartSpan(component, opentracing.ChildOf(parentSpanContext.(opentracing.SpanContext)))
            defer parentSpan.Finish()
        }
    } else {
        spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
        if err != nil {
            fmt.Println(err)
            parentSpan = (*TracerCommon).StartSpan(c.Request.URL.Path)
            defer parentSpan.Finish()
        } else {
            parentSpan = opentracing.StartSpan(
                c.Request.URL.Path,
                opentracing.ChildOf(spCtx),
            )
            defer parentSpan.Finish()
        }
    }
    if rootSpan {
        parentSpan.SetTag("http.url", c.Request.Host+":"+c.Request.RequestURI)
        parentSpan.SetTag("http.method", c.Request.Method)
    }
    if len(extension) > 0 {
        for key, value := range extension {
            parentSpan.SetTag(key, value)
        }
    }
    parentSpan.SetTag("component", component)
    parentSpan.SetTag("span.kind", kind)
    c.Set("Tracer", TracerCommon)
    c.Set("ParentSpanContext", parentSpan.Context())
}
