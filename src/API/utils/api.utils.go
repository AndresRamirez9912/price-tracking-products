package apiUtils

import (
	"encoding/json"
	"io"
	"log"
)

func GetBody(body io.ReadCloser, receiver any) error {
	err := json.NewDecoder(body).Decode(receiver)
	if err != nil {
		log.Println("Error decoding the Body of the request", err)
		return err
	}
	return nil
}
