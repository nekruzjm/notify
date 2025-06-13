package ctxman

import "context"

const _ctxInfoKey = "ctxInfo"

type Info struct {
	ShardName string
	DBName    string
	RequestID string
	UserID    int
	IsReplica bool
}

func Get(ctx context.Context) Info {
	value := ctx.Value(_ctxInfoKey)
	info, _ := value.(Info)
	return info
}

func Save(ctx context.Context, info Info) context.Context {
	return context.WithValue(ctx, _ctxInfoKey, info)
}
