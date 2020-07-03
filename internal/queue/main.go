package queue

import (
	"fmt"
	"strings"
	"sync"
	"time"
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
	ticker := time.NewTicker(time.Millisecond * time.Duration(q.MaxInterval))
	go func() {
		for range ticker.C {
			q.Flush()
		}
	}()
}

func (q *Queue) Add(text string) {
	fmt.Println("queue add: ", text)
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

func (q *Queue) Flush() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.Rows) > 0 {
		q.flush()
	}
}
