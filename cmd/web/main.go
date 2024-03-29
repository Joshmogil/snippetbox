package main

import (
"crypto/tls"
"database/sql"
"flag"
"html/template"
"log"
"net/http"
"time"
"os"
"github.com/Joshmogil/snippetbox/pkg/models/mysql"
_ "github.com/go-sql-driver/mysql"
"github.com/golangcollege/sessions"
)

type contextKey string

var contextKeyUser = contextKey("user")

type Config struct {
	Addr		string
	StaticDir	string
	dsn			string
	secret		string
}

type application struct {
	errorLog 	*log.Logger
	infoLog 	*log.Logger
	session 	*sessions.Session
	snippets 	*mysql.SnippetModel
	templateCache map[string]*template.Template
	users		*mysql.UserModel
}

func main() {
	cfg := new(Config)

	
	flag.StringVar(&cfg.Addr,"addr", ":4000", "HTTP net addr")
	flag.StringVar(&cfg.StaticDir, "static-dir","./ui/static","Path to static assets")
	flag.StringVar(&cfg.dsn,"dsn", "web:password@/snippetbox?parseTime=true", "MySQL database connection string")
	flag.StringVar(&cfg.secret, "secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret for session manager")
	
	//dsn := flag.String("dsn", "web:password@/snippetbox?parseTime=true", "MySQL database connection string")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(cfg.secret))
	session.Lifetime = 12 * time.Hour

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	app := &application{
		errorLog:		errorLog,
		infoLog: 		infoLog,
		session: 		session,
		snippets:		&mysql.SnippetModel{DB: db},
		templateCache: 	templateCache,
		users: 			&mysql.UserModel{DB: db},
	}

	srv := &http.Server {
		Addr: cfg.Addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
		TLSConfig: tlsConfig,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s",cfg.Addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem","./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}