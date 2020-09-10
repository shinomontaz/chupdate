package service

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	Debug    bool
	inserter Inserter
	updater  Updater
	url      string
	client   *http.Client
	parser   Parser
	errs     chan<- error
	wg       sync.WaitGroup
}

func New(ins Inserter, upd Updater, parsr Parser, url string, errs chan<- error) *Service {
	s := &Service{
		inserter: ins,
		updater:  upd,
		parser:   parsr,
		url:      url,
		errs:     errs,
		Debug:    true,
		client: &http.Client{
			Timeout: time.Second * time.Duration(3),
		},
	}

	ins.SetWg(&s.wg)
	upd.SetWg(&s.wg)

	return s
}

func (s *Service) Process(data string) (resp string, err error) {
	// params = "INSERT INTO table3 (c1, c2, c3) FORMAT TabSeparated"
	// content = "v11	v12	v13\nv21	v22	v23"
	data = strings.ToLower(data)
	pq := s.parser.Parse(data)

	// if err != nil {
	// 	return "", err
	// }

	if pq.Insert {
		s.wg.Add(1)
		// add from "list" to inserter.ocache

		go s.inserter.Push(pq)
		return "", nil
	}

	if pq.Update {
		s.wg.Add(1)
		go s.updater.Push(pq)
		return "", nil
	}

	resp, _, err = s.SendQuery(data)
	return resp, err
}

func (s *Service) Handle(w http.ResponseWriter, r *http.Request) {
	// parse query
	// return result
	q, _ := ioutil.ReadAll(r.Body)
	ss := string(q)

	resp, err := s.Process(ss)
	if err != nil {
		s.errs <- err
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(resp))
}

func (s *Service) Shutdown(ctx context.Context) error {
	s.wg.Wait()
	s.updater.Shutdown(ctx)
	s.inserter.Shutdown(ctx)

	log.Printf("service shutting down\n")
	return nil
}

func (s *Service) SendQuery(query string) (response string, status int, err error) {
	resp, err := s.client.Post(s.url, "", strings.NewReader(query))

	log.Debug("send query", s.url)

	if err != nil {
		return err.Error(), http.StatusBadGateway, errors.New("server is down")
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	ss := string(buf)
	if resp.StatusCode >= 502 {
		err = errors.New("server is down")
	} else if resp.StatusCode >= 400 {
		err = fmt.Errorf("Wrong server status %+v:\nresponse: %+v\nrequest: %#v", resp.StatusCode, ss, query)
	}
	return ss, resp.StatusCode, err
}
