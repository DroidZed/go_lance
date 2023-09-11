package internal

import (
	"net/http"

	"github.com/DroidZed/go_lance/internal/auth"
	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/user"
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/MadAppGang/httplog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	Router    *chi.Mux
	DbClient  *mongo.Client
	EnvConfig *config.EnvConfig
}

func CreateNewServer() *Server {
	server := &Server{}
	server.Router = chi.NewRouter()
	server.DbClient = config.GetConnection()
	server.EnvConfig = config.LoadConfig()
	return server
}

func (s *Server) MountHandlers() {

	// Mount all handlers here
	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.JsonResponse(w, 200, utils.DtoResponse{Message: "Hello Go Lance!"})
	})

	s.Router.Mount("/user", user.UserRoutes())
	s.Router.Mount("/auth", auth.AuthRoutes())
}

func (s *Server) ApplyMiddleWares() {

	// Mount all Middleware here
	s.Router.Use(middleware.RequestID)

	s.Router.Use(middleware.CleanPath)

	s.Router.Use(middleware.URLFormat)

	s.Router.Use(middleware.StripSlashes)

	s.Router.Use(httplog.LoggerWithName("Go Lance"))

	s.Router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	s.Router.Use(middleware.Heartbeat("/health"))
}
