package part

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/anden007/af_dp_clean_core/misc"

	"github.com/spf13/viper"
)

type EnumEnv string

const (
	// 开发环境
	ENUM_ENV_DEV EnumEnv = "dev"
	// 测试环境
	ENUM_ENV_TEST EnumEnv = "test"
	// 预发布环境
	ENUM_ENV_PRE EnumEnv = "pre"
	// 生产环境
	ENUM_ENV_PROD EnumEnv = "prod"
)

var (
	// 运行环境,可选项：dev/test/prod
	ENV               EnumEnv = ""
	ENV_LOADED        bool    = false
	USE_PPROF         bool    = false
	APP_NAME          string  = ""
	APP_HOST          string  = "0.0.0.0"
	APP_PORT          int     = 0
	RPCX_SERVICE_PORT int     = 0
	INCLUDE_DEMO      bool    = false
	once              sync.Once
)

func LoadConfig(forceEnv EnumEnv) (err error) {
	once.Do(func() {
		EnvMode := ""
		cfgFile := ""
		flag.StringVar(&APP_HOST, "h", "0.0.0.0", "主机IP")
		flag.StringVar(&EnvMode, "e", "", "运行环境,可选项:dev/test/prod")
		flag.IntVar(&APP_PORT, "p", 0, "运行端口")
		flag.IntVar(&RPCX_SERVICE_PORT, "r", 0, "rpcx服务监听端口")
		flag.BoolVar(&USE_PPROF, "d", false, "启用pprof性能测试")
		flag.Parse()

		if forceEnv != "" {
			EnvMode = string(forceEnv)
		}

		switch EnvMode {
		case "dev":
			ENV = ENUM_ENV_DEV
			cfgFile = "dev"
		case "test":
			ENV = ENUM_ENV_TEST
			cfgFile = "test"
		case "pre":
			ENV = ENUM_ENV_PRE
			cfgFile = "pre"
		case "prod":
			ENV = ENUM_ENV_PROD
			cfgFile = "prod"
		default:
			ENV = ENUM_ENV_PROD
			cfgFile = "prod"
		}

		viper.SetConfigName(cfgFile)
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")

		if cfgErr := viper.ReadInConfig(); cfgErr == nil {
			APP_NAME = viper.GetString("app.name")
			INCLUDE_DEMO = viper.GetBool("app.include_demo")
			if APP_PORT == 0 {
				if ENV == ENUM_ENV_DEV {
					APP_PORT = viper.GetInt("app.dev_port")
				} else {
					APP_PORT, err = misc.GetRndUnUsePortNumber()
					if err != nil {
						panic("随机生成程序端口失败，启动失败！")
					}
				}
			}
			if RPCX_SERVICE_PORT == 0 {
				RPCX_SERVICE_PORT = viper.GetInt("rpcx.port")
			}

			if USE_PPROF {
				_, err := os.Stat("./debug")
				if !os.IsExist(err) {
					// 文件夹不存在则创建
					_ = os.Mkdir("./debug", 0666)
				}
			}
			if ENV == ENUM_ENV_DEV {
				fmt.Println("*********************************************")
				fmt.Println("Welcome use github.com/anden007/af_dp_clean_core rapid development platform")
				fmt.Println("Version: 2.0.0")
				fmt.Println("Mode: Development")
				fmt.Printf("Host: %s\n", APP_HOST)
				fmt.Printf("Port: %d\n", APP_PORT)
				fmt.Println("*********************************************")
			}
		} else {
			if _, ok := cfgErr.(viper.ConfigFileNotFoundError); ok {
				panic(fmt.Sprintf("加载配置文件%s.toml失败, 启动失败!", cfgFile))
			} else {
				panic(fmt.Sprintf("加载配置文件%s.toml发生错误: %s, 启动失败!", cfgFile, cfgErr.Error()))
			}
		}
	})
	return
}
