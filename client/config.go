package client

import (
	"github.com/gorilla/websocket"
	"sync"
)

type socket struct {
	Conn *websocket.Conn
	Lock *sync.RWMutex
}

type Config struct {
	ws          *socket
	AccessKey   string
	HttpBackend string
	WsBackend   string
}

var c *Config

func Configure(a *Config) {
	c = a
}
