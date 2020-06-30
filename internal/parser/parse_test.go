package parser

import (
	"net/url"
	"strings"
	"testing"

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

func TestParse(t *testing.T) {
	prs := New()
	var params string
	var content string
	var insert bool
	var update bool

	body := strings.ToLower(qTitle2) + " " + strings.ToLower(qContent2)
	params, content, insert, update, _ = prs.Parse(body)

	escTitle2 := url.QueryEscape(strings.ToLower(qTitle2))

	assert.Equal(t, "query="+escTitle2, params)
	assert.Equal(t, strings.ToLower(qContent2), content)
	assert.Equal(t, false, update)
	assert.Equal(t, true, insert)

	params, content, insert, _, _ = prs.Parse(strings.ToLower(qTitle) + " " + qContent)

	assert.Equal(t, "query="+escTitle, params)
	assert.Equal(t, qContent, content)
	assert.Equal(t, true, insert)

	params, content, insert, _, _ = prs.Parse(strings.ToLower(qTitle) + " " + qContent)

	assert.Equal(t, "query="+escTitle, params)
	assert.Equal(t, qContent, content)
	assert.Equal(t, true, insert)

	params, content, insert, _, _ = prs.Parse(strings.ToLower(qValuesTitle) + " " + qValuesContent)

	assert.Equal(t, "query="+url.QueryEscape(strings.ToLower(qValuesTitle)), params)
	assert.Equal(t, qValuesContent, content)
	assert.Equal(t, true, insert)

	params, content, insert, _, _ = prs.Parse(strings.ToLower(qSelect))

	assert.Equal(t, "query="+escSelect, params)
	assert.Equal(t, "", content)
	assert.Equal(t, false, insert)

	params, content, insert, _, _ = prs.Parse(strings.ToLower(qTitle) + " " + qContent)

	assert.Equal(t, "query="+strings.ToLower(escTitle), strings.ToLower(params))
	assert.Equal(t, qContent, content)
	assert.Equal(t, true, insert)

	params, content, insert, _, _ = prs.Parse(strings.ToLower(qValuesTitle) + " " + qValuesContent)

	assert.Equal(t, "query="+strings.ToLower(url.QueryEscape(qValuesTitle)), strings.ToLower(params))
	assert.Equal(t, qValuesContent, content)
	assert.Equal(t, true, insert)
}

func TestUpdate(t *testing.T) {
	prs := New()
	qTitleUpdate := "UPDATE table3 SET c1 = 1, 'c2' = 'sdfsdf' WHERE c3 = 2 AND c4 = 'sdfsdf'"
	body := strings.ToLower(qTitleUpdate)
	where, cols, vals := prs.Updateparse(body)

	t.Logf("where: %s", where)
	t.Logf("cols: %v", cols)
	t.Logf("vals: %v", vals)

}
