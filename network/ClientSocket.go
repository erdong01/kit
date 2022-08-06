package network

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/wrong"
	"github.com/xtaci/kcp-go/v5"
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

func (this *ClientSocket) Init(ip string, port int, params ...OpOption) bool {
	if this.Port == port || this.IP == ip {
		return false
	}
	this.Socket.Init(ip, port)
	this.IP = ip
	this.Port = port
	fmt.Println(ip, port)
	return true
}

func (this *ClientSocket) Start() bool {
	if this.IP == "" {
		this.IP = "127.0.0.1"
	}
	if this.Connect() {
		go this.Run()
	}
	return true
}

func (this *ClientSocket) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) int {
	return this.Send(head, rpc.Marshal(head, funcName, params...))
}

func (this *ClientSocket) Send(head rpc.RpcHead, packet rpc.Packet) int {
	defer func() {
		if err := recover(); err != nil {
			wrong.TraceCode(err)
		}
	}()
	if this.Conn == nil {
		return 0
	}
	n, err := this.Conn.Write(this.packetParser.Write(packet.Buff))
	handleError(err)
	if n > 0 {
		return n
	}
	return 0
}

func (this *ClientSocket) DoSend(buff []byte) int {
	if this == nil || this.Conn == nil {
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
	var strRemote = fmt.Sprintf("%s:%d", this.IP, this.Port)
	connectStr := "Tcp"
	if this.bKcp {
		ln, err1 := kcp.Dial(strRemote)
		if err1 != nil {
			return false
		}
		this.SetConn(ln)
		connectStr = "Kcp"
	} else {
		tcpAddr, err := net.ResolveTCPAddr("tcp4", strRemote)
		if err != nil {
			log.Printf("%v", err)
		}
		ln, err1 := net.DialTCP("tcp4", nil, tcpAddr)
		if err1 != nil {
			return false
		}
		this.SetConn(ln)
	}
	fmt.Printf("%s 连接成功，请输入信息！\n", connectStr)
	this.CallMsg(rpc.RpcHead{}, "COMMON_RegisterRequest")
	return true
}
func (this *ClientSocket) OnNetFail(int) {
	this.Stop()
	this.CallMsg(rpc.RpcHead{}, "DISCONNECT", this.clientId)
}

func (this *ClientSocket) Run() bool {
	var buff = make([]byte, this.ReceiveBufferSize)
	loop := func() bool {
		defer func() {
			if err := recover(); err != nil {
				wrong.TraceCode(err)
			}
		}()

		if this.Conn == nil {
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
