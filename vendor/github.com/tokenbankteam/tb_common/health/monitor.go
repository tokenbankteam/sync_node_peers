package health

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"net/http"
)

type Monitor struct {
}

// StartPprof start http monitor.
func InitMonitor(binds []string) {
	m := new(Monitor)
	monitorServeMux := http.NewServeMux()
	monitorServeMux.HandleFunc("/v1/ping", m.Ping)
	for _, addr := range binds {
		log.Infof("start monitor listen: %v", addr)
		go func(bind string) {
			if err := http.ListenAndServe(bind, monitorServeMux); err != nil {
				log.Errorf("http.ListenAndServe(\"%s\", monitorServeMux) error(%v)", addr, err)
				panic(err)
			}
		}(addr)
	}
}

// monitor ping
func (m *Monitor) Ping(w http.ResponseWriter, r *http.Request) {
	type Ret struct {
		Message string `json:"message"`
	}
	marshal, _ := json.Marshal(Ret{
		Message: "ok",
	})
	w.Write(marshal)
}
