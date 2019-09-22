package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

var configPath string

// Config 保存配置信息
type Config struct {
	ListenAddr string `json:"listen"`
	RemoteAddr string `json:"remote"`
	Password   string `json:"password"`
}

// 初始化配置文件路径
func init() {
	home, _ := homedir.Dir()
	configFileName := ".socks.json"
	if len(os.Args) == 2 {
		configFileName = os.Args[1]
	}
	configPath = path.Join(home, configFileName)
}

// SaveConfig 保存配置信息
func (config *Config) SaveConfig() {
	configJSON, _ := json.MarshalIndent(config, "", " ")
	err := ioutil.WriteFile(configPath, configJSON, 0644)
	if err != nil {
		log.Printf("保存配置到文件 %s 出错: %s", configPath, err)
	}
	log.Printf("保存配置到文件 %s 成功\n", configPath)
}

// ReadConfig 读取配置信息
func (config *Config) ReadConfig() {
	// 如果配置文件存在，就读取配置文件中的配置 assign 到 config
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		log.Printf("从文件 %s 中读取配置\n", configPath)
		file, err := os.Open(configPath)
		if err != nil {
			log.Fatalf("打开配置文件 %s 出错:%s", configPath, err)
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(config)
		if err != nil {
			log.Fatalf("格式不合法的 JSON 配置文件:\n%s", file.Name())
		}
	}
}
