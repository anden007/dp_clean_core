package part

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/anden007/dp_clean_core/misc"
	"github.com/anden007/dp_clean_core/pkg"
	"github.com/anden007/dp_clean_core/pkg/base"

	jsoniter "github.com/json-iterator/go"
	jsonTime "github.com/liamylian/jsontime/v2/v2"
	"github.com/spf13/viper"
	"github.com/wI2L/jsondiff"
)

type IDataLogger interface {
	WriteAddLog(entity base.IDBModel, ctx context.Context, transId string) (err error)
	WriteDelLog(entity base.IDBModel, ctx context.Context, transId string) (err error)
	WriteEditLog(oldEntity, newEntity base.IDBModel, ctx context.Context, transId string) (err error)
	WriteUpdateLog(oldEntity base.IDBModel, newValues map[string]interface{}, ctx context.Context, transId string) (err error)
}

type DataUpdateLog struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type DataLogger struct {
	enable      bool
	db          IDataBase
	jsonEncoder jsoniter.API
}

type DataLog struct {
	Id         string    `gorm:"size:20;primaryKey" json:"id"`
	Table      string    `gorm:"size:100" json:"table"`
	Action     string    `gorm:"size:10" json:"action"`
	Executor   string    `gorm:"size:50" json:"executor"`
	Content    string    `gorm:"type:longText" json:"content"`
	CreateTime time.Time `json:"createTime"`
}

func (DataLog) TableName() string {
	return "t_dblog"
}

func NewDataLogger(db IDataBase) IDataLogger {
	// loadTime := time.Now()
	instance := &DataLogger{
		db:          db,
		jsonEncoder: jsonTime.ConfigWithCustomTimeFormat,
		enable:      viper.GetBool("app.data_log"),
	}
	// if ENV == ENUM_ENV_DEV {
	// 	misc.ServiceLoadInfo("DataLog", instance.enable, loadTime)
	// }
	return instance
}

func (m *DataLogger) WriteAddLog(entity base.IDBModel, ctx context.Context, transId string) (err error) {
	if m.enable {
		transId = strings.TrimSpace(transId)
		executorUserName := ""
		if executor, ok := ctx.Value("executor").(pkg.BaseUserInfo); ok {
			executorUserName = executor.UserName
		}
		if logContent, lErr := m.jsonEncoder.MarshalToString(entity); lErr == nil {
			dataLog := DataLog{
				Id:         misc.NewHexId(),
				Table:      entity.TableName(),
				Action:     "新增",
				Executor:   executorUserName,
				Content:    logContent,
				CreateTime: time.Now(),
			}
			if transId != "" {
				if tx, gErr := m.db.GetTransaction(transId); gErr == nil {
					err = tx.Create(&dataLog).Error
				} else {
					err = gErr
				}
			} else {
				err = m.db.GetDB().Create(&dataLog).Error
			}
		} else {
			err = lErr
		}
	}
	return
}

func (m *DataLogger) WriteDelLog(entity base.IDBModel, ctx context.Context, transId string) (err error) {
	if m.enable {
		transId = strings.TrimSpace(transId)
		executorUserName := ""
		if executor, ok := ctx.Value("executor").(pkg.BaseUserInfo); ok {
			executorUserName = executor.UserName
		}
		if logContent, lErr := m.jsonEncoder.MarshalToString(entity); lErr == nil {
			dataLog := DataLog{
				Id:         misc.NewHexId(),
				Table:      entity.TableName(),
				Action:     "删除",
				Executor:   executorUserName,
				Content:    logContent,
				CreateTime: time.Now(),
			}
			if transId != "" {
				if tx, gErr := m.db.GetTransaction(transId); gErr == nil {
					err = tx.Create(&dataLog).Error
				} else {
					err = gErr
				}
			} else {
				err = m.db.GetDB().Create(&dataLog).Error
			}
		} else {
			err = lErr
		}
	}
	return
}

func (m *DataLogger) WriteEditLog(oldEntity, newEntity base.IDBModel, ctx context.Context, transId string) (err error) {
	if m.enable {
		transId = strings.TrimSpace(transId)
		executorUserName := ""
		if executor, ok := ctx.Value("executor").(pkg.BaseUserInfo); ok {
			executorUserName = executor.UserName
		}
		if diff, lErr := jsondiff.Compare(oldEntity, newEntity, jsondiff.Invertible()); lErr == nil {
			diffJsonStr, _ := json.MarshalIndent(diff, "", "    ")
			dataLog := DataLog{
				Id:         misc.NewHexId(),
				Table:      oldEntity.TableName(),
				Action:     "修改",
				Executor:   executorUserName,
				Content:    string(diffJsonStr),
				CreateTime: time.Now(),
			}
			if transId != "" {
				if tx, gErr := m.db.GetTransaction(transId); gErr == nil {
					err = tx.Create(&dataLog).Error
				} else {
					err = gErr
				}
			} else {
				err = m.db.GetDB().Create(&dataLog).Error
			}
		} else {
			err = lErr
		}
	}
	return
}

func (m *DataLogger) WriteUpdateLog(oldEntity base.IDBModel, newValues map[string]interface{}, ctx context.Context, transId string) (err error) {
	if m.enable {
		transId = strings.TrimSpace(transId)
		executorUserName := ""
		if executor, ok := ctx.Value("executor").(pkg.BaseUserInfo); ok {
			executorUserName = executor.UserName
		}
		updateArray := make([]DataUpdateLog, 0)
		for k, v := range newValues {
			updateArray = append(updateArray, DataUpdateLog{
				Op:    "replace",
				Path:  fmt.Sprintf("/%s", k),
				Value: v,
			})
		}
		if logContent, lErr := m.jsonEncoder.MarshalToString(updateArray); lErr == nil {
			dataLog := DataLog{
				Id:         misc.NewHexId(),
				Table:      oldEntity.TableName(),
				Action:     "修改",
				Executor:   executorUserName,
				Content:    logContent,
				CreateTime: time.Now(),
			}
			if transId != "" {
				if tx, gErr := m.db.GetTransaction(transId); gErr == nil {
					err = tx.Create(&dataLog).Error
				} else {
					err = gErr
				}
			} else {
				err = m.db.GetDB().Create(&dataLog).Error
			}
		} else {
			err = lErr
		}
	}
	return
}
