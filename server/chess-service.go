package main

import (
	chess "github.com/schafer14/grpc-chess/service"
	pb "github.com/schafer14/grpc-chess/service"
	"github.com/sirupsen/logrus"
)

type chessService struct {
	l logrus.Entry
	// interface to store chess game state and such
}

// NewChessService creates a new chess service given a logger and a data store
// note: datastore not yet implemented
func NewChessService(l logrus.Entry) pb.ChessApplicationServer {
	return &chessService{l}
}

// UCI handles uci request from an egine. The service acts in the GUI role described in the UCI spec
func (cs chessService) UCI(stream chess.ChessApplication_UCIServer) error {
	logger := cs.l.WithField("request", "UCI")

	logger.Info("Got UCI game request")

	err := stream.Send(&pb.UciResponse{
		MessageType: pb.UciResponse_UCI,
	})
	if err != nil {
		logger.Error(err)
		return err
	}

	// Listen for id messages options messages or uci okay messages
Loop:
	for {
		message, err := stream.Recv()
		if err != nil {
			logger.Error("Could not receive a id/uciok/option message: ", err)
			return err
		}

		switch message.GetMessageType() {
		case pb.UciRequest_ID:
			logger = logger.WithField("engine", message.GetId().GetName())
		case pb.UciRequest_OPTION:
			logger.Infof("Available option %v", message.GetOption().GetName())
		case pb.UciRequest_UCIOK:
			break Loop
		}
	}

	logger.Info("Starting game")

	return nil
}
