package http

import "net/http"

func configTestRoutes() {
	http.HandleFunc("/test/version", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, map[string]string{"Test": "heheheh", "Version": "1.0.1"})
	})
}
