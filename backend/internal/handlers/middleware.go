package handlers

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func TailscaleOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("Eingehende Anfrage von IP: %s\n", r.RemoteAddr)

		host, port, err := net.SplitHostPort(r.RemoteAddr)

		fmt.Printf("Eingehende Anfrage von Host: %s, Port: %s\n", host, port)

		if err != nil {
			http.Error(w, "Ungültige Remote-Adresse", http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(r.RemoteAddr, "100.") && (host != "127.0.0.1") && (host != "::1") {
			http.Error(w, "Zugriff verweigert: Nur Tailscale-Netzwerk erlaubt", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)

	})
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
