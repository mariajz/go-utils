package jaeger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerlog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var Tracer opentracing.Tracer

func InitJaeger(service string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: 100, // 100 traces per second
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "jaeger:6831",
		},
	}

	var err error
	var closer io.Closer
	Tracer, closer, err = cfg.NewTracer()
	if err != nil {
		log.Fatalf("Failed to initialize Jaeger: %v", err)
	}

	return Tracer, closer, err
}

func startSpanFromRequest(tracer opentracing.Tracer, r *http.Request, funcDesc string) opentracing.Span {
	spanCtx, _ := extract(tracer, r)
	return tracer.StartSpan(funcDesc, ext.RPCServerOption(spanCtx))
}

func inject(span opentracing.Span, request *http.Request) error {
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header))
}

func extract(tracer opentracing.Tracer, r *http.Request) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
}

func logHeaders(headers http.Header) string {
	result := ""
	for key, values := range headers {
		for _, value := range values {
			result += key + ": " + value + "\n"
		}
	}
	return result
}

func formatResult(requestBody interface{}) (formattedResult string) {
	res, _ := json.MarshalIndent(requestBody, "", "  ")
	formattedResult = string(res)
	return
}

func createRequestSpan(span opentracing.Span, c *gin.Context, requestBody interface{}) {
	span.LogFields(
		tracerlog.String("event", "request"),
		tracerlog.String("url", c.Request.URL.String()),
		tracerlog.String("method", c.Request.Method),
		tracerlog.String("headers", logHeaders(c.Request.Header)),
		tracerlog.Object("body", formatResult(requestBody)),
	)
}

func createResponseSpan(span opentracing.Span, c *gin.Context, response interface{}) {
	span.LogFields(
		tracerlog.String("event", "response"),
		tracerlog.Int("status", http.StatusCreated),
		tracerlog.Object("body", formatResult(response)),
		tracerlog.Object("headers", logHeaders(c.Writer.Header())),
	)
}

func findSpan(c *gin.Context, spanName string) (span opentracing.Span) {
	value, ok := c.Get(spanName)

	if ok {
		if spanValue, ok := value.(opentracing.Span); ok {
			span = spanValue
		} else {
			log.Println("Error: Unable to convert span to opentracing.Span")
		}
	}

	return span
}

func StartSpan(c *gin.Context, routeDesc string) *gin.Context {
	traceIDFromRequest := c.GetHeader("X-B3-TraceId")

	var span opentracing.Span
	var ctx context.Context

	if traceIDFromRequest != "" {
		traceID, _ := jaeger.TraceIDFromString(traceIDFromRequest)
		spanID := jaeger.SpanID(rand.Uint64())

		jaegerSpanContext := jaeger.NewSpanContext(
			traceID,
			spanID,
			0,
			false,
			nil,
		)

		span, ctx = opentracing.StartSpanFromContext(c.Request.Context(), routeDesc, jaeger.SelfRef(jaegerSpanContext))
	} else {
		span, ctx = opentracing.StartSpanFromContext(c.Request.Context(), routeDesc)
	}

	span.SetTag("trace_id", traceIDFromRequest)
	span.SetTag("api.url", c.Request.URL)

	c.Set("span", span)
	c.Request = c.Request.WithContext(ctx)
	c.Next()

	return c
}

func StartChildSpan(c *gin.Context, request interface{}, funcDesc string) *gin.Context {
	span := findSpan(c, "span")
	childSpan := opentracing.StartSpan(funcDesc, opentracing.ChildOf(span.Context()))

	ctx := opentracing.ContextWithSpan(c.Request.Context(), childSpan)

	createRequestSpan(childSpan, c, request)
	c.Set("child-span", childSpan)
	c.Request = c.Request.WithContext(ctx)
	c.Next()

	return c
}

func EndChildSpan(c *gin.Context, response interface{}, err error) {
	childSpan := findSpan(c, "child-span")

	defer func() {
		childSpan.Finish()
	}()

	if err != nil {
		childSpan.SetTag("error", true)
		childSpan.LogFields(tracerlog.String("event", "error"), tracerlog.String("message", err.Error()))
		return
	}
	createResponseSpan(childSpan, c, response)
	return
}

func EndSpan(c *gin.Context, err error) {
	span := findSpan(c, "span")

	defer func() {
		span.Finish()
	}()

	if err != nil {
		span.SetTag("error", true)
		span.LogFields(tracerlog.Error(err))
		return
	}
	span.SetTag("error", false)
	span.LogFields(tracerlog.Event("Success"))
	return
}
