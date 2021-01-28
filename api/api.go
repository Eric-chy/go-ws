package api

import (
	"fmt"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/gorilla/websocket"
	"time"
	"unsafe"
)

var (
	config = gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}
)

// websocket message
type Message struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Content  string `json:"content,omitempty"`
	Uid      string `json:"Uid,omitempty"`
	CreateAt int64  `json:"createAt,omitempty"`
	UpdateAt int64  `json:"updateAt,omitempty"`
}

var redis *gredis.Redis

func init() {
	gredis.SetConfig(config, "test")
	redis = gredis.Instance("test")
}

type Controller struct {
}

//websocket
func (c *Controller) Ws(r *ghttp.Request) {
	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(err)
		r.Exit()
	}
	userId := gconv.String(r.Get("uid"))
	if userId == "" {
		r.Exit()
	}
	defer func() {
		if _, err := redis.DoVar("HDEL", "conns", userId); err != nil {
			glog.Println(err)
		}
	}()
	if _, err := redis.DoVar("HSET", "conns", userId, gconv.String(unsafe.Pointer(ws.Conn))); err != nil {
		glog.Println(userId+"redis连接失败", err)
		return
	}
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}
		if err = ws.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

//推送给单个用户  localhost:8199/send-to-user?uid=1
func (c *Controller) SendToUser(r *ghttp.Request) {
	uid := gconv.String(r.Get("uid"))
	if uid == "" {
		r.Response.WriteJson("用户id不能为空")
		r.Exit()
	}
	data := &Message{
		"系统",
		"用户" + uid,
		"你好，用户" + uid,
		uid,
		time.Now().Unix(),
		time.Now().Unix(),
	}
	json, _ := gjson.Encode(data)
	v, err := redis.DoVar("HGET", "conns", uid)
	if err != nil {
		glog.Println(err)
		//
		r.Response.WriteJson("redis连接失败")
		r.Exit()
	}
	i := v.Int()
	if i == 0 {
		fmt.Println("uid:" + uid + ":websocket连接不存在")
		r.Response.WriteJson("websocket连接不存在")
		r.Exit()
	}
	conn := *(**websocket.Conn)(unsafe.Pointer(&i))
	if err := conn.WriteMessage(1, json); err != nil {
		fmt.Println("发送消息失败: ", err)
		r.Response.WriteJson("推送失败")
		r.Exit()
	}
}

//推送给所有用户localhost:8199/send-to-users
func (c *Controller) SendToUsers(r *ghttp.Request) {
	v, err := redis.DoVar("HGETAll", "conns")
	if err != nil {
		glog.Println(err)
		return
	}
	i := v.MapStrAny()
	if i == nil {
		fmt.Println("没有连接中的websocket")
		return
	}
	for uid, connect := range i {
		data := &Message{
			"系统",
			"用户" + uid,
			"你好，用户" + uid,
			uid,
			time.Now().Unix(),
			time.Now().Unix(),
		}
		json, _ := gjson.Encode(data)
		c := gconv.Int(connect)
		conn := *(**websocket.Conn)(unsafe.Pointer(&c))
		if err := conn.WriteMessage(1, json); err != nil {
			fmt.Println("发送消息失败: ", err)
		}
	}
	return
}
