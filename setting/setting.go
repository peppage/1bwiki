package setting

import (
	"crypto/rand"

	"github.com/mgutz/logxi/v1"
	"github.com/pelletier/go-toml"
)

var (
	logger        log.Logger
	HttpPort      string
	SessionSecret string
)

const APP_VER = "beta2"

func init() {
	logger = log.New("settings")
}

func Initialize() {
	HttpPort = "8000"
	SessionSecret = randString(20)
	c, err := toml.LoadFile("conf.toml")
	if err != nil {
		logger.Error("Error loading config", "err", err)
	} else {
		if c.Has("server.http_port") {
			HttpPort = c.Get("server.http_port").(string)
		}
		if c.Has("session.secret") {
			SessionSecret = c.Get("sedssion.secret").(string)
		}
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
