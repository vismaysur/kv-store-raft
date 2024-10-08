package raft

import (
	"encoding/gob"
	"errors"
	"net"
	"net/rpc"

	"github.com/vismaysur/kv-store-raft/internal/utils"
)

var (
	ErrIncorrectLeaderStr = "ErrIncorrectLeader"
	ErrDeadNodeStr        = "ErrDeadNode"
)

var (
	ErrIncorrectLeader = errors.New("ErrIncorrectLeader")
	ErrDeadNode        = errors.New("ErrDeadNode")
)

type RequestVoteRequest struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteResponse struct {
	Term        int
	VoteGranted bool
}

type AppendEntriesRequest struct {
	Term         int
	LeaderId     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []map[string]interface{}
	LeaderCommit int
}

type AppendEntriesResponse struct {
	Term    int
	Success bool
	Reason  string
}

type ClientReqRequest struct {
	Entries []map[string]interface{}
}

type ClientReqResponse struct {
	CommitIndex int
	Success     bool
	Err         error
}

func (rf *Raft) call(peer, rpcname string, req interface{}, res interface{}) bool {
	conn, err := net.Dial("tcp", peer)
	if err != nil {
		utils.Dprintf("Failed to dial %s: %s", peer, err)
		return false
	}

	defer conn.Close()

	var id int
	var address string

	rf.withLock("", func() {
		id = rf.me
		address = rf.peers[rf.me]
	})

	client := rpc.NewClient(conn)

	if err := client.Call(rpcname, req, res); err != nil {
		utils.Dprintf(
			"[%d @ %s] RPC call to [%s] failed: %s",
			id, address,
			peer,
			err,
		)
		return false
	}

	return true
}

func init() {
	gob.Register(errors.New(""))
}
