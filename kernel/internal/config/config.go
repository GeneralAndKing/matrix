package config

import "kernel/pkg/config"

type Application struct {
	Debug bool   `mapstructure:"debug"`
	Port  int    `mapstructure:"port"`
	Env   string `mapstructure:"env"`
}

type Datasource struct {
	Sqlite3 string `mapstructure:"sqlite3"`
}

type Log struct {
	AccessLog struct {
		OutputPath    string `mapstructure:"output-path"`
		ErrOutputPath string `mapstructure:"err-output-path"`
	} `mapstructure:"access-log"`
	ErrorLog struct {
		OutputPath    string `mapstructure:"output-path"`
		ErrOutputPath string `mapstructure:"err-output-path"`
	} `mapstructure:"error-log"`
	MaxSize    int  `mapstructure:"max-size"`
	MaxAge     int  `mapstructure:"max-age"`
	MaxBackups int  `mapstructure:"max-backups"`
	Compress   bool `mapstructure:"compress"`
}

type global struct {
	Log         Log `mapstructure:"log"`
	Datasource  `mapstructure:"datasource"`
	Application `mapstructure:",squash"`
}

var (
	g = global{}
)

func GetDataSource() Datasource {
	return g.Datasource
}

func GetApplication() Application {
	return g.Application
}

func init() {
	if err := config.UnmarshalConfig(&g, "matrix", "matrix"); err != nil {
		panic(err)
	}
}
