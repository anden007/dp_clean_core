package part

import (
	"fmt"

	"github.com/anden007/dp_clean_core/misc"
	"github.com/fieldryand/goflow/v2"
	gokv_redis "github.com/philippgille/gokv/redis"

	"github.com/spf13/viper"
)

type IJobCenter interface {
	RegisteJob(job func() *goflow.Job)
	Run()
}

type JobCenter struct {
	Enable     bool
	gfInstance *goflow.Goflow
}

func NewJobCenter() IJobCenter {
	options := goflow.Options{
		UIPath:       "ui/",
		ShowExamples: false,
		WithSeconds:  true,
	}
	gf := goflow.New(options)
	if gokv_redis_client, grcErr := gokv_redis.NewClient(gokv_redis.Options{
		Address:  viper.GetString("redis_cache.server"),
		Password: viper.GetString("redis_cache.password"),
		DB:       viper.GetInt("redis_cache.db"),
	}); grcErr == nil {
		defer gokv_redis_client.Close()
		gf.Store = gokv_redis_client
	} else {
		misc.PrintErrorInfo(fmt.Sprintf("[JobCenter] GoFlow 初始化redis存储发生错误: %s", grcErr.Error()))
	}
	gf.Use(goflow.DefaultLogger())

	instance := &JobCenter{
		gfInstance: gf,
		Enable:     viper.GetBool("job_center.enable"),
	}
	return instance
}

func (m *JobCenter) RegisteJob(job func() *goflow.Job) {
	m.gfInstance.AddJob(job)
}

func (m *JobCenter) Run() {
	if m.Enable {
		gfPort := viper.GetInt("job_center.port")
		go m.gfInstance.Run(fmt.Sprintf(":%d", gfPort))
	}
}
