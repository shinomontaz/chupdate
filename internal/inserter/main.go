package inserter

import (
	"context"
	"sync"

	"github.com/shinomontaz/chupdate/internal/types"

	"github.com/shinomontaz/chupdate/internal/queue"
)

type Service struct {
	List          map[string]*queue.Queue // в го нет генериков, потому придется связать
	mu            sync.RWMutex
	Count         int
	FlushInterval int
	makeReq       func(q, content string, count int)
	errs          chan<- error
	wg            *sync.WaitGroup
	ocache        map[string]map[string]string
}

func New(flush_interval, flush_count int, makeReq func(q, content string, count int), errs chan<- error) *Service {
	return &Service{
		FlushInterval: flush_interval,
		Count:         flush_count,
		List:          make(map[string]*queue.Queue),
		makeReq:       makeReq,
		errs:          errs,
		ocache:        make(map[string]map[string]string, flush_count),
	}
}

func (s *Service) Push(pq *types.ParsedQuery) {
	// log.Debug("inserter push ", query, params)
	// s.mu.Lock()
	// defer s.mu.Unlock()
	// defer s.wg.Done()

	// q, ok := s.List[query]
	// if !ok {
	// 	q = queue.Create(s.Count, s.FlushInterval, query, s.makeReq)
	// 	s.List[query] = q
	// }

	// q.Add(params)
}

func (s *Service) SetWg(wg *sync.WaitGroup) {
	s.wg = wg
}

func (s *Service) Shutdown(ctx context.Context) {
	for _, q := range s.List {
		q.Shutdown(ctx)
	}
}
