package api

import (
	"github.com/dubass83/go-micro-logger/data"
	"github.com/dubass83/go-micro-logger/util"
	"github.com/rs/zerolog/log"
)

type RPCService struct {
	LogStorage data.LogStorage
	Config     util.Config
}

type RPCPayload struct {
	Name string
	Data string
}

func CreateNewRPCService(config util.Config, logStorage data.LogStorage) *RPCService {

	s := &RPCService{
		LogStorage: logStorage,
		Config:     config,
	}
	return s
}

// LogInfo processes an RPC payload and logs the information into storage.
// It takes an RPCPayload and a reply string pointer as arguments.
// If the log entry is successfully inserted into storage, it sets the reply to a success message.
// If there is an error during insertion, it logs the error and sets the reply to an error message.
//
// Parameters:
//   - payload: A pointer to RPCPayload containing the log entry data.
//   - reply: A pointer to a string where the response message will be stored.
//
// Returns:
//   - error: An error if the log entry insertion fails, otherwise nil.
func (s *RPCService) LogInfo(payload *RPCPayload, reply *string) error {

	entry := data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}

	err := s.LogStorage.Insert(entry)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert a log into storage by RPC")
		*reply = "failed to insert a log into storage"
		return err
	}

	*reply = "Processed payload via RPC: " + payload.Name
	return nil
}
