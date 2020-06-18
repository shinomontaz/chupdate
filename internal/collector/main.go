package collector

import (
	"net/url"
	"regexp"
	"strings"
	"sync"
)

const formatValues = "values"
const formatTabSeparated = "tabseparated"

var regexFormat = regexp.MustCompile("(?i)format\\s\\S+(\\s+)")
var regexValues = regexp.MustCompile("(?i)\\svalues\\s")
var regexGetFormat = regexp.MustCompile("(?i)format\\s(\\S+)")

// Collector - query collector
type Collector struct {
	Tables        map[string]*Table
	mu            sync.RWMutex
	Count         int
	FlushInterval int
}

// NewCollector - default collector constructor
func New(flush_interval, flush_count int) (c *Collector) {
	c = new(Collector)
	c.Tables = make(map[string]*Table)
	c.Count = flush_count
	c.FlushInterval = flush_interval
	return c
}

// Empty - check if all tables are empty
func (c *Collector) Empty() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, t := range c.Tables {
		if ok := t.Empty(); !ok {
			return false
		}
	}
	return true
}

// FlushAll - flush all tables to clickhouse
func (c *Collector) FlushAll() (count int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	count = 0
	for _, t := range c.Tables {
		if ok := t.CheckFlush(); ok {
			count++
		}

	}
	return count
}

// WaitFlush - wait for flush all tables
func (c *Collector) WaitFlush() (err error) {
	return nil
}

// AddTable - adding table to collector
func (c *Collector) AddTable(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.addTable(name)
}

func (c *Collector) separateQuery(name string) (query string, params string) {
	items := strings.Split(name, "&")
	for _, p := range items {
		if HasPrefix(p, "query=") {
			query = p[6:]
		} else {
			params += "&" + p
		}
	}
	if len(params) > 0 {
		params = strings.TrimSpace(params[1:])
	}
	q, err := url.QueryUnescape(query)
	if err != nil {
		return "", name
	}
	return q, params
}

func (c *Collector) getFormat(query string) (format string) {
	format = formatValues
	f := regexGetFormat.FindSubmatch([]byte(query))
	if len(f) > 1 {
		format = strings.TrimSpace(string(f[1]))
	}
	return format
}

func (c *Collector) addTable(name string) *Table {
	t := NewTable(name, c.Count, c.FlushInterval)
	query, params := c.separateQuery(name)
	t.Query = query
	t.Params = params
	t.Format = c.getFormat(query)
	c.Tables[name] = t
	t.RunTimer()
	return t
}

// Push - adding query to collector with query params (with query) and rows
func (c *Collector) Push(params string, content string) {
	c.mu.RLock()
	table, ok := c.Tables[params]
	if ok {
		table.Add(content)
		c.mu.RUnlock()
		//		pushCounter.Inc()
		return
	}
	c.mu.RUnlock()
	c.mu.Lock()
	table, ok = c.Tables[params]
	if !ok {
		table = c.addTable(params)
	}
	table.Add(content)
	c.mu.Unlock()
	//	pushCounter.Inc()
}

// ParseQuery - parsing inbound query to unified format (params/query), content (query data)
func (c *Collector) ParseQuery(queryString string, body string) (params string, content string, insert bool) {
	i := strings.Index(queryString, "query=")
	if i >= 0 {
		if HasPrefix(queryString[i+6:], "insert") {
			insert = true
		}
		var q string
		eoq := strings.Index(queryString[i+6:], "&")
		if eoq >= 0 {
			q = queryString[i+6 : eoq+6]
			params = queryString[:i] + queryString[eoq+7:]
		} else {
			q = queryString[i+6:]
			params = queryString[:i]
		}
		uq, err := url.QueryUnescape(q)
		if body != "" {
			uq += " " + body
		}
		if err != nil {
			return queryString, body, false
		}
		prefix, cnt := c.Parse(uq)
		if strings.HasSuffix(params, "&") || params == "" {
			params += "query=" + url.QueryEscape(strings.TrimSpace(prefix))
		} else {
			params += "&query=" + url.QueryEscape(strings.TrimSpace(prefix))
		}
		content = cnt
	} else {
		var q string
		q, content = c.Parse(body)
		q = strings.TrimSpace(q)
		if HasPrefix(q, "insert") {
			insert = true
		}
		if queryString != "" {
			params = queryString + "&query=" + url.QueryEscape(q)
		} else {
			params = "query=" + url.QueryEscape(q)
		}
	}
	return strings.TrimSpace(params), strings.TrimSpace(content), insert
}

// Parse - parsing text for query and data
func (c *Collector) Parse(text string) (prefix string, content string) {
	i := strings.Index(text, "FORMAT")
	k := strings.Index(text, "VALUES")
	if k == -1 {
		k = strings.Index(text, "values")
	}
	if i >= 0 && i < k {
		w := false
		off := -1
		for c := i + 7; c < len(text); c++ {
			if !w && text[c] != ' ' && text[c] != '\n' && text[c] != ';' {
				w = true
			}
			if w && (text[c] == ' ' || text[c] == '\n' || text[c] == ';') {
				off = c + 1
				break
			}
		}
		if off >= 0 {
			prefix = text[:off]
			content = text[off:]
		}
	} else {
		if k >= 0 {
			prefix = strings.TrimSpace(text[:k+6])
			content = strings.TrimSpace(text[k+6:])
		} else {
			off := regexFormat.FindStringSubmatchIndex(text)
			if len(off) > 3 {
				prefix = text[:off[3]]
				content = text[off[3]:]
			} else {
				off := regexValues.FindStringSubmatchIndex(text)
				if len(off) > 0 {
					prefix = text[:off[1]]
					content = text[off[1]:]
				} else {
					prefix = text
				}
			}
		}
	}
	return prefix, content
}

func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.ToLower(s[0:len(prefix)]) == strings.ToLower(prefix)
}
