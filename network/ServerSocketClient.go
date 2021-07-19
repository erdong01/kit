package network

import (
	"fmt"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools"
	"github.com/erDong01/micro-kit/wrong"
	"google.golang.org/protobuf/proto"
	"hash/crc32"
	"io"
	"log"
	"net"
)

var (
	DISCONNECTINT = crc32.ChecksumIEEE([]byte("DISCONNECT"))
)

type IServerSocketClient interface {
	ISocket
}

type ServerSocketClient struct {
	Socket
	ServerSocket *ServerSocket
	sendChan     chan []byte //对外缓冲队列
}

func (this *ServerSocketClient) Init(ip string, port int) bool {
	this.sendChan = make(chan []byte, MAX_SEND_CHAN)
	this.Socket.Init(ip, port)
	return true
}

func (this *ServerSocketClient) Start() bool {
	if this.state != SSF_SHUT_DOWN {
		return false
	}
	if this.ServerSocket == nil {
		return false
	}
	if this.PacketFuncList.Len() == 0 {
		this.PacketFuncList = this.ServerSocket.PacketFuncList
	}
	this.state = SSF_CONNECT
	this.Conn.(*net.TCPConn).SetNoDelay(true)
	go this.Run()
	if this.connectType == CLIENT_CONNECT {
		go this.SendLoop()
	}
	return true
}

func (this *ServerSocketClient) Run() bool {
	var buff = make([]byte, this.ReceiveBufferSize)
	loop := func() bool {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		if this.shuttingDown || this.Conn == nil {
			return false
		}
		n, err := this.Conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("远程链接：%s已经关闭！\n", this.Conn.RemoteAddr().String())
			return false
		}
		if err != nil {
			handleError(err)
			this.OnNetFail()
			return false
		}
		if n > 0 {
			if !this.packetParser.Read(buff[:n]) && this.connectType == CLIENT_CONNECT {
				this.OnNetFail()
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
	this.Close()
	fmt.Printf("%s关闭连接 \n", this.IP)
	return true
}

func (this *ServerSocketClient) SendLoop() bool {
	for {
		select {
		case buff := <-this.sendChan:
			if buff == nil {
				return false
			} else {
				this.DoSend(buff)
			}
		}
	}
	return true
}
func (this *ServerSocketClient) DoSend(buff []byte) int {
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

func (this *ServerSocketClient) SendPacket(head rpc3.RpcHead, funcName string, packet proto.Message) {
	buff := rpc.MarshalPacket(head, funcName, packet)
	this.Send(rpc3.RpcHead{}, buff)
}

func (this *ServerSocketClient) SendMsg(head rpc3.RpcHead, funcName string, params ...interface{}) int {
	buff := rpc.Marshal(head, funcName, params...)
	return this.Send(head, buff)
}

func (this *ServerSocketClient) Send(head rpc3.RpcHead, buff []byte) int {
	defer func() {
		if err := recover(); err != nil {
			wrong.TraceCode(err)
		}
	}()

	if this.connectType == CLIENT_CONNECT { //对外链接send不阻塞
		select {
		case this.sendChan <- buff:
		default: //网络太卡,tcp send缓存满了并且发送队列也满了
			this.OnNetFail()
		}
	} else {
		return this.DoSend(buff)

	}
	return 0
}

func (this *ServerSocketClient) Close() {
	this.Socket.Close()
	if this.ServerSocket != nil {
		this.ServerSocket.DelClient(this)
	}
}

func (this *ServerSocketClient) OnNetFail() {
	this.Stop()
	if this.connectType == CLIENT_CONNECT {
		stream := tools.NewBitStream(make([]byte, 32), 32)
		stream.WriteInt(int(DISCONNECTINT), 32)
		stream.WriteInt(int(this.clientId), 32)
		this.HandlePacket(stream.GetBuffer())
	} else {
		this.CallMsg("DISCONNECT", this.clientId)
	}
	if this.ServerSocket != nil {
		this.ServerSocket.DelClient(this)
	}
}
