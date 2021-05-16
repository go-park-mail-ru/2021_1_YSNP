package configs

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	consulApi "github.com/hashicorp/consul/api"
	vaultApi "github.com/hashicorp/vault/api"
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

	Vault struct {
		Host string
		Port int
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

	vaultClient *vaultApi.Client

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

		if err := remoteConsulConfig(); err != nil {
			return err
		}

	case "development":
		prefix = "development"
		consulConfigs.Host = "89.208.199.170"
		//consulConfigs.Host = "localhost"
		//consulConfigs.Host = "consul"
		consulConfigs.Port = "8500"

		if err := remoteConsulConfig(); err != nil {
			return err
		}

	case "production":
		prefix = "production"
		consulConfigs.Host = "consul"
		consulConfigs.Port = "8500"

		if err := remoteConsulConfig(); err != nil {
			return err
		}

		if err := remoteVaultConfig(); err != nil {
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

func remoteConsulConfig() error {
	var err error

	consulClient, err = consulApi.NewClient(&consulApi.Config{
		Address: consulConfigs.getConsulAddr(),
	})
	if err != nil {
		return err
	}

	err = loadConsulConfig()
	if err != nil {
		return err
	}

	//go reloadConsulConfig()

	return nil
}

func loadConsulConfig() error {
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

//func reloadConsulConfig() {
//	ticker := time.Tick(10 * time.Second)
//	for range ticker {
//		if err := loadConsulConfig(); err != nil {
//			logger.Warn(err)
//		}
//	}
//}

func remoteVaultConfig() error {
	var err error

	vaultClient, err = vaultApi.NewClient(&vaultApi.Config{
		Address: Configs.GetVaultConfig(),
	})
	if err != nil {
		return err
	}

	token := os.Getenv("VAULT_TOKEN")
	vaultClient.SetToken(token)

	err = loadVaultConfig()
	if err != nil {
		return err
	}

	return nil
}

func loadVaultConfig() error {
	postgresValues, err := vaultClient.Logical().Read("secret/data/postgres")
	if err != nil {
		return err
	}

	postgresData, ok := postgresValues.Data["data"].(map[string]interface{})
	if !ok {
		return errors.New("wrong postgres data")
	}

	Configs.Postgres.User = postgresData["user"].(string)
	Configs.Postgres.Password = postgresData["password"].(string)
	Configs.Postgres.Name = postgresData["name"].(string)

	tarantoolValues, err := vaultClient.Logical().Read("secret/data/tarantool")
	if err != nil {
		return err
	}

	tarantoolData, ok := tarantoolValues.Data["data"].(map[string]interface{})
	if !ok {
		return errors.New("wrong tarantool data")
	}

	Configs.Tarantool.User = tarantoolData["user"].(string)
	Configs.Tarantool.Password = tarantoolData["password"].(string)

	vkValues, err := vaultClient.Logical().Read("secret/data/vk")
	if err != nil {
		return err
	}

	vkData, ok := vkValues.Data["data"].(map[string]interface{})
	if !ok {
		return errors.New("wrong vk data")
	}

	Configs.VKOAuth.AppID = vkData["app_id"].(string)
	Configs.VKOAuth.AppKey = vkData["app_key"].(string)
	Configs.VKOAuth.AppSecret = vkData["app_secret"].(string)

	return nil
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

func (c config) GetVaultConfig() string {
	return fmt.Sprint("http://", c.Vault.Host, ":", c.Vault.Port)
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
