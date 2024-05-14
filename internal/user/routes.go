package user

import (
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

func UserRoutes() chi.Router {
	userRouter := chi.NewRouter()

	userRouter.Route("/", func(r chi.Router) {
		r.Get("/", GetAllUsers)
		r.Post("/", CreateUser)
		r.Put("/", UpdateUser)

	})

	userRouter.With(httpin.NewInput(utils.UserIdPath{})).
		Route("/{userId}", func(r chi.Router) {
			r.Get("/", GetUserById)
			r.Delete("/", DeleteUserById)
		})

	return userRouter
}
