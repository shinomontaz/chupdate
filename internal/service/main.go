package service

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	Debug    bool
	inserter Inserter
	updater  Updater
	url      string
	db       *sqlx.DB
	client   *http.Client
	parser   Parser
	errs     chan<- error
}

func New(ins Inserter, upd Updater, parsr Parser, url string, db *sqlx.DB, errs chan<- error) *Service {
	return &Service{
		inserter: ins,
		updater:  upd,
		parser:   parsr,
		db:       db,
		url:      url,
		errs:     errs,
		Debug:    true,
		client: &http.Client{
			Timeout: time.Second * time.Duration(3),
		},
	}
}

func (s *Service) Process(data string) (resp string, err error) {
	// params = "INSERT INTO table3 (c1, c2, c3) FORMAT TabSeparated"
	// content = "v11	v12	v13\nv21	v22	v23"
	query, content, insert, update, err := s.parser.Parse(data) // query=INSERT INTO table3 (c1, c2, c3) VALUES ('v1', 'v2', 'v3') and this string entirely is in 'qs' or 'ss'
	if err != nil {
		return "", err
	}

	if insert {
		go s.inserter.Push(query, content)
		return "", nil
	}

	if update {
		go s.updater.Push(query)
		return "", nil
	}

	resp, _, err = s.SendQuery(query)
	return resp, err
}

func (s *Service) Handle(w http.ResponseWriter, r *http.Request) {
	// parse query
	// return result
	q, _ := ioutil.ReadAll(r.Body)
	ss := string(q)

	query, content, insert, update, err := s.parser.Parse(ss) // query=INSERT INTO table3 (c1, c2, c3) VALUES ('v1', 'v2', 'v3') and this string entirely is in 'qs' or 'ss'
	if err != nil {
		s.errs <- err
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	if insert {
		go s.inserter.Push(query, content)
		w.WriteHeader(http.StatusOK)
		return
	}

	if update {
		go s.updater.Push(query)
		w.WriteHeader(http.StatusOK)
		return
	}

	resp, _, err := s.SendQuery(ss)
	//	return resp, err

	if s.Debug {
		log.Printf("query %+v %+v\n", r.URL.String(), ss)
	}

	//	resp, err := s.Process(ss)
	// if err != nil {
	// 	s.errs <- err
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(fmt.Sprintf("%s", err)))
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(resp))

}

func (s *Service) Shutdown(ctx context.Context) error {
	log.Printf("service shutting down\n")
	return nil
}

func (s *Service) SendQuery(query string) (response string, status int, err error) {
	resp, err := s.client.Post(s.url, "", strings.NewReader(query))

	fmt.Println("send query", s.url)

	if err != nil {
		return err.Error(), http.StatusBadGateway, errors.New("server is down")
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	ss := string(buf)
	if resp.StatusCode >= 502 {
		err = errors.New("server is down")
	} else if resp.StatusCode >= 400 {
		err = fmt.Errorf("Wrong server status %+v:\nresponse: %+v\nrequest: %#v", resp.StatusCode, s, query)
	}
	return ss, resp.StatusCode, err
}
