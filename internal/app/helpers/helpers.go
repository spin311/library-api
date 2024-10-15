package helpers

import (
	"github.com/spin311/library-api/internal/repository/models"
	"log"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	log.Println(err)
	http.Error(w, err.Error(), statusCode)
}

func WriteHttpErrorResponse(w http.ResponseWriter, err models.HttpError) {
	log.Println(err)
	http.Error(w, err.Message, err.StatusCode)
}
