package conf

import "gopkg.in/ini.v1"

type Configure struct {
	Host            string `ini:"host"`
	Port            string `ini:"port"`
	LogDir          string `ini:"log_dir"`
	LogMaxAge       int64  `ini:"log_max_age"`
	LogRotationTime int64  `ini:"log_rotation_time"` //日志切割时间间隔（小时）
	LogFileName     string `ini:"log_filename"`
	MySQLConfigure
}

var Config *Configure = nil

// LogConfigure @description: 日志配置文件
type LogConfigure struct {
}

type MySQLConfigure struct {
	MySQLHostname          string `ini:"mysql_hostname"`
	MySQLPort              string `ini:"mysql_port"`
	MySQLUserName          string `ini:"mysql_user"`
	MySQLPass              string `ini:"mysql_pass"`
	MySQLDBName            string `ini:"mysql_dbname"`
	MySQLMaxIdleConnection int    `ini:"max_idle_connection"`
	MySQLMaxOpenConnection int    `ini:"max_open_connection"`
	Timeout                string `ini:"timeout"`
	MaxLifeTime            string `ini:"max_life_time"`
	TablePre               string `ini:"table_pre"`
	SlowSqlTime            string `ini:"slow_sql_time"`
	PrintSqlLog            bool   `ini:"print_sql_log"`
}

func init() {
	Config = new(Configure)
	Config.Host = "127.0.0.1"
	Config.Port = "8080"
}

func LoadConfigure(filePath string) (*Configure, error) {
	if fd, e := ini.Load(filePath); e != nil {
		return nil, e
	} else {
		fd.Section("basic").MapTo(Config)
		fd.Section("log").MapTo(Config)
		fd.Section("mysql").MapTo(Config)
	}

	return Config, nil
}
