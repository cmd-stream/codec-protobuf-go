package test

import (
	"context"

	"github.com/cmd-stream/cmd-stream-go/core"
)

func (c *Cmd1) Exec(ctx context.Context, receiver any, proxy core.Proxy) error {
	return nil
}

func (c *Cmd2) Exec(ctx context.Context, receiver any, proxy core.Proxy) error {
	return nil
}

func (r *Result1) LastOne() bool {
	return true
}

func (r *Result2) LastOne() bool {
	return true
}
