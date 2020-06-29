package parser

import (
	"net/url"
	"regexp"
	"strings"
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
func (p *Service) Parse(body string) (query, content string, insert, update bool, err error) {
	var q string
	q, content = p.dataparse(body)
	q = strings.TrimSpace(q)
	if strings.HasPrefix(q, "insert") {
		insert = true
	}
	if strings.HasPrefix(q, "update") {
		update = true
	}

	query = "query=" + url.QueryEscape(q)
	return strings.TrimSpace(query), strings.TrimSpace(content), insert, update, nil
}

func (p *Service) dataparse(text string) (prefix string, content string) {
	i := strings.Index(text, "FORMAT")
	k := strings.Index(text, "VALUES")
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
			prefix = text[:off]
			content = text[off:]
		}
	} else {
		if k >= 0 {
			prefix = strings.TrimSpace(text[:k+6])
			content = strings.TrimSpace(text[k+6:])
		} else {
			off := p.regexFormat.FindStringSubmatchIndex(text)
			if len(off) > 3 {
				prefix = text[:off[3]]
				content = text[off[3]:]
			} else {
				off := p.regexValues.FindStringSubmatchIndex(text)
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
