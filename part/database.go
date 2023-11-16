package part

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type IDataBase interface {
	GetDB() *gorm.DB
	GetTransaction(transId string) (result *gorm.DB, err error)
	Commit(transId string) (err error)
	Rollback(transId string) (err error)
}

type DataBase struct {
	db              *gorm.DB
	transactionPool ILruCache
	Enable          bool
}

func NewDataBase(lruCache ILruCache) IDataBase {
	instance := &DataBase{
		transactionPool: lruCache,
	}
	// loadTime := time.Now()
	enable := viper.GetBool("mysql.enable")
	if enable {
		instance.Enable = true
		logLevel := logger.Silent
		if ENV == ENUM_ENV_DEV {
			logLevel = logger.Info
		}
		myLogger := logger.Default.LogMode(logLevel)
		dsn := viper.GetString("mysql.server")
		if mySqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger:                                   myLogger,
			DisableForeignKeyConstraintWhenMigrating: true,
		}); err == nil {
			mySqlDB.Set("gorm:table_options", fmt.Sprintf("ENGINE=%s DEFAULT CHARSET=%s", viper.GetString("mysql.engine"), viper.GetString("mysql.charset")))
			if db, err := mySqlDB.DB(); err == nil {
				db.SetMaxIdleConns(10)
				db.SetMaxOpenConns(100)
				db.SetConnMaxLifetime(time.Minute * 30)
			}
			instance.db = mySqlDB
		} else {
			fmt.Println("Connect to MySQL error", err)
			return nil
		}
	}

	// if ENV == ENUM_ENV_DEV {
	// misc.ServiceLoadInfo("DataBase", enable, loadTime)
	// }
	return instance
}

func (m *DataBase) checkMe() {
	if !m.Enable {
		panic("DataBase is Disabled.\n")
	}
}

func (m *DataBase) GetDB() *gorm.DB {
	m.checkMe()
	return m.db
}

func (m *DataBase) GetTransaction(transId string) (result *gorm.DB, err error) {
	m.checkMe()
	if exists := m.transactionPool.Contains(transId); exists {
		if intface, success := m.transactionPool.Get(transId); success {
			result = intface.(*gorm.DB)
		}
	}
	if result == nil {
		result = m.db.Begin().SavePoint(transId)
		m.transactionPool.Add(transId, result)
	}
	return
}

func (m *DataBase) Commit(transId string) (err error) {
	if tx, gErr := m.GetTransaction(transId); gErr == nil {
		err = tx.Commit().Error
		m.transactionPool.Remove(transId)
	} else {
		err = gErr
	}
	return
}

func (m *DataBase) Rollback(transId string) (err error) {
	if tx, gErr := m.GetTransaction(transId); gErr == nil {
		err = tx.RollbackTo(transId).Error
		m.transactionPool.Remove(transId)
	} else {
		err = gErr
	}
	return
}
