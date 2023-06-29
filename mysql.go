package ezgo

import (
	"fmt"
	"github.com/taerc/ezgo/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
	"time"
)

var M string = "EZGO"
var myClient *gorm.DB = nil

// ConnectMysql
// @description:连接mysql
func initMySQL(c *conf.Configure) error {

	// 用户名:密码@tcp(IP:port)/数据库?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		c.MySQLUserName, c.MySQLPass, c.MySQLHostname, c.MySQLPort,
		c.MySQLDBName, c.Charset, c.ParseTime, c.Loc)

	Info(nil, M, dsn)

	// 连接额外配置信息
	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //使用单数表名，启用该选项时，`User` 的表名应该是 `user`而不是users
		},
	}
	// 打印SQL设置
	//if MysqlConfigInstance.PrintSqlLog {
	//	slowSqlTime, err := time.ParseDuration(MysqlConfigInstance.SlowSqlTime)
	//	if nil != err {
	//		log.Errorln("打印SQL设置失败：", err)
	//		return err
	//	}
	//	loggerNew := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
	//		SlowThreshold: slowSqlTime, //慢SQL阈值
	//		LogLevel:      logger.Info,
	//		Colorful:      true, // 彩色打印开启
	//	})
	//	gormConfig.Logger = loggerNew
	//}
	var err error
	// 建立连接
	myClient, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gormConfig)
	if nil != err {
		Error(nil, M, err.Error())
		return err
	}
	Info(nil, M, "open database")
	// 设置连接池信息
	db, err2 := myClient.DB()
	if nil != err2 {
		Error(nil, M, err2.Error())
		return err2
	}
	Info(nil, M, "open db")
	//// 设置空闲连接池中连接的最大数量
	db.SetMaxIdleConns(c.MySQLMaxIdleConnection)
	//// 设置打开数据库连接的最大数量
	db.SetMaxOpenConns(c.MySQLMaxOpenConnection)
	//// 设置了连接可复用的最大时间
	duration, err3 := time.ParseDuration(c.MaxLifeTime)
	if nil != err3 {
		Error(nil, M, err3.Error())
		return err3
	}
	Info(nil, M, "open duration")
	db.SetConnMaxLifetime(duration)
	return nil
}

func DB() *gorm.DB {
	return myClient
}

func WithComponentMySQL(c *conf.Configure) Component {
	return func(wg *sync.WaitGroup) {
		wg.Done()
		initMySQL(c)
		Info(nil, M, "Finished Load MySQL !")
	}
}
