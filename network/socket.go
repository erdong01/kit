package network

type Socket struct {
	IP                string
	Port              int
	Zone              string
	ReceiveBufferSize int //单次接收缓存
}

func (this *Socket) Init(string, int) {
	this.ReceiveBufferSize = 1024
}
