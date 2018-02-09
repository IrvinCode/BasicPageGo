package main

import (
	"database/sql"
	"flag"
	"log"
	"time"

	"snippetbox.org/pkg/models"

	"github.com/alexedwards/scs"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "sb:pass@/snippetbox?parseTime=true", "MySQL DSN")
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	secret := flag.String("secret", "sb04y4ER5irMeOppyf5qdJG9kQSjWw2F", "Secret key")
	top := flag.String("top", "sb04y4ER5irMeOppyf5qdJG9kQSjWw8G", "Secret key top")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")
	tlsCert := flag.String("tls-cert", "./tls/cert.pem", "Path to TLS certificate")
	tlsKey := flag.String("tls-key", "./tls/key.pem", "Path to TLS key")

	flag.Parse()

	db := connect(*dsn)
	defer db.Close()

	sessionManager := scs.NewCookieManager(*secret)
	sessionManager.Lifetime(12 * time.Hour)
	sessionManager.Persist(true)
	sessionManager.Secure(true)

	sessionAdmin := scs.NewCookieAdmin(*top)
	sessionAdmin.Lifetime(12 * time.Hour)
	sessionAdmin.Persist(true)
	sessionAdmin.Secure(true)

	app := &App{
		Addr:      *addr,
		Database:  &models.Database{db},
		HTMLDir:   *htmlDir,
		Sessions:   sessionManager,
		Admin:      sessionAdmin,
		StaticDir: *staticDir,
		TLSCert:   *tlsCert,
		TLSKey:    *tlsKey,
	}

	app.RunServer()

}

func connect(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
