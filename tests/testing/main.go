package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var client *http.Client
var chUrl string

var wg sync.WaitGroup

func main() {
	// создаем сервер, который слушает POST-сообщения и пробрасывает их на http-ручку CH и возвращает ответ.

	client = &http.Client{
		Timeout: time.Second * time.Duration(3),
	}

	chUrl = "http://localhost:8123?user=default&password=qwe&database=test"

	var query string

	for i := 0; i < 100; i++ {
		query = fmt.Sprintf("SELECT * FROM test.test2 WHERE id = %d ORDER BY time LIMIT 1 BY id", uint32(i+1))
		fmt.Println("send query: ", query)
		wg.Add(1)
		go func(query string) {
			resp, err := Process(query)
			defer wg.Done()
			if err != nil {
				panic(err)
			}
			fmt.Println(resp)
		}(query)
	}

	wg.Wait()
}

func Process(query string) (string, error) {
	resp, err := client.Post(chUrl, "", strings.NewReader(query))

	if err != nil {
		return "", errors.New("server is down")
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	ss := string(buf)
	if resp.StatusCode >= 502 {
		err = errors.New("server is down")
	} else if resp.StatusCode >= 400 {
		err = fmt.Errorf("Wrong server status %+v:\nresponse: %+v\nrequest: %#v", resp.StatusCode, ss, query)
	}

	return ss, nil
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
