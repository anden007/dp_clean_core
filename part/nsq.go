package part

import (
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
)

type INsq interface {
	GetProducer() (result *nsq.Producer, err error)
	GetConsumer(topic string, channel string, handler nsq.Handler) (result *nsq.Consumer, err error)
}

type Nsq struct {
	nsqAddress          string // nsq服务器地址
	nsqLookupAddress    string // nsqlookup地址
	lookupdPollInterval int64  // 消费者轮询间隔
}

func NewNsq() INsq {
	return &Nsq{
		nsqAddress:          viper.GetString("nsq.nsqAddress"),
		nsqLookupAddress:    viper.GetString("nsq.nsqLookupAddress"),
		lookupdPollInterval: viper.GetInt64("nsq.lookupd_poll_interval"),
	}
}

func (m *Nsq) GetProducer() (result *nsq.Producer, err error) {
	config := nsq.NewConfig()
	result, err = nsq.NewProducer(m.nsqAddress, config)
	return
}

func (m *Nsq) GetConsumer(topic string, channel string, handler nsq.Handler) (result *nsq.Consumer, err error) {
	config := nsq.NewConfig()
	config.LookupdPollInterval = time.Duration(m.lookupdPollInterval) * time.Second
	if result, err = nsq.NewConsumer(topic, channel, config); err == nil {
		result.AddHandler(handler)
		err = result.ConnectToNSQLookupd(m.nsqLookupAddress)
	}
	return
}
