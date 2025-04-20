package httpserver

import (
	"net/http"

	"github.com/gsouza97/my-bots/internal/logger"
)

func StartHealthCheckServer() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Log.Fatalf("Erro ao iniciar o health check server: %v", err)
	}
}
