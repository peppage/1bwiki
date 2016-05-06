package setting

import (
	"crypto/rand"
	"os"

	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/pelletier/go-toml"
)

var (
	logger        log.Logger
	HttpPort      string
	SessionSecret string
	LogLevel      string
	ServerLogging bool
	config        *toml.TomlTree
)

const APP_VER = "beta3"

func init() {
	var err error

	config, err = toml.LoadFile("conf.toml")
	if err != nil {
		log.WithError(err).Error("local config file error")
		box, err := rice.FindBox("")
		if err != nil {
			log.WithError(err).Error("can't find setup rice box")
		}
		conf, err := box.String("conf.toml")
		if err != nil {
			log.WithError(err).Error("conf.toml file error")
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
	LogLevel = "error"
	ServerLogging = false
	if config.Has("server.http_port") {
		HttpPort = config.Get("server.http_port").(string)
	}
	if config.Has("session.secret") {
		SessionSecret = config.Get("session.secret").(string)
	}
	if config.Has("server.log_level") {
		LogLevel = config.Get("server.log_level").(string)
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
