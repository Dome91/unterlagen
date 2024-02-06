package web

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"io/fs"
	"net/http"
	"unterlagen/pkg/config"
	"unterlagen/pkg/domain"
	"unterlagen/pkg/web/auth"
	"unterlagen/pkg/web/handlers"
	"unterlagen/views"
)

func StartServer(documents *domain.Documents, folders *domain.Folders, users *domain.Users) error {
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

	addr := fmt.Sprintf(":%s", config.Port)
	return http.ListenAndServe(addr, router)
}

func serveAssets(router chi.Router) {
	if config.Development {
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
	auth.ConfigureSession(router, users)
}
