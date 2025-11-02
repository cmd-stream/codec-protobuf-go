package internal

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
