package graceful

import (
	"context"
)

type App interface {
	SetUp(ctx context.Context, makeGo MakeGoroutine) error
	TearDown()
}
