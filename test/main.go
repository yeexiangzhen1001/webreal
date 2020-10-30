package main

import (
	"fmt"
	"github.com/myxtype/webreal"
	"log"
	"time"
)

type PushingBusiness struct {
}

var (
	Sh   = webreal.NewSubscriptionHub()
	Push = PushingBusiness{}
)

func (p *PushingBusiness) OnConnect(c *webreal.Client) {


	if c.Query().Get("imei") == "" {
		fmt.Println("请输入订阅设备的Imei！")
	} else {
		fmt.Println("订阅设备的Imei为：",c.Query().Get("imei"))
		c.Write([]byte("success！"))
		fmt.Println("客户端id：",c.Id())
		fmt.Println("获取请求对象的URL.RawQuery：",c.Request().URL.RawQuery)
		fmt.Println("获取请求对象的RemoteAddr：",c.Request().RemoteAddr)
		//fmt.Println("获取请求对象的方法：",c.Request().URL.Path)
		_, Sh.Subscribers = c.Subscribe(c.Query().Get("imei"))
		fmt.Println(Sh.Subscribers)
		go func() {
			for {
				tik := time.NewTicker(time.Second*20)

				select {
				case <-tik.C:
					Sh.Publish(c.Query().Get("imei"), []byte("ok"))
				}
			}
		}()
		log.Printf("New client %d", c.Id())
	}
}

func (p *PushingBusiness) OnMessage(c *webreal.Client, msg *webreal.Message) {
	log.Printf("Client %d Message: %v", c.Id(), msg.Data)
}

func (p *PushingBusiness) OnClose(c *webreal.Client) error {
	defer c.UnsubscribeAll()
	fmt.Println("断开连接")
	log.Printf("Client %d closed.", c.Id())
	return nil
}

func main() {
	server := webreal.NewServer(&Push, Sh)
	server.Run("192.168.110.27:8080", "/ws")
}
