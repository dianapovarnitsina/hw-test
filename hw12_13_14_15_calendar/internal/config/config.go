package config

type Config struct {
	Logger   LoggerConf   `json:"logger"`
	FilePath string       `json:"file_path"`
	Database DataBaseConf `json:"database"`
	Http     HTTPConf     `json:"http"`
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

type HTTPConf struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

func NewConfig() Config {
	return Config{
		Logger:   LoggerConf{},
		FilePath: "",
		Database: DataBaseConf{},
		Http:     HTTPConf{},
	}
}

func (c *Config) SetFilePath(path string) {
	c.FilePath = path
}

func (c *Config) SetLogLevel(level string) {
	c.Logger.Level = level
}

func (c *Config) SetDataBase(host string, port int, dbname string, username string, password string) {
	c.Database.Host = host
	c.Database.Port = port
	c.Database.Dbname = dbname
	c.Database.Username = username
	c.Database.Password = password
}

func (c *Config) setHttp(host string, port uint16) {
	c.Http.Host = host
	c.Http.Port = port
}
