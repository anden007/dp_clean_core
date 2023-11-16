package pkg

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	jsonTime "github.com/liamylian/jsontime/v2/v2"
	"github.com/sirupsen/logrus"
)

type ZincHook struct {
	zincServerUrl string
	zincIndex     string
	user          string
	password      string
	jsonEncoder   jsoniter.API
	levels        []logrus.Level
}

func NewZincHook(level logrus.Level, zincIndex, zincServerUrl, user, password string) (*ZincHook, error) {
	var levels []logrus.Level
	for _, l := range []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	} {
		if l <= level {
			levels = append(levels, l)
		}
	}
	jsonTime.SetDefaultTimeFormat("2006-01-02 15:04:05", time.Local)
	return &ZincHook{
		zincServerUrl: zincServerUrl,
		zincIndex:     zincIndex,
		user:          user,
		password:      password,
		levels:        levels,
		jsonEncoder:   jsonTime.ConfigWithCustomTimeFormat,
	}, nil
}

func (hook *ZincHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *ZincHook) Fire(entry *logrus.Entry) (err error) {
	client := &http.Client{}
	if fieldsJson, fErr := hook.jsonEncoder.MarshalToString(entry.Data); fErr == nil {
		var data = strings.NewReader(fmt.Sprintf(`{ "level": "%s", "msg": "%s", "fields": %s }`, entry.Level, entry.Message, fieldsJson))
		if req, rErr := http.NewRequest("PUT", fmt.Sprintf("%s/api/%s/_doc", hook.zincServerUrl, hook.zincIndex), data); rErr == nil {
			req.Header.Set("Content-Type", "application/json")
			req.SetBasicAuth(hook.user, hook.password)
			_, err = client.Do(req)
		} else {
			err = rErr
		}
	} else {
		err = fErr
	}
	return
}
