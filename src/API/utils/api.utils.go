package apiUtils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	apiModels "price-tracking-products/src/API/models"
)

func GetBody(body io.ReadCloser, receiver any) error {
	err := json.NewDecoder(body).Decode(receiver)
	if err != nil {
		log.Println("Error decoding the Body of the request", err)
		return err
	}
	return nil
}

func CreateResponse(w http.ResponseWriter, statusCode int, response *apiModels.GenericResponse) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error serializing the response ", err)
		return
	}
}
