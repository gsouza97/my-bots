package httpserver

import (
	"log"
	"net/http"
)

func StartHealthCheckServer() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Erro ao iniciar o health check server: %v", err)
	}
}
