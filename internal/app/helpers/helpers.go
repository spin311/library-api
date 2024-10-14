package helpers

import (
	"log"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	log.Println(err)
	http.Error(w, http.StatusText(statusCode), statusCode)
}
