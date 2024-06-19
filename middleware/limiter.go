package middleware

import (
	"fmt"
	"sync"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

var (
	limiterMap      sync.Map
	httpPassCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_pass_count",
			Help: "http pass count",
		},
		[]string{"service", "path"},
	)
	httpReachMaxLimitationCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_reach_max_limitation_count",
			Help: "reach max limitation count",
		},
		[]string{"service", "path"},
	)
)

func NewLimiter(max float64) iris.Handler {
	return func(ctx iris.Context) {
		APP_NAME := viper.GetString("app.name")
		routeName := ctx.RouteName()
		mapResult, exists := limiterMap.Load(routeName)
		if !exists || mapResult == nil {
			mapResult = tollbooth.NewLimiter(max, nil)
			limiterMap.Store(routeName, mapResult)
		}
		limiter := mapResult.(*limiter.Limiter)
		if limiter != nil {
			httpError := tollbooth.LimitByRequest(limiter, ctx.ResponseWriter(), ctx.Request())
			if httpError != nil {
				httpReachMaxLimitationCounter.WithLabelValues(APP_NAME, routeName).Inc()
				ctx.StatusCode(httpError.StatusCode)
				ctx.ContentType("application/json;charset=UTF-8")
				ctx.WriteString(fmt.Sprintf("{\"success\":false,\"message\":\"%s\",\"result\":\"\"}", httpError.Message))
				ctx.StopExecution()
				return
			} else {
				httpPassCounter.WithLabelValues(APP_NAME, routeName).Inc()
			}
		}
		ctx.Next()
	}
}

func init() {
	prometheus.MustRegister(httpPassCounter)
	prometheus.MustRegister(httpReachMaxLimitationCounter)
}
