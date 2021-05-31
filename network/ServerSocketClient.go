package network

type ServerSocketClient struct {
	Socket
	ServerSocket *ServerSocket
	SendChan     chan []byte //对外缓冲队列
	ClientId     uint32
}

func (this *ServerSocketClient) Run() {

}
