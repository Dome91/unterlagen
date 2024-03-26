package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unterlagen/pkg/config"
	"unterlagen/pkg/domain"
	"unterlagen/pkg/web/auth"
	"unterlagen/pkg/web/handlers"
	"unterlagen/views"
)

func StartServer(documents *domain.Documents, folders *domain.Folders, users *domain.Users) {
	router := chi.NewRouter()
	registerMiddleware(router, users)

	executor := handlers.NewTemplateExecutor()
	showLogin := handlers.ShowLogin(executor)
	loginUser := handlers.LoginUser(users)
	logoutUser := handlers.LogoutUser()
	createFolder := handlers.CreateFolder(folders)
	getFolder := handlers.GetFolder(documents, folders, executor)
	createDocument := handlers.UploadDocument(documents)
	downloadDocument := handlers.DownloadDocument(documents)

	router.Get("/", showLogin)
	router.Post("/", loginUser)
	router.Post("/logout", logoutUser)
	router.Post("/folders", createFolder)
	router.Get("/folders", getFolder)
	router.Post("/documents", createDocument)
	router.Get("/documents/:id", downloadDocument)
	serveAssets(router)

	err := users.CreateAdmin()
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf(":%s", config.Get().Port)
	server := &http.Server{Addr: addr, Handler: router}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				panic(err)
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			panic(err)
		}

		serverStopCtx()
	}()

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	<-serverCtx.Done()
	log.Info().Msg("Stopped server")
}

func StopServer() {
	log.Info().Msg("Stopping server")
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if err != nil {
		panic(err)
	}
}

func serveAssets(router chi.Router) {
	if config.Get().Development {
		assets := http.FileServer(http.Dir("views/public"))
		router.Handle("/*", assets)
	} else {
		assets, err := fs.Sub(views.Assets, "public")
		if err != nil {
			panic(err)
		}

		router.Handle("/*", http.FileServer(http.FS(assets)))
	}
}

func configureLogging(router chi.Router) {
	loggingMw := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				log.Info().Fields(map[string]interface{}{
					"remote_addr": r.RemoteAddr,
					"path":        r.URL.Path,
					"method":      r.Method,
					"user_agent":  r.UserAgent(),
					"status":      http.StatusText(ww.Status()),
				}).Msg("Request")
			}()
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}

	router.Use(loggingMw)
}

func registerMiddleware(router chi.Router, users *domain.Users) {
	router.Use(middleware.Recoverer)
	configureLogging(router)
	auth.ConfigureSession(router)
}
