package gid

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	log "github.com/cihub/seelog"
	"strconv"
	"strings"
	"time"
	"errors"
)

var client = &http.Client{}

func init() {
	client.Timeout = time.Second
}

type Server struct {
	UrlPrefix string
}

type Result struct {
	Id        int64 `json:"id"`
	MachineId int64 `json:"machine-id"`
	Msb       int64 `json:"msb"`
	Sequence  int64 `json:"sequence"`
	Timestamp int64 `json:"time"`
}

func NewServer(urlPrefix string) (*Server) {
	if urlPrefix == "" {
		urlPrefix = "http://sonyflake.live.xunlei.com/"
	}
	var S = &Server{
		UrlPrefix: urlPrefix,
	}
	return S
}
func (s *Server) GetId() (int64, error) {
	result, err := s.Get()
	if err != nil {
		return -1, err
	}
	return result.Id, nil
}

func (s *Server) Get() (*Result, error) {
	req, err := http.NewRequest("POST", s.UrlPrefix, strings.NewReader(""))
	if err != nil {
		log.Errorf("getId request error %v", err)
		return nil, err
	}
	resp, err := client.Do(req)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if status := resp.StatusCode; status < 200 || status >= 300 {
		log.Warnf("status code is not 200")
		return nil, errors.New("status code is " + strconv.Itoa(status))
	}
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		ret := &Result{}
		if err = json.Unmarshal([]byte(body), ret); err != nil {
			return nil, err
		}
		return ret, nil
	}
}
