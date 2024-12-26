package gapi

import (
	"context"

	"github.com/dubass83/go-micro-logger/data"
	"github.com/dubass83/go-micro-logger/pb"
	"github.com/rs/zerolog/log"
)

func (l *LogServer) WriteLog(ctx context.Context, req *pb.LogRequest) (*pb.LogResponse, error) {
	entry := data.LogEntry{
		Name: req.LogEntry.Name,
		Data: req.LogEntry.Data,
	}

	err := l.LogStorage.Insert(entry)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert a log into storage by gRPC")
		return &pb.LogResponse{Result: "failed to insert a log into storage"}, err
	}

	return &pb.LogResponse{Result: "Processed payload via gRPC: " + req.LogEntry.Name}, nil
}
