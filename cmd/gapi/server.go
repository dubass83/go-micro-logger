package gapi

import (
	"github.com/dubass83/go-micro-logger/data"
	"github.com/dubass83/go-micro-logger/pb"
	"github.com/dubass83/go-micro-logger/util"
)

type LogServer struct {
	pb.UnimplementedLogServiceServer
	LogStorage data.LogStorage
	Config     util.Config
}

func CreateNewLogServer(conf util.Config, logStorage data.LogStorage) *LogServer {
	s := &LogServer{
		LogStorage: logStorage,
		Config:     conf,
	}
	return s
}
