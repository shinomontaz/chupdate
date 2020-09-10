package updater

import (
	"context"
	"sync"

	"github.com/shinomontaz/chupdate/internal/types"
)

// Collector - query collector
type Service struct {
	instr Inserter
	parsr Parser
	errs  chan<- error
	//	chConn *clickhouse.Conn
	wg  *sync.WaitGroup
	sem chan struct{}
}

// NewCollector - default collector constructor
func New(instr Inserter, parsr Parser, chUrl string, errs chan<- error) (c *Service) {
	return &Service{
		instr: instr,
		parsr: parsr,
		errs:  errs,
		sem:   make(chan struct{}, 100),
	}
}

func (s *Service) getColumns(table string) []string {
	return []string{"id", "event", "another_field", "time"}
}

// Push - adding query to collector with query params (with query) and rows
func (s *Service) Push(pq *types.ParsedQuery) {

}

func (s *Service) SetWg(wg *sync.WaitGroup) {
	s.wg = wg
}

func (s *Service) Shutdown(ctx context.Context) {
	s.wg.Wait()
}
