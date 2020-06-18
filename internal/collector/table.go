package collector

import (
	"strings"
	"sync"
	"time"
)

// Table - store query table info
type Table struct {
	Name          string
	Format        string
	Query         string
	Params        string
	Rows          []string
	count         int
	FlushCount    int
	FlushInterval int
	mu            sync.Mutex
	// todo add Last Error
}

// NewTable - default table constructor
func NewTable(name string, count int, interval int) (t *Table) {
	t = new(Table)
	t.Name = name
	t.FlushCount = count
	t.FlushInterval = interval
	return t
}

// Content - get text content of rowsfor query
func (t *Table) Content() string {
	rowDelimiter := "\n"
	return t.Query + "\n" + strings.Join(t.Rows, rowDelimiter)
}

// Flush - sends collected data in table to clickhouse
func (t *Table) Flush() {
	// req := ClickhouseRequest{
	// 	Params:  t.Params,
	// 	Query:   t.Query,
	// 	Content: t.Content(),
	// 	Count:   t.count,
	// }
	// t.Sender.Send(&req)
	t.Rows = make([]string, 0, t.FlushCount)
	t.count = 0
}

// CheckFlush - check if flush is need and sends data to clickhouse
func (t *Table) CheckFlush() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.count > 0 {
		t.Flush()
		return true
	}
	return false
}

// Empty - Checks if table is empty
func (t *Table) Empty() bool {
	return t.GetCount() == 0
}

// GetCount - Checks if table is empty
func (t *Table) GetCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.count
}

// RunTimer - timer for periodical savings data
func (t *Table) RunTimer() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(t.FlushInterval))
	go func() {
		for range ticker.C {
			t.CheckFlush()
		}
	}()
}

// Add - Adding query to table
func (t *Table) Add(text string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.count++
	t.Rows = append(t.Rows, text)
	if len(t.Rows) >= t.FlushCount {
		t.Flush()
	}
}
