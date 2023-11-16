package misc

import (
	"context"

	"github.com/kataras/iris/v12"
	jsonTime "github.com/liamylian/jsontime/v2/v2"
)

func WriteJson(ctx iris.Context, v interface{}) {
	jsonEncoder := jsonTime.ConfigWithCustomTimeFormat
	if body, err := jsonEncoder.Marshal(v); err == nil {
		ctx.ContentType("application/json;charset=utf-8")
		ctx.Write(body)
	} else {
		ctx.StopWithError(500, err)
	}
}

func GetExecutorContext(executor string) context.Context {
	return context.WithValue(context.TODO(), "executor", executor)
}
