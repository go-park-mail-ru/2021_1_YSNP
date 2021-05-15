package configs

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"

	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
)

type config struct {
	Main struct {
		URL  string
		Host string
		Port string
	}

	Auth struct {
		Host string
		Port string
	}

	Chat struct {
		Host string
		Port string
	}

	Logger struct {
		Mode string
		Host string
	}

	Postgres struct {
		User     string
		Password string
		Name     string
		Host     string
		Port     int
	}

	Tarantool struct {
		User     string
		Password string
		Host     string
		Port     int
	}

	Jaeger struct {
		Host string
		Port int
	}

	VKOAuth struct {
		AppID       string `json:"app_id" mapstructure:"app_id"`
		AppKey      string `json:"app_key" mapstructure:"app_key"`
		AppSecret   string `json:"app_secret" mapstructure:"app_secret"`
		AppUrl      string `json:"app_url" mapstructure:"app_url"`
		RedirectURL string `json:"redirect_url" mapstructure:"redirect_url"`
		FrontURL    string `json:"front_url" mapstructure:"front_url"`
	}
}

type consulConfig struct {
	Host string
	Port string
}

var (
	Configs       config
	consulConfigs consulConfig

	consulClient    *consulApi.Client
	consulLastIndex uint64 = 0

	prefix = ""

	logger = log.GetDefaultLogger()
)

func checkFlags() (string, string, string, string) {
	mode := flag.String("m", "dev", "config mode: local (read config from disk) | dev (localhost, go not in docker) | development (dev mode, go in docker) | production (ykoya config)")
	filePath := flag.String("p", ".", "file config path")
	fileName := flag.String("f", "config", "file config name")
	fileType := flag.String("t", "json", "file config type")

	flag.Parse()

	envMode := os.Getenv("APP_MODE")
	if envMode != "" {
		*mode = envMode
	}

	return *mode, *filePath, *fileName, *fileType
}

func LoadConfig() error {
	mode, filePath, fileName, fileType := checkFlags()

	viper.AddConfigPath(filePath)
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)

	logger.Info(" mode = ", mode, " filePath = ", filePath, " fileName = ", fileName, " fileType = ", fileType)

	//if err := viper.ReadInConfig(); err != nil {
	//	return err
	//}

	switch mode {
	case "local":
		if err := localConfig(); err != nil {
			return err
		}

	case "dev":
		prefix = "dev"
		consulConfigs.Host = "89.208.199.170"
		//consulConfigs.Host = "localhost"
		//consulConfigs.Host = "consul"
		consulConfigs.Port = "8500"
		if err := remoteConfig(); err != nil {
			return err
		}

	case "development":
		prefix = "development"
		consulConfigs.Host = "89.208.199.170"
		//consulConfigs.Host = "localhost"
		//consulConfigs.Host = "consul"
		consulConfigs.Port = "8500"
		if err := remoteConfig(); err != nil {
			return err
		}

	case "production":
		prefix = "production"
		consulConfigs.Host = "consul"
		consulConfigs.Port = "8500"
		if err := remoteConfig(); err != nil {
			return err
		}
	}

	logger.Info(Configs)

	return nil
}

func localConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&Configs); err != nil {
		return err
	}

	return nil
}

func remoteConfig() error {
	config := consulApi.DefaultConfig()
	config.Address = consulConfigs.getConsulAddr()

	var err error
	consulClient, err = consulApi.NewClient(config)
	if err != nil {
		return err
	}

	err = loadRemoteConfig()
	if err != nil {
		return err
	}

	//go reloadReloadConfig()

	return nil
}

func loadRemoteConfig() error {
	qo := &consulApi.QueryOptions{
		WaitIndex: consulLastIndex,
	}

	kv, qm, err := consulClient.KV().Get(prefix, qo)
	if err != nil {
		return err
	}

	if consulLastIndex == qm.LastIndex {
		return err
	}

	err = json.Unmarshal(kv.Value, &Configs)
	if err != nil {
		return err
	}

	consulLastIndex = qm.LastIndex

	return nil
}

func reloadReloadConfig() {
	ticker := time.Tick(10 * time.Second)
	for range ticker {
		if err := loadRemoteConfig(); err != nil {
			logger.Warn(err)
		}
	}
}

func (cc *consulConfig) getConsulAddr() string {
	return fmt.Sprint(cc.Host, ":", cc.Port)
}

func (c *config) GetMainUrl() string {
	return c.Main.URL
}

func (c *config) GetMainHost() string {
	return c.Main.Host
}

func (c *config) GetMainPort() string {
	return c.Main.Port
}

func (c *config) GetAuthHost() string {
	return c.Auth.Host
}

func (c *config) GetAuthPort() string {
	return c.Auth.Port
}

func (c *config) GetChatHost() string {
	return c.Chat.Host
}

func (c *config) GetChatPort() string {
	return c.Chat.Port
}

func (c *config) GetLoggerMode() string {
	return c.Logger.Mode
}

func (c *config) GetLoggerHost() string {
	return c.Logger.Host
}

func (c *config) GetPostgresConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.Name)
}

func (c *config) GetTarantoolUser() string {
	return c.Tarantool.User
}

func (c *config) GetTarantoolPassword() string {
	return c.Tarantool.Password
}

func (c *config) GetTarantoolConfig() string {
	return fmt.Sprint(c.Tarantool.Host, ":", c.Tarantool.Port)
}

func (c config) GetJaegerConfig() string {
	return fmt.Sprint(c.Jaeger.Host, ":", c.Jaeger.Port)
}

func (c *config) GetVKAppID() string {
	return c.VKOAuth.AppID
}

func (c *config) GetVKAppKey() string {
	return c.VKOAuth.AppKey
}

func (c *config) GetVKAppSecret() string {
	return c.VKOAuth.AppSecret
}

func (c *config) GetVKAppUrl() string {
	return c.VKOAuth.AppUrl
}

func (c *config) GetVKRedirectUrl() string {
	return fmt.Sprint(c.Main.URL, c.VKOAuth.RedirectURL)
}

func (c *config) GetVKFrontUrl() string {
	return c.VKOAuth.FrontURL
}
