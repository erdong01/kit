package network

import (
	"fmt"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/wrong"
	"io"
	"log"
	"net"
)

type (
	IClientSocket interface {
		ISocket
	}

	ClientSocket struct {
		Socket
		maxClients int
		minClients int
	}
)

func handleError(err error) {
	if err == nil {
		return
	}
	log.Printf("错误：%s\n", err.Error())
}

func (this *ClientSocket) Init(ip string, port int) bool {
	if this.Port == port || this.IP == ip {
		return false
	}
	this.Socket.Init(ip, port)
	this.IP = ip
	this.Port = port
	return true
}

func (this *ClientSocket) Start() bool {

	this.shuttingDown = false
	if this.IP == "" {
		this.IP = "127.0.0.1"
	}
	if this.Connect() {
		this.Conn.(*net.TCPConn).SetNoDelay(true)
		go this.Run()
	}
	return true
}

func (this *ClientSocket) Stop() bool {
	if this.shuttingDown {
		return true
	}
	this.shuttingDown = true
	this.Close()
	return true
}

func (this *ClientSocket) SendMsg(head rpc3.RpcHead, funcName string, params ...interface{}) {
	buff := rpc.Marshal(head, funcName, params...)
	this.Send(head, buff)
}

func (this *ClientSocket) Send(head rpc3.RpcHead, buff []byte) int {
	defer func() {
		if err := recover(); err != nil {
			wrong.TraceCode(err)
		}
	}()
	if this.Conn == nil {
		return 0
	}
	n, err := this.Conn.Write(this.packetParser.Write(buff))
	handleError(err)
	if n > 0 {
		return n
	}
	return 0
}

func (this *ClientSocket) Restart() bool {
	return true
}
func (this *ClientSocket) Connect() bool {
	if this.state == SSF_CONNECT {
		return false
	}
	var strRemote = fmt.Sprintf("%s:%d", this.IP, this.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", strRemote)
	if err != nil {
		log.Println("%v", err)
	}
	ln, err1 := net.DialTCP("tcp", nil, tcpAddr)
	if err1 != nil {
		return false
	}
	this.state = SSF_CONNECT
	this.SetTcpConn(ln)
	this.CallMsg("COMMON_RegisterRequest")
	return true
}
func (this *ClientSocket) OnNetFail(int) {
	this.Stop()
	this.CallMsg("DISCONNECT", this.clientId)
}

func (this *ClientSocket) Run() bool {
	var buff = make([]byte, this.ReceiveBufferSize)
	loop := func() bool {
		if err := recover(); err != nil {
			wrong.TraceCode(err)
		}

		if this.shuttingDown || this.Conn == nil {
			return false
		}

		n, err := this.Conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("远程链接：%s已经关闭！\n", this.Conn.RemoteAddr().String())
			this.OnNetFail(0)
			return false
		}
		if err != nil {
			handleError(err)
			this.OnNetFail(0)
		}
		if n > 0 {
			this.packetParser.Read(buff[:n])
		}
		return true
	}
	for {
		if !loop() {
			break
		}
	}
	this.Close()
	return true
}
