package service

import (
	"context"
	"sync"
)

type Inserter interface {
	Push(q, params string)
	Shutdown(ctx context.Context)
	SetWg(wg *sync.WaitGroup)
}

type Updater interface {
	Push(q string)
	Shutdown(ctx context.Context)
	SetWg(wg *sync.WaitGroup)
}

type Parser interface {
	Parse(body string) (params, content string, insert, update bool, err error)
}
