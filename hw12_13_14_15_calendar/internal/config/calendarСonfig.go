package config

type CalendarConfig struct {
	Logger   LoggerConf   `json:"logger"`
	FilePath string       `json:"file_path"` //nolint:tagliatelle
	Database DataBaseConf `json:"database"`
	Net      NetConf      `json:"http"`
	Storage  StorageConf  `json:"storage"`
}

type LoggerConf struct {
	Level string `json:"level"`
}

type StorageConf struct {
	Type      string `json:"type"`
	Migration string `json:"migration"`
}

type DataBaseConf struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Dbname   string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type NetConf struct {
	API  string `json:"api"`
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

func NewCalendarConfig() CalendarConfig {
	return CalendarConfig{
		Logger:   LoggerConf{},
		FilePath: "",
		Database: DataBaseConf{},
		Net:      NetConf{},
	}
}

func (c *CalendarConfig) SetFilePath(path string) {
	c.FilePath = path
}

func (c *CalendarConfig) SetLogLevel(level string) {
	c.Logger.Level = level
}

func (c *CalendarConfig) SetDataBase(host string, port int, dbname string, username string, password string) {
	c.Database.Host = host
	c.Database.Port = port
	c.Database.Dbname = dbname
	c.Database.Username = username
	c.Database.Password = password
}
