package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	RumCertPath  string
	RumUrl       string
	RumReadGroup string
	RumPostGroup string

	MixinBotConfig *Config
)

type (
	Config struct {
		Bot        string `json:"bot"`
		Pin        string `json:"pin"`
		ClientId   string `json:"client_id"`
		SessionId  string `json:"session_id"`
		PinToken   string `json:"pin_token"`
		PrivateKey string `json:"private_key"`
	}
)

func init() {
	path, _ := os.Getwd()

	p := filepath.Join(path, "config")
	var err error
	MixinBotConfig, err = readConfig(filepath.Join(p, "config.json"))
	if err != nil {
		log.Printf("fatal error config file: %s", err)
		os.Exit(0)
	}
	viper.SetConfigName("rum")
	viper.AddConfigPath(p)
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("fatal error rum file: %s", err)
		return
	}
	RumCertPath = viper.GetString("rum.cert.file")
	RumPostGroup = viper.GetString("rum.post.group.id")
	RumReadGroup = viper.GetString("rum.read.group.id")
	RumUrl = fmt.Sprintf("https://%s:%s", viper.GetString("rum.host"), viper.GetString("rum.port"))

}

func readConfig(name string) (*Config, error) {
	c := &Config{}
	f, err := os.Open(name)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println("找不到文件", name, err)
		return c, err
	}
	// defer the closing of our c so that we can parse it later on
	defer f.Close()

	// read our opened c as a byte array.
	byteValue, _ := ioutil.ReadAll(f)

	err = json.Unmarshal(byteValue, &c)
	if err != nil {
		fmt.Println("文件格式错误")
		return c, err
	}
	return c, nil
}
