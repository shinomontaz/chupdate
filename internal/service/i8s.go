package service

import (
	"context"
	"sync"

	"github.com/shinomontaz/chupdate/internal/types"
)

type Inserter interface {
	Push(pq *types.ParsedQuery)
	Shutdown(ctx context.Context)
	SetWg(wg *sync.WaitGroup)
}

type Updater interface {
	Push(pq *types.ParsedQuery)
	Shutdown(ctx context.Context)
	SetWg(wg *sync.WaitGroup)
}

type Parser interface {
	Parse(body string) *types.ParsedQuery
}
