package gapi

import (
	"context"

	"github.com/dubass83/go-micro-logger/data"
	"github.com/dubass83/go-micro-logger/pb"
	"github.com/rs/zerolog/log"
)

func (l *LogServer) WriteLog(ctx context.Context, req *pb.LogRequest) (*pb.LogResponse, error) {
	input := req.GetLogEntry()
	entry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}
	log.Debug().Msgf("Received a log entry: %+v", entry)

	err := l.LogStorage.Insert(entry)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert a log into storage by gRPC")
		return &pb.LogResponse{Result: "failed to insert a log into storage"}, err
	}

	return &pb.LogResponse{Result: "Processed payload via gRPC: " + input.Name}, nil
}
