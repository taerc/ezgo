package conf

import "gopkg.in/ini.v1"

// Configure @description:
// @note: ini 不支持 embed 模式
type Configure struct {
	// [basic]
	Host         string `ini:"host"`
	Port         string `ini:"port"`
	ResourcePath string `ini:"resource_path"`
	// [mysql]
	SQL MySQLConf

	// [log]
	Log LogConf

	// sqlite
	SQLite SQLiteConf

	// dingding
	Ding DingConf
	// [sqlmonitor]
	SQLMonitor SqlMonitorConf
}

type LogConf struct {
	LogDir        string `ini:"log_dir"`
	LogFileName   string `ini:"log_filename"`
	LogMaxSize    int    `ini:"log_max_size"`
	LogMaxAge     int    `ini:"log_max_age"`
	LogMaxBackups int    `ini:"log_max_backups"`
	LogCompress   bool   `ini:"log_compress"`
	LogStderr     bool   `ini:"log_stderr"`
}

type SQLiteConf struct {
	SQLitePath string `ini:"sqlite_path"`
}

type DingConf struct {
	// [dingding]
	Token  string `ini:"dingding_token"`
	Secret string `ini:"dingding_secret"`
}

type MySQLConf struct {
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
}

// MQTTConf

type MQTTConf struct {
	MQTTHost     string `ini:"mqtt_host"`
	MQTTPort     int    `ini:"mqtt_port"`
	MQTTUser     string `ini:"mqtt_user"`
	MQTTPwd      string `ini:"mqtt_pwd"`
	MQTTSubTopic string `ini:"mqtt_sub_topic"`
}

// RedisConf

type RedisConf struct {
	RedisDSN string `ini:"redis_dsn"`
}

type SqlMonitorConf struct {
	TableSchema       string `ini:"table_schema"`
	HistoryTablePath  string `ini:"history_table_path"`
	HistoryColumnPath string `ini:"history_column_path"`
}

var Config *Configure = nil

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
		fd.Section("log").MapTo(&Config.Log)
		fd.Section("mysql").MapTo(&Config.SQL)
		fd.Section("dingding").MapTo(&Config.Ding)
		fd.Section("sqlmonitor").MapTo(&Config.SQLMonitor)
	}

	return Config, nil
}
