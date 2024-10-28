package api

import (
	"net/http"

	"github.com/dubass83/go-micro-logger/data"
	"github.com/rs/zerolog/log"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Test api Handler
func Test(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Massage: "Hello from Logger!",
	}

	_ = writeJSON(w, http.StatusAccepted, payload)
}

// WriteLog save geting logs from POST request to LogStorage
func (s *Server) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload
	err := readJSON(w, r, &requestPayload)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request payload")
		err = errorJSON(w, err)
		if err != nil {
			log.Error().Err(err).Msg("failed to write error into ")
		}
		return
	}

	entry := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = s.LogStorage.Insert(entry)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert a log into storage")
		_ = errorJSON(w, err)
		return
	}

	response := jsonResponse{
		Error:   false,
		Massage: "successfully incerted a log into storage",
	}

	_ = writeJSON(w, http.StatusAccepted, response)
}
