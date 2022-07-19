package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"party/models"
	"strconv"

	"github.com/gorilla/mux"
)

func SendError(w http.ResponseWriter, status int, err models.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func SendServerErrorIfErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		errorMessage := models.Error{Message: "Server Error"}
		SendError(w, http.StatusInternalServerError, errorMessage)
		return true
	}
	return false
}

func HandleResponse(w http.ResponseWriter, err error, data interface{}) {
	if err != nil {
		log.Default().Panic(err)
		errorMessage := models.Error{Message: "Server Error"}
		SendError(w, http.StatusInternalServerError, errorMessage)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	SendSuccess(w, data)
}

func RemoveFromArray(arr []int, item int) []int {
	for i, v := range arr {
		if v == item {
			arr = append(arr[:i], arr[i+1:]...)
			break
		}
	}
	return arr
}

func GetIntVar(r *http.Request, key string) int {
	value := mux.Vars(r)[key]
	intVal, _ := strconv.Atoi(value)
	return intVal
}

func GetReqBody[T any](r *http.Request, data T) T {
	json.NewDecoder(r.Body).Decode(&data)
	return data
}
