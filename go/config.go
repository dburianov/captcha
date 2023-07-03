package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"os"
)

var (
	loggingLevel   *string
	configFileJson *string
	seviceConfig   Config
	bindAddress    string
)

type Config struct {
	Salt        string        `json:"salt,omitempty"`
	PrimeNumber int64         `json:"primenumber,omitempty"`
	AuthExpire  time.Duration `json:"authexpire,omitempty"`
	Rand        struct {
		RandLenght    int `json:"rand_lenght,omitempty"`
		RandCount     int `json:"rand_count,omitempty"`
		CharacterUp   int `json:"characterUp,omitempty"`
		CharacterDown int `json:"characterDown,omitempty"`
		CharacterInt  int `json:"characterInt,omitempty"`
		CharacterSpec int `json:"characterSpec,omitempty"`
	} `json:"rand,omitempty"`
	Captcha struct {
		TextLength   int    `json:"textlength,omitempty"`
		Salt         string `json:"salt,omitempty"`
		EncryptTimes int16  `json:"encrypttimes,omitempty"`
		CharPreset   string `json:"charpreset,omitempty"`
		Text         struct {
			Width int `json:"width,omitempty"`
			Hight int `json:"hight,omitempty"`
		} `json:"text,omitempty"`
		Math struct {
			Width int `json:"width,omitempty"`
			Hight int `json:"hight,omitempty"`
		} `json:"math,omitempty"`
		Cookie struct {
			Expire time.Duration `json:"expire,omitempty"`
		} `json:"cookie,omitempty"`
	} `json:"captcha,omitempty"`
	LoggingLevel string `json:"logginglevel,omitempty"`
	Bind         struct {
		Port      uint16 `json:"port,omitempty"`
		IpAddress string `json:"ip,omitempty"`
	} `json:"bind,omitempty"`
}

func LoadConfiguration(file string) Config {
	config := new(Config).Init()
	//	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	//jsonParser := json.NewDecoder(configFile)
	//jsonParser.Decode(&config)
	byteValue, _ := ioutil.ReadAll(configFile)
	if *loggingLevel == "debug" {
		log.Println(string([]byte(byteValue)))
	}

	json.Unmarshal([]byte(byteValue), &config)
	if *loggingLevel == "debug" {
		log.Println(config)
	}

	return config
}

func init() {
	configFileJson = flag.String("config.file", "config.json", "a string")
	loggingLevel = flag.String("logging.level", "info", "a string")
	flag.Parse()

	seviceConfig = LoadConfiguration(*configFileJson)
	bindAddress = seviceConfig.Bind.IpAddress + ":" + strconv.FormatUint(uint64(seviceConfig.Bind.Port), 10)

}

func (cfg Config) Init() Config {
	cfg.Captcha.TextLength = 4
	cfg.Captcha.Salt = "localhost"
	cfg.Captcha.EncryptTimes = 100 // do not use now
	cfg.Captcha.CharPreset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	cfg.Bind.IpAddress = "0.0.0.0"
	cfg.Bind.Port = 18781
	cfg.PrimeNumber = 115969591
	cfg.AuthExpire = 86400 // second
	cfg.Rand.RandLenght = 24
	cfg.Rand.RandCount = 1
	cfg.Rand.CharacterUp = 8
	cfg.Rand.CharacterDown = 8
	cfg.Rand.CharacterInt = 8
	cfg.Rand.CharacterSpec = 8
	cfg.Salt = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	cfg.LoggingLevel = "info"
	cfg.Captcha.Text.Hight = 50
	cfg.Captcha.Text.Width = 150
	cfg.Captcha.Math.Hight = 50
	cfg.Captcha.Math.Width = 150
	cfg.Captcha.Cookie.Expire = 120 // second
	return cfg
}
