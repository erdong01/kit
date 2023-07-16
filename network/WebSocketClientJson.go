package network

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/erdong01/kit/api"
	"github.com/erdong01/kit/base"
	"github.com/erdong01/kit/common/timer"
	"github.com/erdong01/kit/rpc"
	"golang.org/x/net/websocket"
)

type WebSocketClientJson struct {
	Socket
	server   *WebSocket
	sendChan chan []byte
	timerId  *int64
}

func (w *WebSocketClientJson) Init(ip string, port int, params ...OpOption) bool {
	w.Socket.Init(ip, port, params...)
	w.timerId = new(int64)
	return true
}

func (w *WebSocketClientJson) Start() bool {
	if w.server == nil {
		return false
	}

	if w.connectType == CLIENT_CONNECT {
		w.sendChan = make(chan []byte, MAX_SEND_CHAN)
		timer.StoreTimerId(w.timerId, int64(w.clientId)+1<<32)
		timer.RegisterTimer(w.timerId, (HEART_TIME_OUT/3)*time.Second, func() {
			w.Update()
		})
	}
	if w.packetFuncList.Len() == 0 {
		w.packetFuncList = w.server.packetFuncList
	}
	if w.connectType == CLIENT_CONNECT {
		go w.SendLoop()

	}
	w.Run()
	return true
}

func (w *WebSocketClientJson) Stop() bool {
	timer.RegisterTimer(w.timerId, timer.TICK_INTERVAL, func() {
		timer.StopTimer(w.timerId)
		if atomic.CompareAndSwapInt32(&w.state, SSF_NULL, SSF_STOP) {
			if w.conn != nil {
				w.conn.Close()
			}
		}
	})
	return false
}

func (w *WebSocketClientJson) Send(head rpc.RpcHead, packet rpc.Packet) int {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	if w.connectType == CLIENT_CONNECT { //对外链接send不阻塞
		select {
		case w.sendChan <- packet.Buff:
		default: //网络太卡,tcp send缓存满了并且发送队列也满了
			w.OnNetFail(1)
		}
	} else {
		return w.DoSend(packet.Buff)
	}
	return 0
}

func (w *WebSocketClientJson) DoSend(buff []byte) int {
	if w.conn == nil {
		return 0
	}

	n, err := w.conn.Write(buff)
	handleError(err)
	if n > 0 {
		return n
	}

	return 0
}

func (s *WebSocketClientJson) SendJson(head api.JsonHead, funcName string, params ...interface{}) int {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()
	packet := rpc.MarshalJson(head, funcName, params...)
	if s.connectType == CLIENT_CONNECT { //对外链接send不阻塞
		select {
		case s.sendChan <- packet.Buff:
		default: //网络太卡,tcp send缓存满了并且发送队列也满了
			s.OnNetFail(1)
		}
	} else {
		return s.DoSend(packet.Buff)
	}
	return 0
}

func (w *WebSocketClientJson) OnNetFail(error int) {
	w.SetState(SSF_NULL)
	w.Stop()

	if w.connectType == CLIENT_CONNECT { //netgate对外格式统一
		stream := base.NewBitStream(make([]byte, 32), 32)
		stream.WriteInt(int(DISCONNECTINT), 32)
		stream.WriteInt(int(w.clientId), 32)
		w.HandlePacketJson(stream.GetBuffer())
	} else {
		w.CallMsg(rpc.RpcHead{}, "DISCONNECT", w.clientId)
	}
	if w.server != nil {
		w.server.DelClient(w)
	}
}

func (w *WebSocketClientJson) Close() {
	if w.connectType == CLIENT_CONNECT {
		//close(w.sendChan)
	}
	w.Socket.Close()
	if w.server != nil {
		w.server.DelClient(w)
	}
}

func (w *WebSocketClientJson) Run() bool {
	var buff = make([]byte, w.receiveBufferSize)
	w.SetState(SSF_RUN)
	loop := func() bool {
		defer func() {
			if err := recover(); err != nil {
				base.TraceCode(err)
			}
		}()
		if w.conn == nil {
			return false
		}
		n, err := w.conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("远程链接：已经关闭！\n")
			w.OnNetFail(0)
			return false
		}
		if err != nil {
			handleError(err)
			w.OnNetFail(0)
			return false
		}
		if n > 0 {
			w.packetParser.Read(buff[:n])
		}
		w.heartTime = int(time.Now().Unix()) + HEART_TIME_OUT
		return true
	}
	for {
		if !loop() {
			break
		}
	}
	w.Close()
	fmt.Printf("%s关闭连接", w.ip)
	return true
}

// heart
func (w *WebSocketClientJson) Update() bool {
	now := int(time.Now().Unix())
	if w.heartTime < now {
		w.OnNetFail(2)
		return false
	}
	return true
}

func (w *WebSocketClientJson) SendLoop() bool {
	for {
		defer func() {
			if err := recover(); err != nil {
				base.TraceCode(err)
			}
		}()
		select {
		case buff := <-w.sendChan:
			if buff == nil { //信道关闭
				return false
			} else {
				w.DoSend(buff)
			}
		case buff := <-w.sendChan:
			if buff == nil { //信道关闭
				return false
			} else {
				w.DoSend(buff)
			}
		}
	}
	return true
}
