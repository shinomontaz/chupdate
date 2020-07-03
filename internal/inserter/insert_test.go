package parser

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/shinomontaz/chupdate/internal/errorer"
	"github.com/shinomontaz/chupdate/internal/inserter"
	"github.com/shinomontaz/chupdate/internal/parser"
	"github.com/stretchr/testify/assert"
)

const qTitle = "INSERT INTO table3 (c1, c2, c3) FORMAT TabSeparated"
const qContent = "v11	v12	v13\nv21	v22	v23"
const qValuesTitle = "INSERT INTO table3 (c1, c2, c3) Values"
const qValuesTitleUpper = "INSERT INTO table3 (c1, c2, c3) VALUES"
const qValuesContent = "(v11,v12,v13),(v21,v22,v23)"
const qSelect = "SELECT 1"
const qParams = "user=user&password=111"
const qSelectAndParams = "query=" + qSelect + "&" + qParams

const qFormatInQuotesQuery = "INSERT INTO test (date, args) VALUES"
const qFormatInQuotesValues = "('2019-06-13', 'query=select%20args%20from%20test%20group%20by%20date%20FORMAT%20JSON')"

const qTSNamesTitle = "INSERT INTO table3 (c1, c2, c3) FORMAT TabSeparatedWithNames"
const qNames = "field1	field2	field3"

const qTitle2 = "INSERT INTO table3 (c1, c2, c3) VALUES"
const qContent2 = "(1, 2, 3)"

var escTitle = url.QueryEscape(strings.ToLower(qTitle))
var escSelect = url.QueryEscape(strings.ToLower(qSelect))
var escParamsAndSelect = qParams + "&query=" + escSelect

func TestInsert(t *testing.T) {
	var makeReq func(q, content string, count int)
	makeReq = func(q, content string, count int) {
		fmt.Println("now we send a request", q, content)
	}

	go errorer.Listen(errors)
	parsr := parser.New()
	instr := inserter.New(env.Config.FlushInterval, env.Config.FlushCount, makeReq, errors)
}

func TestUpdate(t *testing.T) {
	prs := New()
	qTitleUpdate := "UPDATE table3 SET c1 = 1, 'c2' = 'sdfsdf' WHERE c3 = 2 AND c4 = 'sdfsdf'"
	body := strings.ToLower(qTitleUpdate)
	table, where, cols, vals := prs.Updateparse(body)

	assert.Equal(t, "table3", table)
	assert.Equal(t, "c3 = 2 and c4 = 'sdfsdf'", where)

	t.Logf("where: %s", where)
	t.Logf("cols: %v", cols)
	t.Logf("vals: %v", vals)
}
