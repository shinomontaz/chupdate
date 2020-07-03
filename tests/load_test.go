package test

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"testing"

	"github.com/shinomontaz/chupdate/config"
	"github.com/shinomontaz/chupdate/internal/errorer"
	"github.com/shinomontaz/chupdate/internal/inserter"
	"github.com/shinomontaz/chupdate/internal/parser"

	"github.com/shinomontaz/chupdate/internal/service"
	"github.com/shinomontaz/chupdate/internal/updater"
)

var env *config.Env
var errors chan error

func init() {
	env = config.NewEnv("../config")
	env.InitDb()
	errors = make(chan error, 1000)
}

func BenchmarkInsert(b *testing.B) {
	go errorer.Listen(errors)

	var makeReq func(q, content string, count int)
	makeReq = func(q, content string, count int) {
		fmt.Println("now we send a request", q, content)
	}

	parsr := parser.New()
	instr := inserter.New(env.Config.FlushInterval, env.Config.FlushCount, makeReq, errors)
	updtr := updater.New(instr, parsr, env.Db, errors)

	prc := service.New(instr, updtr, parsr, env.Config.CHUrl, env.Db, errors)
	var query, randName string
	var randField uint8
	for i := 0; i < b.N; i++ {
		randName = randString(10)
		randField = randBool()
		query = fmt.Sprintf("INSERT INTO test2 (id, event, another_field) VALUES (%d, '%s', %d)", uint32(i+1), randName, randField)

		fmt.Println(query)

		_, err := prc.Process(query)
		if err != nil {
			log.Fatal(err)
		}
	}
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
