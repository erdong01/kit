package network

type Socket struct {
	IP string
	Port int
	Zone string
}

const (
	CLIENT_CONNECT = iota //对外
	SERVER_CONNECT = iota //对内
)
