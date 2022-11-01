package raft

import (
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/erDong01/micro-kit/common"

	"github.com/erDong01/micro-kit/base"
	"github.com/hashicorp/raft"
	raftBoltdb "github.com/hashicorp/raft-boltdb/v2"
)

type (
	Raft struct {
		*raft.Raft
		*common.ClusterInfo
		hashRing       *base.HashRing //hash一致性
		clusterInfoMap map[uint32]*common.ClusterInfo
	}
)

func (this *Raft) InitRaft(info *common.ClusterInfo, Endpoints []string, fsm raft.FSM) {
	this.ClusterInfo = info
	this.hashRing = base.NewHashRing()
	this.clusterInfoMap = make(map[uint32]*common.ClusterInfo)

	this.Raft, _ = NewRaft(info.IpString(), info.IpString(), "./node", fsm)
	var configuration raft.Configuration
	for _, v := range Endpoints {
		server := raft.Server{ID: raft.ServerID(v), Address: raft.ServerAddress(v)}
		configuration.Servers = append(configuration.Servers, server)
	}
	this.BootstrapCluster(configuration)
}
func (this *Raft) IsLeader() bool {
	return string(this.Leader()) == this.IpString()
}

func (this *Raft) GetHashRing(Id int64) (error, uint32) {
	return this.hashRing.Get64(Id)
}

func NewRaft(raftAddr, raftId, raftDir string, fsm raft.FSM) (*raft.Raft, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(raftId)
	addr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(raftAddr, addr, 3, 5*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	snapshots, err := raft.NewFileSnapshotStore(raftDir, 1, os.Stderr)
	if err != nil {
		return nil, err
	}

	logStore, err := raftBoltdb.NewBoltStore(filepath.Join(raftDir, "raft-log.db"))
	if err != nil {
		return nil, err
	}
	stableStore, err := raftBoltdb.NewBoltStore(filepath.Join(raftDir, "raft-stable.db"))
	if err != nil {
		return nil, err
	}

	rf, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return nil, err
	}
	return rf, nil
}
