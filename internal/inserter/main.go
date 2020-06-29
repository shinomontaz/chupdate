package inserter

import (
	"sync"

	"github.com/shinomontaz/chupdate/internal/queue"
)

type Service struct {
	List          map[string]*queue.Queue // в го нет генериков, потому придется связать
	mu            sync.RWMutex
	Count         int
	FlushInterval int
	makeReq       func(q, content string, count int)
}

func New(flush_interval, flush_count int, makeReq func(q, content string, count int)) *Service {
	return &Service{
		FlushInterval: flush_interval,
		Count:         flush_count,
		List:          make(map[string]*queue.Queue),
		makeReq:       makeReq,
	}
}

func (s *Service) Push(query, params string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.List[query]
	if !ok {
		q = queue.Create(s.Count, s.FlushInterval, query, s.makeReq)
	}
	q.Add(params)
}
