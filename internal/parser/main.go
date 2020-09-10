package parser

import (
	"regexp"
	"strings"

	"github.com/shinomontaz/chupdate/internal/types"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	regexFormat    *regexp.Regexp
	regexValues    *regexp.Regexp
	regexGetFormat *regexp.Regexp
}

func New() *Service {
	return &Service{
		regexFormat:    regexp.MustCompile("(?i)format\\s\\S+(\\s+)"),
		regexValues:    regexp.MustCompile("(?i)\\svalues\\s"),
		regexGetFormat: regexp.MustCompile("(?i)format\\s(\\S+)"),
	}
}

// parse returns string of params and total content
// func (p *Service) Parse(body string) (query, content string, list map[string]string, insert, update bool, err error) {
func (p *Service) Parse(body string) *types.ParsedQuery {
	result := &types.ParsedQuery{}
	q := strings.TrimSpace(body)
	if strings.HasPrefix(q, "insert") {
		result.Insert = true
		p.dataparse(q, result)
	}
	if strings.HasPrefix(q, "update") {
		result.Update = true
		p.updateparse(q, result)
	}

	log.Debug("Parse: ", q)

	return result
}

func (p *Service) dataparse(text string, pq *types.ParsedQuery) {
	i := strings.Index(text, "FORMAT")
	k := strings.Index(text, "VALUES")
	l := strings.Index(text, "INTO")

	pq.Table = strings.TrimSpace(text[l:k]) // Criteria: tables with no quotes

	if k == -1 {
		k = strings.Index(text, "values")
	}
	if i >= 0 && i < k { // we have a FORMAT and VALUES keywords and VALUES after FORMAT
		w := false
		off := -1
		for c := i + 7; c < len(text); c++ { // start looking from end of keyword 'FORMAT'
			if !w && text[c] != ' ' && text[c] != '\n' && text[c] != ';' {
				w = true
			}
			if w && (text[c] == ' ' || text[c] == '\n' || text[c] == ';') {
				off = c + 1
				break
			}
		}
		if off >= 0 {
			pq.Query = text[:off]
			pq.InsertContent = text[off:]
		}
	} else {
		if k >= 0 {
			pq.Query = strings.TrimSpace(text[:k+6])
			pq.InsertContent = strings.TrimSpace(text[k+6:])
		} else {
			off := p.regexFormat.FindStringSubmatchIndex(text)
			if len(off) > 3 {
				pq.Query = text[:off[3]]
				pq.InsertContent = text[off[3]:]
			} else {
				off := p.regexValues.FindStringSubmatchIndex(text)
				if len(off) > 0 {
					pq.Query = text[:off[1]]
					pq.InsertContent = text[off[1]:]
				} else {
					pq.Query = text
				}
			}
		}
	}

	partCols := text[strings.Index(text, "(")+1 : strings.Index(text, ")")]
	cols := strings.Split(partCols, ",")

	log.Debugln("content: ", pq.InsertContent)

	vals := strings.Split(strings.Trim(pq.InsertContent, "()"), ",")
	log.Debugln(cols, vals)

	pq.UpdateMap = make(map[string]string, len(vals))
	for i, c := range cols {
		log.Debugln(i, "vals: ", vals[i])
		c = strings.TrimSpace(c)
		pq.UpdateMap[c] = strings.TrimSpace(vals[i])
	}
}

// (table, where string, cols []string, values []string, condition_cols []string)
func (p *Service) updateparse(text string, pq *types.ParsedQuery) {
	whIdx := strings.Index(text, "where")
	where := text[whIdx+6:]

	// parse where, get column names

	cond_cols := strings.Split(where, "AND")
	pq.Conditions = make([]string, len(cond_cols))
	for _, cond := range cond_cols {
		col_val := strings.Split(strings.TrimSpace(cond), "=")
		pq.Conditions = append(pq.Conditions, strings.TrimSpace(col_val[0]))
	}

	setIdx := strings.Index(text, "set")
	set := text[setIdx+3 : whIdx]

	//update
	pq.Table = strings.TrimSpace(text[6:setIdx])

	// split set substr

	pairs := strings.Split(set, ",")
	pq.UpdateMap = make(map[string]string, len(pairs))

	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		pq.UpdateMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
}
