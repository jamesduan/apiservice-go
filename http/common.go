package http

import (
	"apiservice/g"
	"net/http"
	"strings"

	"github.com/toolkits/file"
)

func configCommonRoutes() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, map[string]string{"Version": g.VERSION})
		// w.Write([]byte(g.VERSION))
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, map[string]string{"Workdir": file.SelfDir()})
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			g.ParseConfig(g.ConfigFile)
			RenderDataJson(w, g.Config())
		} else {
			// w.Write([]byte("no privilege"))
			RenderDataJson(w, map[string]string{"Permission": "no privilege"})
		}
	})
}
