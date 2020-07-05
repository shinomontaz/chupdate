package queue

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"testing"

	"github.com/shinomontaz/chupdate/config"
	"github.com/shinomontaz/chupdate/internal/parser"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var env *config.Env

func init() {
	env = config.NewEnv("../../config")
	env.Config.TestFlag = false
	env.InitLog()
}

func TestAdd(t *testing.T) {
	parsr := parser.New()
	var q *Queue

	var makeReq func(q, content string, count int)
	makeReq = func(q, content string, count int) {
		log.Debug("now we send a request", q, content)
	}

	var data, randName, params string
	var randField uint8

	var query = "INSERT INTO test2 (id, event, another_field) VALUES"
	content := "query=" + url.QueryEscape(query) + "\n"
	var params_arr []string
	for i := 0; i < 100; i++ {
		randName = randString(10)
		randField = randBool()
		params = fmt.Sprintf("(%d, '%s', %d)", uint32(i+1), randName, randField)
		data = fmt.Sprintf("%s %s", query, params)

		params_arr = append(params_arr, params)
		query, content, _, _, err := parsr.Parse(data) // query=INSERT INTO table3 (c1, c2, c3) VALUES ('v1', 'v2', 'v3') and this string entirely is in 'qs' or 'ss'
		if i == 0 {
			q = Create(1000, -1, query, makeReq)
		}

		if err != nil {
			panic(err)
		}
		q.Add(content)
	}

	assert.Equal(t, 100, len(q.Rows))
	assert.Equal(t, content+strings.Join(params_arr, "\n"), q.Content())
}

func TestCount(t *testing.T) {
	parsr := parser.New()
	var q *Queue

	var makeReq func(q, content string, count int)
	makeReq = func(q, content string, count int) {
		log.Debug("now we send a request", q, content)
	}

	var data, randName, params string
	var randField uint8

	var query = "INSERT INTO test2 (id, event, another_field) VALUES"
	var params_arr []string

	for i := 0; i < 100; i++ {
		randName = randString(10)
		randField = randBool()
		params = fmt.Sprintf("(%d, '%s', %d)", uint32(i+1), randName, randField)
		data = fmt.Sprintf("%s %s", query, params)

		params_arr = append(params_arr, params)
		query, content, _, _, err := parsr.Parse(data) // query=INSERT INTO table3 (c1, c2, c3) VALUES ('v1', 'v2', 'v3') and this string entirely is in 'qs' or 'ss'
		if i == 0 {
			q = Create(101, -1, query, makeReq)
		}

		if err != nil {
			panic(err)
		}
		q.Add(content)
	}

	ln := q.Flush()

	assert.Equal(t, 100, ln)
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
