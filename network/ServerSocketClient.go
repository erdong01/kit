package network

import (
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"sync/atomic"
	"time"

	"github.com/erdong01/kit/base"
	"github.com/erdong01/kit/base/timer"
	"github.com/erdong01/kit/rpc"
	"google.golang.org/protobuf/proto"
)

const (
	IDLE_TIMEOUT    = iota
	CONNECT_TIMEOUT = iota
	CONNECT_TYPE    = iota
)

var (
	DISCONNECTINT = crc32.ChecksumIEEE([]byte("DISCONNECT"))
	HEART_PACKET  = crc32.ChecksumIEEE([]byte("heardpacket"))
)

type IServerSocketClient interface {
	ISocket
}

type ServerSocketClient struct {
	Socket
	server   *ServerSocket
	sendChan chan []byte //对外缓冲队列
	timerId  *int64
	Property any
}

func handleError(err error) {
	if err == nil {
		return
	}
	log.Printf("错误：%s\n", err.Error())
}

func (s *ServerSocketClient) Init(ip string, port int, params ...OpOption) bool {
	s.timerId = new(int64)

	s.Socket.Init(ip, port, params...)
	return true
}

func (s *ServerSocketClient) StartJson() bool {
	if s.server == nil {
		return false
	}
	if s.connectType == CLIENT_CONNECT {
		s.sendChan = make(chan []byte, MAX_SEND_CHAN)
		timer.StoreTimerId(s.timerId, int64(s.clientId)+1<<32)
		timer.RegisterTimer(s.timerId, (HEART_TIME_OUT/2)*time.Second, func() {
			s.Update()
		})
	}
	if s.packetFuncList.Len() == 0 {
		s.packetFuncList = s.server.packetFuncList
	}
	//s.m_Conn.SetKeepAlive(true)
	//s.m_Conn.SetKeepAlivePeriod(5*time.Second)
	go s.Run()
	if s.connectType == CLIENT_CONNECT {
		go s.SendLoop()
	}
	return true
}

func (s *ServerSocketClient) Start() bool {
	if s.server == nil {
		return false
	}
	if s.connectType == CLIENT_CONNECT {
		s.sendChan = make(chan []byte, MAX_SEND_CHAN)
		timer.StoreTimerId(s.timerId, int64(s.clientId)+1<<32)
		timer.RegisterTimer(s.timerId, (HEART_TIME_OUT/2)*time.Second, func() {
			s.Update()
		})
	}
	if s.packetFuncList.Len() == 0 {
		s.packetFuncList = s.server.packetFuncList
	}
	//s.m_Conn.SetKeepAlive(true)
	//s.m_Conn.SetKeepAlivePeriod(5*time.Second)
	go s.Run()
	if s.connectType == CLIENT_CONNECT {
		go s.SendLoop()
	}
	return true
}

func (s *ServerSocketClient) Send(head rpc.RpcHead, packet rpc.Packet) int {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

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

func (s *ServerSocketClient) DoSend(buff []byte) int {
	if s.conn == nil {
		return 0
	}

	n, err := s.conn.Write(s.packetParser.Write(buff))
	handleError(err)
	if n > 0 {
		return n
	}

	return 0
}

func (s *ServerSocketClient) OnNetFail(error int) {
	s.Stop()
	if s.connectType == CLIENT_CONNECT { //netgate对外格式统一
		stream := base.NewBitStream(make([]byte, 32), 32)
		stream.WriteInt(int(DISCONNECTINT), 32)
		stream.WriteInt(int(s.clientId), 32)
		s.HandlePacket(stream.GetBuffer())
	} else {
		s.CallMsg(rpc.RpcHead{}, "DISCONNECT", s.clientId)
	}
	if s.server != nil {
		s.server.DelClinet(s)
	}
}

func (s *ServerSocketClient) Stop() bool {
	// timer.RegisterTimer(s.timerId, timer.TICK_INTERVAL, func() {
	// 	timer.StopTimer(s.timerId)
	time.Sleep(timer.TICK_INTERVAL)
	if atomic.CompareAndSwapInt32(&s.state, SSF_RUN, SSF_STOP) {
		if s.conn != nil {
			s.conn.Close()
		}
	}
	// })
	return false
}

func (s *ServerSocketClient) Close() {
	if s.connectType == CLIENT_CONNECT {
		s.sendChan <- nil
		//close(s.sendChan)
	}
	s.Socket.Close()
	if s.server != nil {
		s.server.DelClinet(s)
	}
}

func (s *ServerSocketClient) Run() bool {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()
	var buff = make([]byte, s.receiveBufferSize)
	s.SetState(SSF_RUN)
	loop := func() bool {
		defer func() {
			if err := recover(); err != nil {
				base.TraceCode(err)
			}
		}()

		if s.conn == nil {
			return false
		}

		n, err := s.conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("远程链接：%s已经关闭！\n", s.conn.RemoteAddr().String())
			s.OnNetFail(0)
			return false
		}
		if err != nil {
			handleError(err)
			s.OnNetFail(0)
			return false
		}
		if n > 0 {
			//熔断
			if !s.packetParser.Read(buff[:n]) && s.connectType == CLIENT_CONNECT {
				s.OnNetFail(1)
				return false
			}
		}
		s.heartTime = int(time.Now().Unix()) + HEART_TIME_OUT

		return true
	}

	for {
		if !loop() {
			break
		}
	}
	if s.server.clientClose != nil {
		s.server.clientClose(s.clientId)
	}
	s.Close()
	fmt.Printf("%s关闭连接;socketId:%d \n", s.ip, s.GetId())
	return true
}

// heart
func (s *ServerSocketClient) Update() {
	now := int(time.Now().Unix())
	timer.StopTimer(s.timerId)
	// timeout
	if s.heartTime < now {
		s.OnNetFail(2)
		return
	}
}

func (s *ServerSocketClient) SendLoop() bool {
	for {
		defer func() {
			if err := recover(); err != nil {
				base.TraceCode(err)
			}
		}()

		select {
		case buff := <-s.sendChan:
			if buff == nil { //信道关闭
				return false
			} else {
				s.DoSend(buff)
			}
		}
	}

	return true
}

func (s *ServerSocketClient) SendPacket(head rpc.RpcHead, funcName string, msg proto.Message) int {
	rpcPacketByte, _ := rpc.MarshalPacket(head, funcName, msg)
	var packet = rpc.Packet{
		Buff: rpcPacketByte,
	}
	return s.Send(head, packet)
}

func (s *ServerSocketClient) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) int {
	buff := rpc.Marshal(&head, &funcName, params...)
	return s.Send(head, buff)
}

// 设置链接属性
func (s *ServerSocketClient) SetProperty(p any) {
	s.Property = p
}

// 获取链接属性
func (s *ServerSocketClient) GetProperty() (p any) {
	if s.Property == nil {
		return nil
	}
	return s.Property
}

// 移除链接属性
func (s *ServerSocketClient) RemoveProperty() {
	s.Property = nil
}
