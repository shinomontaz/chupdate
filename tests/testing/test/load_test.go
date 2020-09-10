package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"math/rand"
	"strings"
	"testing"
)

var in chan string
var client *http.Client
var chUrl string

var wg sync.WaitGroup

func init() {
	in = make(chan string, 100000000)
	//	go ListenResults(in)

	client = &http.Client{
		Timeout: time.Second * time.Duration(10),
	}

	chUrl = "http://localhost:8123?user=default&password=qwe&database=test"
}

func ListenResults(in chan string) {
	// for res := range in {
	// 	//		fmt.Println(res)
	// }
}

func BenchmarkSelectHttp(b *testing.B) {
	id := rand.Intn(100)
	query := fmt.Sprintf("SELECT * FROM test.test2 WHERE id = %d ORDER BY time LIMIT 1 BY id", uint32(id))
	err := Process(query, in)
	if err != nil {
		panic(err)
	}
}

func BenchmarkSelectHttpParallel(b *testing.B) {
	for i := 0; i < 100; i++ {
		id := rand.Intn(100)
		query := fmt.Sprintf("SELECT * FROM test.test2 WHERE id = %d ORDER BY time LIMIT 1 BY id", uint32(id))
		wg.Add(1)
		go func() {
			err := Process(query, in)
			defer wg.Done()

			if err != nil {
				panic(err)
			}

		}()
	}
	wg.Wait()
}

func BenchmarkParallelSelect(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := rand.Intn(100)
			query := fmt.Sprintf("SELECT * FROM test.test2 WHERE id = %d ORDER BY time LIMIT 1 BY id", uint32(id))
			err := Process(query, in)
			if err != nil {
				panic(err)
			}

		}
	})
}

func Process(query string, out chan string) error {
	resp, err := client.Post(chUrl, "", strings.NewReader(query))

	if err != nil {
		return errors.New("server is down")
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	ss := string(buf)
	out <- ss

	if resp.StatusCode >= 502 {
		err = errors.New("server is down")
	} else if resp.StatusCode >= 400 {
		err = fmt.Errorf("Wrong server status %+v:\nresponse: %+v\nrequest: %#v", resp.StatusCode, ss, query)
	}

	return nil
}

func randBool() uint8 {
	if rand.Float32() < 0.5 {
		return 1
	}

	return 0
}

func randString(length int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
