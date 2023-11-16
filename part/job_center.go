package part

import (
	"fmt"
	"time"

	"github.com/anden007/af_dp_clean_core/misc"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

type JobOption struct {
	// Immediately 立即执行，当设置为true时，任务将在系统启动时立即执行，不用等待Schedule，反之亦然
	Name        string
	Immediately bool
	Schedule    cron.Schedule
}

type IJob interface {
	Option() JobOption
	Run()
}

type IJobCenter interface {
	RegisteJob(job IJob)
	StartJobs()
}

type JobCenter struct {
	Enable       bool
	cronInstance *cron.Cron
	JobArray     []IJob
}

func NewJobCenter() IJobCenter {
	// loadTime := time.Now()
	instance := &JobCenter{
		cronInstance: cron.New(),
		Enable:       viper.GetBool("cronjob.enable"),
	}
	// if ENV == ENUM_ENV_DEV {
	// misc.ServiceLoadInfo("CronJob", instance.Enable, loadTime)
	// }
	return instance
}

func (m *JobCenter) RegisteJob(job IJob) {
	option := job.Option()
	nextTime := option.Schedule.Next(time.Now())
	second := time.Until(nextTime).Seconds()
	if !m.Enable {
		misc.PrintErrorInfo(fmt.Sprintf("[CronJob] %s(每%d秒执行一次)注册失败，原因：定时任务功能尚未启用，无法执行定时任务！", option.Name, int(second)))
	} else {
		misc.PrintInfo(fmt.Sprintf("[CronJob] %s (每%d秒执行一次)，注册成功", option.Name, int(second)))
		m.JobArray = append(m.JobArray, job)
	}
}

func (m *JobCenter) StartJobs() {
	if m.Enable {
		for _, j := range m.JobArray {
			option := j.Option()
			if option.Immediately {
				// 立即执行
				j.Run()
			}
			m.cronInstance.Schedule(option.Schedule, j)
		}
		m.cronInstance.Start()
	}
}
