package network

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

type ServerSocketClient struct {
	Socket
	ServerSocket *ServerSocket
	SendChan     chan []byte //对外缓冲队列
}

func (this *ServerSocketClient) Init(ip string, port int) bool {
	this.SendChan = make(chan []byte, MAX_SEND_CHAN)
	this.Socket.Init(ip, port)
	return true
}

func (this *ServerSocketClient) Start() bool {
	if this.PacketFuncList.Len() == 0 {
		this.PacketFuncList = this.ServerSocket.PacketFuncList
	}
	this.Conn.(*net.TCPConn).SetNoDelay(true)

	go this.Run()
	go this.SendLoop()
	return true
}

func (this *ServerSocketClient) Run() bool {
	var buff = make([]byte, this.ReceiveBufferSize)
	loop := func() bool {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()
		n, err := this.Conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("远程链接：%s已经关闭！\n", this.Conn.RemoteAddr().String())
			return false
		}
		if n > 0 {
			if !this.packetParser.Read(buff[:n]) {
				return false
			}
		}

		return true
	}
	for {
		if !loop() {
			break
		}
	}
	this.ServerSocket.TCPListener.Close()
	fmt.Printf("%s关闭连接", this.IP)
	return true
}

func (this *ServerSocketClient) SendLoop() bool {
	for {
		select {
		case buff := <-this.SendChan:
			if buff == nil {
				return false
			} else {
				this.Send(buff)
			}
		}
	}
	return true
}
func (this *ServerSocketClient) Send(buff []byte) int {
	if this.Conn == nil {
		return 0
	}
	n, err := this.Conn.Write(this.packetParser.Write(buff))
	if err != nil {
		log.Error(err)
	}
	if n > 0 {
		return n
	}
	return 0
}

func (this *ServerSocketClient) Close() {
	this.Conn.Close()
	this.ServerSocket.DelClient(this)
}
