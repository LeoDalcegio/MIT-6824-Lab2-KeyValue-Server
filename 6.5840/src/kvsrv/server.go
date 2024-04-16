package kvsrv

import (
	"log"
	"sync"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type Operation struct {
	idx int
	preAppendIdx int
}

type KVServer struct {
	mu sync.Mutex

	data map[string]string
	operations map[int64]Operation
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	reply.Value = kv.data[args.Key]
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	prev, ok := kv.operations[args.ClerkID]
	
	if ok && args.Idx == prev.idx {
		reply.Value = ""
	} else {
		reply.Value = kv.data[args.Key]
		kv.data[args.Key] = args.Value
		
		kv.operations[args.ClerkID] = Operation{
			idx: args.Idx,
			preAppendIdx: -1, // -1 because we are initializing the value here, no append happened before
		}
	}
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	prev, ok := kv.operations[args.ClerkID]

	if ok && args.Idx == prev.idx {
		reply.Value = kv.data[args.Key][:prev.preAppendIdx]
	} else {
		reply.Value = kv.data[args.Key]
		kv.data[args.Key] += args.Value
		
		preAppendIdx := len(reply.Value)

		kv.operations[args.ClerkID] = Operation{
			idx: args.Idx,
			preAppendIdx: preAppendIdx, // -1 because we are initializing the value here, no append happened before
		}
	}
}

func StartKVServer() *KVServer {
	kv := new(KVServer)
	kv.data = make(map[string]string)
	kv.operations = make(map[int64]Operation)
	return kv
}
