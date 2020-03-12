package skip

import pb "rxt/cmd/barrier/proto/sc"

type ScSkip struct {
	IScSkip
}

func New() ScSkip {
	scSkip := &ScSkipV1{}
	scSkip.Init()
	return ScSkip{scSkip}
}

type IScSkip interface {
	Skip(request *pb.Request) (bool, error)
}
