package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/db"
	"github.com/DroidZed/go_lance/internal/user"
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/MadAppGang/httplog"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Router *chi.Mux
}

func CreateNewServer() *Server {
	server := &Server{}
	server.Router = chi.NewRouter()
	return server
}

func (s *Server) MountHandlers() {
	// Mount all Middleware here
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.CleanPath)
	s.Router.Use(middleware.URLFormat)
	s.Router.Use(httplog.LoggerWithName("CHI API"))

	// Mount all handlers here
	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.JsonResponse(w, 200, utils.DtoResponse{Message: "Hello World!"})
	})

	s.Router.Mount("/user", user.UserRoutes())

	s.Router.Get("/dev", func(w http.ResponseWriter, r *http.Request) {
		utils.LogAllRoutes(s.Router)
		utils.JsonResponse(w, 200, utils.DtoResponse{Message: "Nothing will be returned. This is just a dummy message. If you're a developer, check your console."})
	})
}

func init() {
	// Register a directive named "path" to retrieve values from `chi.URLParam`,
	// i.e. decode path variables.
	httpin.UseGochiURLParam("path", chi.URLParam)
}

// Entry point, setting up chi and graceful shutdown <3
func main() {

	log := config.InitializeLogger().LogHandler

	port, err := config.EnvDbPORT()

	if err != nil {
		log.Fatal("Could not retrieve port!\n")
	}

	addr := utils.SetupHostWithPort(config.EnvHost(), port)

	// The HTTP Server
	server := &http.Server{Addr: addr, Handler: service(port)}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	client := db.GetConnection()

	// Clean up to disconnect
	defer func() {
		if err := client.Disconnect(serverCtx); err != nil {
			log.Fatal(err)
		}
	}()

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 3 seconds
		shutdownCtx, cancelFunc := context.WithTimeout(serverCtx, 3*time.Second)
		defer cancelFunc()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {

				log.Fatal("Graceful shutdown timed out.. forcing exit.\n")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	log.Info("Goodbye 🧩 👋")
}

func service(port int64) http.Handler {

	log := config.InitializeLogger().LogHandler

	server := CreateNewServer()

	server.MountHandlers()

	log.Infof("Listening on port: %d\n", port)

	return server.Router

}
