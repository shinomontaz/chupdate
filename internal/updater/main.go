package updater

import "github.com/jmoiron/sqlx"

// Collector - query collector
type Service struct {
	db    *sqlx.DB
	instr Inserter
}

// NewCollector - default collector constructor
func New(instr Inserter, db *sqlx.DB) (c *Service) {
	return &Service{
		instr: instr,
		db:    db,
	}
}

// Push - adding query to collector with query params (with query) and rows
func (s *Service) Push(query string) {
	// parse query
	// get table, get PK
	// perform select and get all the row
	// construct insert with all the data
	// push insert query into inserter
}
