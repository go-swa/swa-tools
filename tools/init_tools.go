package tools

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Export struct {
	DbType   string `mapstructure:"db-type" json:"db-type" toml:"db-type"`
	Ip       string `mapstructure:"ip" json:"ip" toml:"ip"`
	Port     string `mapstructure:"port" json:"port" toml:"port"`
	Username string `mapstructure:"username" json:"username" toml:"username"`
	Password string `mapstructure:"password" json:"password" toml:"password"`
	Dbname   string `mapstructure:"db-name" json:"db-name" toml:"db-name"`
	Config   string `mapstructure:"config" json:"config" toml:"config"`
}
type Import struct {
	DbType   string `mapstructure:"db-type" json:"db-type" toml:"db-type"`
	Ip       string `mapstructure:"ip" json:"ip" toml:"ip"`
	Port     string `mapstructure:"port" json:"port" toml:"port"`
	Username string `mapstructure:"username" json:"username" toml:"username"`
	Password string `mapstructure:"password" json:"password" toml:"password"`
	Dbname   string `mapstructure:"db-name" json:"db-name" toml:"db-name"`
	Config   string `mapstructure:"config" json:"config" toml:"config"`
}

type Tables struct {
	TbString string `mapstructure:"table-string" json:"table-string" toml:"table-string"`
}

type ToolConfig struct {
	Export Export `mapstructure:"export" json:"export" toml:"export"`
	Import Import `mapstructure:"import" json:"import" toml:"import"`
	Tables Tables `mapstructure:"tables" json:"tables" toml:"tables"`
}


var SwaConfig ToolConfig

func InitSwaTools() error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取swa tools toml配置文件失败: %s \n", err))
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:")
		if err = v.Unmarshal(&SwaConfig); err != nil {
			fmt.Printf("系统配置数据被修改:%v", zap.Error(err))
		}
	})

	if err = v.Unmarshal(&SwaConfig); err != nil {
		fmt.Printf("系统配置数据转换失败：%v", zap.Error(err))
	}
	fmt.Printf("工具配置:\n%+v\n", SwaConfig)
	return nil
}
