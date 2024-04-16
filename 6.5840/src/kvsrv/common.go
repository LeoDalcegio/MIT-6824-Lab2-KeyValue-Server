package kvsrv

// Put or Append
type PutAppendArgs struct {
	Key   string
	Value string
	ClerkID int64
	Idx int
}

type PutAppendReply struct {
	Value string
}

type GetArgs struct {
	Key string
	ClerkID int64
	Idx int
}

type GetReply struct {
	Value string
}
