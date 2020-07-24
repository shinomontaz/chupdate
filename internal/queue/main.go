package queue

import (
	"context"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Queue struct {
	Query       string
	Rows        []string
	MaxCount    int
	MaxInterval int
	mu          sync.Mutex
	makeReq     func(q, content string, count int)
}

func Create(Count, FlushInterval int, query string, makeReq func(q, content string, count int)) *Queue {
	log.Debug("creating queue", query)
	q := &Queue{
		Query:       query,
		MaxCount:    Count,
		MaxInterval: FlushInterval,
		makeReq:     makeReq,
	}
	q.RunTimer()
	return q
}

func (q *Queue) RunTimer() {
	if q.MaxInterval < 0 {
		return
	}

	ticker := time.NewTicker(time.Millisecond * time.Duration(q.MaxInterval))
	go func() {
		for range ticker.C {
			q.Flush()
		}
	}()
}

func (q *Queue) Add(text string) {
	log.Debug("queue add", text)
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Rows = append(q.Rows, text)
	if len(q.Rows) >= q.MaxCount {
		q.flush()
	}
}

func (q *Queue) Content() string {
	rowDelimiter := "\n"
	return q.Query + "\n" + strings.Join(q.Rows, rowDelimiter)
}

func (q *Queue) flush() {
	q.makeReq(q.Query, q.Content(), len(q.Rows)) // тут и запрос подготовим, тут и в канал запишем
	q.Rows = make([]string, 0, q.MaxCount)
}

func (q *Queue) Flush() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	ln := len(q.Rows)
	if ln > 0 {
		q.flush()
	}
	return ln
}

func (q *Queue) Shutdown(ctx context.Context) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.flush()
}
