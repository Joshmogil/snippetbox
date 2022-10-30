package main

import (
"database/sql"
"flag"
"log"
"net/http"
"os"
"github.com/Joshmogil/snippetbox/pkg/models/mysql"
_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Addr		string
	StaticDir	string
	dsn			string
}

type application struct {
	errorLog 	*log.Logger
	infoLog 	*log.Logger
	snippets 	*mysql.SnippetModel
}

func main() {
	cfg := new(Config)

	
	flag.StringVar(&cfg.Addr,"addr", ":4000", "HTTP net addr")
	flag.StringVar(&cfg.StaticDir, "static-dir","./ui/static","Path to static assets")
	flag.StringVar(&cfg.dsn,"dsn", "web:password@/snippetbox?parseTime=true", "MySQL database connection string")
	//dsn := flag.String("dsn", "web:password@/snippetbox?parseTime=true", "MySQL database connection string")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog:	errorLog,
		infoLog: 	infoLog,
		snippets:	&mysql.SnippetModel{DB: db},
	}

	srv := &http.Server {
		Addr: cfg.Addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s",cfg.Addr)
	err = srv.ListenAndServe()
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