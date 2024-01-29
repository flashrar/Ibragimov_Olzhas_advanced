package main
import (
"crypto/tls"
"database/sql"
"flag"
"html/template"
"log/slog"
"net/http"
"os"
"time"
"snippetbox.alexedwards.net/internal/models"
"github.com/alexedwards/scs/mysqlstore"
"github.com/alexedwards/scs/v2"
"github.com/go-playground/form/v4"
_ "github.com/go-sql-driver/mysql"
)
type application struct {
logger *slog.Logger
snippets models.SnippetModelInterface // Use our new interface type.
users models.UserModelInterface // Use our new interface type.
templateCache map[string]*template.Template
formDecoder *form.Decoder
sessionManager *scs.SessionManager
}

		

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:flarar22@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB(*dsn)
	if err != nil {
	logger.Error(err.Error())
	os.Exit(1)
	}
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
	logger.Error(err.Error())
	os.Exit(1)
	}
	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Make sure that the Secure attribute is set on our session cookies.
	// Setting this means that the cookie will only be sent by a user's web
	// browser when a HTTPS connection is being used (and won't be sent over an
	// unsecure HTTP connection).
	sessionManager.Cookie.Secure = true
	app := &application{
	logger: logger,
	snippets: &models.SnippetModel{DB: db},
	users: &models.UserModel{DB: db},
	templateCache: templateCache,
	formDecoder: formDecoder,
	sessionManager: sessionManager,
	}
	tlsConfig := &tls.Config{
	CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := &http.Server{
	Addr: *addr,
	Handler: app.routes(),
	ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	TLSConfig: tlsConfig,
	IdleTimeout: time.Minute,
	ReadTimeout: 5 * time.Second,
	WriteTimeout: 10 * time.Second,
	}
	logger.Info("starting server", "addr", srv.Addr)
	err = srv.ListenAndServeTLS("C:/Users/User/code/snippetbox/tls/cert.pem", "C:/Users/User/code/snippetbox/tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}
	
	
	

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
	
	
