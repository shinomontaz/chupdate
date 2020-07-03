package updater

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Collector - query collector
type Service struct {
	db    *sqlx.DB
	instr Inserter
	parsr Parser
	errs  chan<- error
}

// NewCollector - default collector constructor
func New(instr Inserter, parsr Parser, db *sqlx.DB, errs chan<- error) (c *Service) {
	return &Service{
		instr: instr,
		parsr: parsr,
		db:    db,
		errs:  errs,
	}
}

// Push - adding query to collector with query params (with query) and rows
func (s *Service) Push(query string) {
	table, where, cols, values := s.parsr.Updateparse(query)

	// get table, get PK
	// perform select and get all the row

	selectsql := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, where)

	rows, err := s.db.Query(selectsql)
	if err != nil {
		s.errs <- err
		return
	}
	defer rows.Close()

	res_cols, err := rows.Columns()
	if err != nil {
		s.errs <- err
		return
	}
	var result [][]string
	pointers := make([]interface{}, len(res_cols))
	container := make([]string, len(res_cols))
	for i, _ := range pointers {
		pointers[i] = &container[i]
	}
	for rows.Next() {
		rows.Scan(pointers...)
		result = append(result, container)
	}

	if len(result) > 1 {
		s.errs <- errors.New("Malformed update request - non precise PK")
		return
	}

	res := result[0]
	founded := -1
	for idx, col := range res_cols {
		for i, update_col := range cols {
			if update_col == col {
				founded = i
				break
			}
		}
		if founded > 0 {
			res[idx] = values[founded]
		}

		founded = -1
	}

	insert_query := fmt.Sprintf("INSERT INTO %s (%s) VALUES", table, strings.Join(res_cols, ", "))
	params := fmt.Sprintf("(%s)", strings.Join(res, ", "))
	s.instr.Push(insert_query, params)
}
