package user

// CRUD: user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/ggicci/httpin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController interface {
	GetAllUsers(w http.ResponseWriter, _ *http.Request)
	GetUserById(w http.ResponseWriter, _ *http.Request)
	DeleteUserById(w http.ResponseWriter, r *http.Request)
	UpdateUserById(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	log := config.InitializeLogger().LogHandler

	user := &User{ID: primitive.NewObjectID()}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.Fatal(err)
	}

	userId, err := SaveOne(user)
	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("Error creating the user!\n%s", err.Error())})
		return
	}

	utils.JsonResponse(w, http.StatusCreated, userId)

}

func GetAllUsers(w http.ResponseWriter, _ *http.Request) {

	log := config.InitializeLogger().LogHandler

	users, err := FindAllUsers()

	if err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusOK, []User{})
		return
	}

	utils.JsonResponse(w, http.StatusOK, users)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value(httpin.Input).(*utils.UserIdPath).UserId

	log := config.InitializeLogger().LogHandler

	log.Infof("id: %s", id)

	user, err := FindUserByID(id)

	if err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(httpin.Input).(*utils.UserIdPath).UserId

	result := DeleteOne(id)

	if !result {
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, utils.DtoResponse{Message: fmt.Sprintf("User with id: %s has been deleted successfully!", id)})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	log := config.InitializeLogger().LogHandler

	user := &User{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid JSON payload"})
		return
	}

	log.Debug(*user)

	err := UpdateOneUser(*user)
	if err != nil {
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("Invalid update!\n %s", err.Error())})
		return
	}

	utils.JsonResponse(w, http.StatusOK, utils.DtoResponse{Message: "Success !!"})
}
