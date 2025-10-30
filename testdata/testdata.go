package testdata

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
)

func (c *Cmd1) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver any, proxy core.Proxy,
) error {
	return nil
}

func (c *Cmd2) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver any, proxy core.Proxy,
) error {
	return nil
}

func (r *Result1) LastOne() bool {
	return true
}

func (r *Result2) LastOne() bool {
	return true
}
