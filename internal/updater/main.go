package updater

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
)

// Collector - query collector
type Service struct {
	db    *sqlx.DB
	instr Inserter
	parsr Parser
	errs  chan<- error
	wg    *sync.WaitGroup
	sem   chan struct{}
}

// NewCollector - default collector constructor
func New(instr Inserter, parsr Parser, chUrl string, db *sqlx.DB, errs chan<- error) (c *Service) {
	return &Service{
		instr: instr,
		parsr: parsr,
		db:    db,
		errs:  errs,
		sem:   make(chan struct{}, 100),
	}
}

func (s *Service) getColumns(table string) []string {
	return []string{"id", "event", "another_field", "time"}
}

// Push - adding query to collector with query params (with query) and rows
func (s *Service) Push(query string) {
	s.sem <- struct{}{}
	defer func() {
		<-s.sem
	}()
	defer s.wg.Done()
	table, where, cols, values, condition_cols := s.parsr.Updateparse(query)

	// get table, get PK
	// perform select and get all the row

	timecolumn := "time"
	selectsql := fmt.Sprintf("SELECT * FROM %s WHERE %s ORDER BY %s LIMIT 1 BY %s", table, where, timecolumn, strings.Join(condition_cols, ", "))

	//	selectsql := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, where)

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
	for i := range pointers {
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

	res := values

	if len(result) > 0 {
		res = result[0]
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
	} else {
		res_cols = cols
	}

	var res2 []string
	for i, v := range res {
		if res_cols[i] == timecolumn {
			res_cols = append(res_cols[:i], res_cols[i+1:]...)
			continue
		}
		res2 = append(res2, "'"+v+"'") // note the = instead of :=
	}

	insert_query := fmt.Sprintf("INSERT INTO %s (%s) VALUES", table, strings.Join(res_cols, ", "))
	params := fmt.Sprintf("(%s)", strings.Join(res2, ", "))

	s.wg.Add(1)
	s.instr.Push(insert_query, params)
}

func (s *Service) SetWg(wg *sync.WaitGroup) {
	s.wg = wg
}

func (s *Service) Shutdown(ctx context.Context) {
	s.wg.Wait()
}
