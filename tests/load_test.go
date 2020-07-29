package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"math/rand"
	"strings"
	"testing"

	"github.com/shinomontaz/chupdate/config"
	"github.com/shinomontaz/chupdate/internal/errorer"
	"github.com/shinomontaz/chupdate/internal/inserter"
	"github.com/shinomontaz/chupdate/internal/parser"

	log "github.com/sirupsen/logrus"

	"github.com/shinomontaz/chupdate/internal/service"
	"github.com/shinomontaz/chupdate/internal/updater"
)

var env *config.Env
var errors chan error

func init() {
	env = config.NewEnv("../config")
	env.Config.TestFlag = false
	env.InitLog()
	env.InitDb()
	errors = make(chan error, 1000)
}

func makeReq(q, content string, count int) {
	cl := &http.Client{
		Timeout: time.Second * time.Duration(30),
	}

	log.Debug("send query ", content)
	resp, err := cl.Post(env.Config.CHUrl, "", strings.NewReader(content))
	if err != nil {
		fmt.Println("panic!")
		panic(err)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	s := string(buf)
	if resp.StatusCode >= 502 {
		panic("502")
	} else if resp.StatusCode >= 400 {
		panic(fmt.Sprintf("Wrong server status %+v:\nresponse: %+v\nrequest: %#v", resp.StatusCode, s, content))
	}
}

func BenchmarkInsert(b *testing.B) {
	go errorer.Listen(errors)

	parsr := parser.New()
	instr := inserter.New(env.Config.FlushInterval, env.Config.FlushCount, makeReq, errors)
	updtr := updater.New(instr, parsr, env.Config.CHUrl, env.Db, errors)

	prc := service.New(instr, updtr, parsr, env.Config.CHUrl, env.Db, errors)
	var query, randName string
	var randField uint8
	for i := 0; i < b.N; i++ {
		randName = randString(10)
		randField = randBool()
		query = fmt.Sprintf("INSERT INTO test2 (id, event, another_field) VALUES (%d, '%s', %d)", uint32(i+1), randName, randField)

		_, err := prc.Process(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestInsert(t *testing.T) {
	go errorer.Listen(errors)

	parsr := parser.New()
	instr := inserter.New(-1, env.Config.FlushCount, makeReq, errors)
	updtr := updater.New(instr, parsr, env.Config.CHUrl, env.Db, errors)

	prc := service.New(instr, updtr, parsr, env.Config.CHUrl, env.Db, errors)
	var query, randName string
	var randField uint8

	start := time.Now()

	for i := 0; i < 10000; i++ {
		randName = randString(10)
		randField = randBool()
		query = fmt.Sprintf("INSERT INTO test2 (id, event, another_field) VALUES (%d, '%s', %d)", uint32(i+1), randName, randField)
		_, err := prc.Process(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx := context.Background()
	prc.Shutdown(ctx)

	elapsed := time.Since(start)

	fmt.Println(int64(elapsed))
}

func BenchmarkSelectHttp(b *testing.B) {
	go errorer.Listen(errors)

	parsr := parser.New()
	instr := inserter.New(env.Config.FlushInterval, env.Config.FlushCount, makeReq, errors)
	updtr := updater.New(instr, parsr, env.Config.CHUrl, env.Db, errors)

	prc := service.New(instr, updtr, parsr, env.Config.CHUrl, env.Db, errors)
	var query string

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		query = fmt.Sprintf("SELECT * FROM test.test2 WHERE id = %d ORDER BY time LIMIT 1 BY id", uint32(i+1))
		wg.Add(1)
		go func(query string) {
			defer wg.Done()
			resp, err := prc.Process(query)
			if err != nil {
				panic(err)
			}
			log.Debug(resp)
		}(query)
	}

	wg.Wait()
}

func BenchmarkSelectParallel(b *testing.B) {
	go errorer.Listen(errors)

	parsr := parser.New()
	instr := inserter.New(env.Config.FlushInterval, env.Config.FlushCount, makeReq, errors)
	updtr := updater.New(instr, parsr, env.Config.CHUrl, env.Db, errors)

	prc := service.New(instr, updtr, parsr, env.Config.CHUrl, env.Db, errors)
	var query string
	var i uint32

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			query = fmt.Sprintf("SELECT * FROM test.test2 WHERE id = %d ORDER BY time LIMIT 1 BY id", i+1)
			resp, err := prc.Process(query)
			if err != nil {
				panic(err)
			}
			log.Debug(resp)
			i++
		}
	})
}

func TestUpdate(t *testing.T) {
	go errorer.Listen(errors)

	parsr := parser.New()
	instr := inserter.New(env.Config.FlushInterval, env.Config.FlushCount, makeReq, errors)
	updtr := updater.New(instr, parsr, env.Config.CHUrl, env.Db, errors)

	prc := service.New(instr, updtr, parsr, env.Config.CHUrl, env.Db, errors)
	var query string
	var randField uint8

	start := time.Now()
	for i := 0; i < 10000; i++ {

		randField = randBool()
		query = fmt.Sprintf("UPDATE test.test2 SET another_field = %d WHERE id = %d", randField, uint32(i))
		resp, err := prc.Process(query)
		if err != nil {
			panic(err)
		}
		log.Debug(resp)
	}

	ctx := context.Background()
	prc.Shutdown(ctx)

	elapsed := time.Since(start)

	fmt.Println(int64(elapsed))
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
