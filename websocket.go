package main

import (
	"encoding/json"
	"net"
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

type WsPool struct {
	sync.Mutex
	Peers map[net.Conn]peer
}

type peer struct {
	pubs map[topic]WsMsg
	subs map[topic]struct{}
}

var wsPool = &WsPool{Peers: map[net.Conn]peer{}}

func (w *WsPool) SetSub(c net.Conn, t topic) {
	w.Lock()
	defer w.Unlock()
	if p, ok := w.Peers[c]; ok {
		p.subs[t] = struct{}{}
	} else {
		w.Peers[c] = peer{subs: map[topic]struct{}{t: {}}, pubs: map[topic]WsMsg{}}
	}
}

func (w *WsPool) getSub(t topic) net.Conn {
	w.Lock()
	defer w.Unlock()
	for c, p := range w.Peers {
		if _, ok := p.subs[t]; ok {
			return c
		}
	}
	return nil
}

func (w *WsPool) removePeer(c net.Conn) {
	w.Lock()
	defer w.Unlock()
	delete(w.Peers, c)
	if err := c.Close(); err != nil {
		log.Errorf("conn close error: %s", err.Error())
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

func (w *WsPool) getPub(t topic) *WsMsg {
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

	pubMsg := wsPool.getPub(t)
	if pubMsg == nil {
		return
	}

	if err := wsutil.WriteServerText(conn, pubMsg.Marshal()); err != nil {
		log.Error("subscribe controller write err", err.Error())
	} else {
		log.Info("outgoing msg1", conn.RemoteAddr(), string(pubMsg.Marshal()))
	}
}

func PublishController(conn net.Conn, msg WsMsg) {
	if c := wsPool.getSub(msg.Topic); c != nil {
		if err := wsutil.WriteServerText(c, msg.Marshal()); err != nil {
			log.Error("publish controller write err", err.Error())
		} else {
			log.Info("outgoing msg2", c.RemoteAddr(), string(msg.Marshal()))
		}
	} else {
		wsPool.SetPub(conn, msg)
	}
}

func WebSocketHandler(ctx *gin.Context) {
	conn, _, _, err := ws.UpgradeHTTP(ctx.Request, ctx.Writer)
	if err != nil {
		log.Errorf("web socket connect err: %s", err.Error())
		return
	}

	go func() {
		defer wsPool.removePeer(conn)
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
			log.Info("incoming", conn.RemoteAddr(), string(msg.Marshal()))
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
