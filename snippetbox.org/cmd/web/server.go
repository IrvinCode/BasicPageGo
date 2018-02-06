package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)
func (app *App) RunServer() {
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:
		[]tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Add Idle, Read and Write timeouts to the server.
	srv := &http.Server{
		Addr:         app.Addr,
		Handler:      app.Routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on %s", app.Addr)
	err := srv.ListenAndServeTLS(app.TLSCert, app.TLSKey)
	log.Fatal(err)
}
