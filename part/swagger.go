package part

import (
	"fmt"
	"strings"

	"github.com/anden007/af_dp_clean_core/misc"

	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
	"github.com/swaggo/swag"
)

type ISwaggerService interface {
}

type SwaggerService struct {
}

func NewSwaggerService(app *iris.Application, spec *swag.Spec) ISwaggerService {
	instance := &SwaggerService{}
	swaggerEnable := viper.GetBool("swagger.enable")
	if swaggerEnable {
		vhost := viper.GetString("app.vhost")
		if !strings.HasSuffix(vhost, "/") {
			vhost = fmt.Sprintf("%s/", vhost)
		}
		config := &swagger.Config{
			URL: fmt.Sprintf("%s%s", vhost, "swagger/doc.json"), //The url pointing to API definition
		}
		if spec != nil {
			spec.Title = viper.GetString("app.name")
			spec.BasePath = vhost
		}
		app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(config, swaggerFiles.Handler))
		fmt.Println("Swagger: Enabled")
		fmt.Printf("Swagger Address: http://127.0.0.1:%d/swagger/index.html\n", APP_PORT)
	} else {
		misc.PrintErrorInfo("Swagger 尚未启用")
	}
	return instance
}
