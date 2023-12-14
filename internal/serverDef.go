package internal

import (
	"fmt"
	"net/http"

	_ "github.com/DroidZed/go_lance/docs"
	"github.com/DroidZed/go_lance/internal/auth"
	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/pigeon"
	"github.com/DroidZed/go_lance/internal/signup"
	"github.com/DroidZed/go_lance/internal/user"
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/MadAppGang/httplog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
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
	server.EnvConfig = config.LoadEnv()
	pigeon.GetSmtp()
	return server
}

func (s *Server) MountHandlers() {

	// Mount all handlers here
	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.JsonResponse(w, 200, utils.DtoResponse{Message: "Hello Go Lance!"})
	})

	s.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", s.EnvConfig.Port)), //The url pointing to API definition
	))

	s.Router.Mount("/user", user.UserRoutes())
	s.Router.Mount("/auth", auth.AuthRoutes())
	s.Router.Mount("/signup", signup.SignUpRoutes())
}

func (s *Server) ApplyMiddleWares() {

	s.Router.Use(middleware.RequestID)

	s.Router.Use(middleware.CleanPath)

	s.Router.Use(middleware.URLFormat)

	s.Router.Use(middleware.StripSlashes)

	s.Router.Use(httplog.LoggerWithName("Go Lance"))

	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.Router.Use(middleware.Heartbeat("/health"))
}
