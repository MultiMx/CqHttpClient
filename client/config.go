package client

import (
	"github.com/Mmx233/tool"
	"github.com/MultiMx/CqHttpClient/util"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type socket struct {
	Conn *websocket.Conn
	Lock *sync.RWMutex
}

type Config struct {
	HttpClient  *http.Client
	ws          *socket
	AccessKey   string
	HttpBackend string
	WsBackend   string
}

var c *Config

func Configure(a *Config) {
	if a.HttpClient == nil {
		a.HttpClient = tool.GenHttpClient(&tool.HttpClientOptions{
			Transport: tool.GenHttpTransport(&tool.HttpTransportOptions{
				Timeout: time.Second * 30,
			}),
			Timeout: time.Second * 30,
		})
	}
	util.Http = tool.NewHttpTool(a.HttpClient)
	c = a
}
