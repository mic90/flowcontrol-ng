package graph

import (
	"context"
	"github.com/rs/xid"
)

type Runner interface {
	Run(ctx context.Context, nodes map[xid.ID]Node)
	Stop()
	State() State
}
