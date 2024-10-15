package helpers

import (
	"encoding/json"
	"github.com/spin311/library-api/internal/repository/models"
	"log"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	htmlError := models.NewHttpError(err.Error(), statusCode)
	WriteHttpErrorResponse(w, htmlError)
}

func WriteHttpErrorResponse(w http.ResponseWriter, htmlError models.HttpError) {
	log.Println(htmlError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(htmlError.StatusCode)
	jsonErr := json.NewEncoder(w).Encode(htmlError)
	if jsonErr != nil {
		http.Error(w, htmlError.Message, htmlError.StatusCode)
	}

}
