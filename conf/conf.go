package conf

import "gopkg.in/ini.v1"

// Configure @description:
// @note: ini 不支持 embed 模式
type Configure struct {
	// [basic]
	Host         string `ini:"host"`
	Port         string `ini:"port"`
	ResourcePath string `ini:"resource_path"`

	// [log]
	LogDir        string `ini:"log_dir"`
	LogFileName   string `ini:"log_filename"`
	LogMaxSize    int    `ini:"log_max_size"`
	LogMaxAge     int    `ini:"log_max_age"`
	LogMaxBackups int    `ini:"log_max_backups"`
	LogCompress   bool   `ini:"log_compress"`
	LogStderr     bool   `ini:"log_stderr"`

	// [mysql]
	MySQLHostname          string `ini:"mysql_hostname"`
	MySQLPort              string `ini:"mysql_port"`
	MySQLUserName          string `ini:"mysql_user"`
	MySQLPass              string `ini:"mysql_pass"`
	MySQLDBName            string `ini:"mysql_dbname"`
	MySQLMaxIdleConnection int    `ini:"max_idle_connection"`
	MySQLMaxOpenConnection int    `ini:"max_open_connection"`
	Charset                string `ini:"charset"`
	Loc                    string `ini:"loc"`
	ParseTime              string `ini:"parse_time"`
	Timeout                string `ini:"timeout"`
	MaxLifeTime            string `ini:"max_life_time"`
	TablePre               string `ini:"table_pre"`
	SlowSqlTime            string `ini:"slow_sql_time"`
	PrintSqlLog            bool   `ini:"print_sql_log"`
	// sqlite
	SQLitePath string `ini:"sqlite_path"`

	// [dingding]
	Token  string `ini:"dingding_token"`
	Secret string `ini:"dingding_secret"`
}

var Config *Configure = nil

// LogConfigure @description: 日志配置文件
type LogConfigure struct {
}

type MySQLConfigure struct {
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
		fd.Section("dingding").MapTo(Config)
	}

	return Config, nil
}
