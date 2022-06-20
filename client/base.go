package client

import (
	"context"
	"github.com/Mmx233/tool"
	"github.com/MultiMx/CqHttpClient/botUtil"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

func init() {
	go coolDownRoutine()
}

var coolDownChan = make(chan bool)

func coolDownRoutine() {
	for {
		<-coolDownChan
		time.Sleep(time.Millisecond * 300)
	}
}

func doWs(action string, query map[string]interface{}) error {
	return c.ws.Conn.WriteJSON(&map[string]interface{}{
		"action": action,
		"params": query,
	})
}

func do(ws bool, action string, query map[string]interface{}) (map[string]interface{}, error) {
	coolDownChan <- false
	if ws && c.ws != nil {
		c.ws.Lock.RLock()
		defer c.ws.Lock.RUnlock()
		if c.ws.Conn != nil {
			e := doWs(action, query)
			if e == nil {
				return nil, nil
			}
		}
	}
	if query == nil {
		query = make(map[string]interface{})
	}
	query["access_token"] = c.AccessKey
	_, d, e := Http.Get(&tool.DoHttpReq{
		Url:   c.HttpBackend + action,
		Query: query,
	})
	return d, e
}

func initWs() error {
	c.ws.Lock.Lock()
	defer c.ws.Lock.Unlock()
	var err error
	c.ws.Conn, _, err = websocket.DefaultDialer.Dial(
		c.WsBackend+"?access_token="+c.AccessKey,
		nil)
	return err
}

func RunWsSwitcher(Switcher func(ctx context.Context)) error {
	defer botUtil.Recover()

	if c.WsBackend != "" {
		//connect websocket
		c.ws = &socket{
			Lock: &sync.RWMutex{},
		}
	start:
		if e := initWs(); e != nil {
			log.Println("bot链接失败: ", e)
			time.Sleep(time.Second * 10)
			goto start
		}
		c.ws.Lock.RLock()
		conn := c.ws.Conn
		c.ws.Lock.RUnlock()
		for {
			var data map[string]interface{}
			err := conn.ReadJSON(&data)
			if err != nil {
				c.ws.Lock.Lock()
				_ = c.ws.Conn.Close()
				c.ws.Conn = nil
				c.ws.Lock.Unlock()
				log.Println("BOT链接中断：", err)
				goto start
			}
			var info = context.WithValue(context.Background(), "D", data)
			go Switcher(info)
		}
	}
	return nil
}
