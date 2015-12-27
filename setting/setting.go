package setting

import (
	"crypto/rand"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/mgutz/logxi/v1"
	"github.com/pelletier/go-toml"
)

var (
	logger        log.Logger
	HttpPort      string
	SessionSecret string
	LogLevel      int
	ServerLogging bool
	config        *toml.TomlTree
)

const APP_VER = "beta2"

func init() {
	logger = log.New("settings")
	var err error

	config, err = toml.LoadFile("conf.toml")
	if err != nil {
		logger.Error("local config file error", "err", err)
		box, err := rice.FindBox("")
		if err != nil {
			logger.Error("can't find setup rice box", "err", err)
		}
		conf, err := box.String("conf.toml")
		if err != nil {
			logger.Error("conf.toml file error", "err", err)
		}
		config, err = toml.Load(conf)
		f, err := os.Create("conf.toml")
		defer f.Close()
		f.Write([]byte(conf))
	}
}

func Initialize() {
	HttpPort = "8000"
	SessionSecret = randString(20)
	LogLevel = 3
	ServerLogging = false
	if config.Has("server.http_port") {
		HttpPort = config.Get("server.http_port").(string)
	}
	if config.Has("session.secret") {
		SessionSecret = config.Get("session.secret").(string)
	}
	if config.Has("server.log_level") {
		LogLevel = int(config.Get("server.log_level").(int64))
		logger.SetLevel(LogLevel)
	}
	if config.Has("server.server_logging") {
		ServerLogging = config.Get("server.server_logging").(bool)
	}
}

func randString(size int) string {
	const alpha = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, size)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = alpha[v%byte(len(alpha))]
	}
	return string(bytes)
}
