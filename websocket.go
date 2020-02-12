package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

const (
	WsSubEvent = "sub"
	WsPubEvent = "pub"
)

type topic string

type WsMsg struct {
	Topic   topic  `json:"topic"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
	Silent  bool   `json:"silent"`
}

func (msg *WsMsg) Marshal() []byte {
	bytes, _ := json.Marshal(msg)
	return bytes
}

func (msg *WsMsg) String() string {
	return fmt.Sprintf("topic: %s, type: %s, silent: %v", msg.Topic, msg.Type, msg.Silent)
}

type WsPool struct {
	sync.Mutex
	Peers map[net.Conn]peer
}

type peer struct {
	pubs map[topic]WsMsg
	subs map[topic]struct{}
}

var wsPool = &WsPool{Peers: make(map[net.Conn]peer)}

func (w *WsPool) SetSub(c net.Conn, t topic) {
	w.Lock()
	defer w.Unlock()
	if p, ok := w.Peers[c]; ok {
		p.subs[t] = struct{}{}
	} else {
		w.Peers[c] = peer{subs: map[topic]struct{}{t: {}}, pubs: map[topic]WsMsg{}}
	}
}

func (w *WsPool) GetSub(t topic) net.Conn {
	w.Lock()
	defer w.Unlock()
	for c, p := range w.Peers {
		if _, ok := p.subs[t]; ok {
			return c
		}
	}
	return nil
}

func (w *WsPool) RemovePeer(c net.Conn) {
	w.Lock()
	defer w.Unlock()
	_, ok := w.Peers[c]
	if ok {
		log.Info("find peer: ", c.RemoteAddr())
		delete(w.Peers, c)
	}
	if err := c.Close(); err != nil {
		log.Errorf("close conn error: %s", err.Error())
	}
}

func (w *WsPool) SetPub(c net.Conn, msg WsMsg) {
	w.Lock()
	defer w.Unlock()
	if p, ok := w.Peers[c]; ok {
		p.pubs[msg.Topic] = msg
	} else {
		w.Peers[c] = peer{pubs: map[topic]WsMsg{msg.Topic: msg}, subs: map[topic]struct{}{}}
	}
}

func (w *WsPool) GetPub(t topic) *WsMsg {
	w.Lock()
	defer w.Unlock()
	for _, p := range w.Peers {
		if msg, ok := p.pubs[t]; ok {
			delete(p.pubs, t)
			return &msg
		}
	}
	return nil
}

func SubscribeController(conn net.Conn, t topic) {
	wsPool.SetSub(conn, t)

	pubMsg := wsPool.GetPub(t)
	if pubMsg == nil {
		return
	}

	if err := wsutil.WriteServerText(conn, pubMsg.Marshal()); err != nil {
		log.Error("subscribe controller write err", err.Error())
	} else {
		log.Info("outgoing msg1", conn.RemoteAddr(), pubMsg.String())
	}
}

func PublishController(conn net.Conn, msg WsMsg) {
	if c := wsPool.GetSub(msg.Topic); c != nil {
		if err := wsutil.WriteServerText(c, msg.Marshal()); err != nil {
			log.Error("publish controller write err", err.Error())
		} else {
			log.Info("outgoing msg2", c.RemoteAddr(), msg.String())
		}
	} else {
		wsPool.SetPub(conn, msg)
	}
}

func WebSocketHandler(ctx *gin.Context) {
	if !ctx.IsWebsocket() {
		ctx.String(http.StatusBadRequest, "This is a the websocket API")
		return
	}
	conn, _, _, err := ws.UpgradeHTTP(ctx.Request, ctx.Writer)
	if err != nil {
		log.Errorf("web socket connect err: %s", err.Error())
		return
	}
	log.Info("success to peer connection:", conn.RemoteAddr())

	go func() {
		defer wsPool.RemovePeer(conn)
		for {
			msgBt, err := wsutil.ReadClientText(conn)
			if err != nil {
				log.Errorf("web socket read client data err: %s", err.Error())
				return
			}
			var msg WsMsg
			if err = json.Unmarshal(msgBt, &msg); err != nil {
				log.Errorf("web socket unmarshal msg err: %s", err.Error())
				break
			}
			log.Info("incoming", conn.RemoteAddr(), msg.String())
			switch msg.Type {
			case WsSubEvent:
				SubscribeController(conn, msg.Topic)
			case WsPubEvent:
				PublishController(conn, msg)
			default:
				log.Warning("web socket msg type is not exist!", msg.Type)
			}
		}
	}()
}
