package setting

import (
	"github.com/mgutz/logxi/v1"
	"github.com/pelletier/go-toml"
)

var (
	logger   log.Logger
	HttpPort string
)

const APP_VER = "beta1"

func init() {
	logger = log.New("settings")
}

func Initialize() {
	c, err := toml.LoadFile("conf.toml")
	if err != nil {
		logger.Error("Error loading config", "err", err)
		HttpPort = "8000"
	} else {
		if c.Has("server.http_port") {
			HttpPort = c.Get("server.http_port").(string)
		} else {
			HttpPort = "8000"
		}
	}
}
