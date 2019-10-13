package main

import (
	"context"
	"fmt"

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

	// At this point  the client can send a message of type: ID, Option, or UCIOK
	// So the serve accepts any one of these until the UCIOK comes through
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

	// The server can send any options it wants and then sends a ISREADY
	err = stream.Send(&pb.UciResponse{
		MessageType: pb.UciResponse_SETOPTION,
		SetOption: &pb.UciResponse_SetOption{
			Name:  "Hash",
			Value: "500",
		},
	})

	// Send is ready
	err = stream.Send(&pb.UciResponse{
		MessageType: pb.UciResponse_ISREADY,
	})
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("Listening for a `readyok` message")

	// Listen for ready ok message
	readyOk, err := stream.Recv()
	if err != nil {
		logger.Error(err)
		return err
	}
	if readyOk.GetMessageType() != pb.UciRequest_READYOK {
		logger.Warningf("Invalid message expecting readyok got %v", readyOk.GetMessageType().String())
		return fmt.Errorf("Invalid message expecting readyok got %v", readyOk.GetMessageType().String())
	}

	logger.Info("Recieved `readyok` message")

	return cs.handleGameLogic(stream, logger)
}

// handleGameLogic is responsible for the logic of adjudicating a game
func (cs chessService) handleGameLogic(stream pb.ChessApplication_UCIServer, logger *logrus.Entry) error {
	// Setup a new context
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	// Setup a channel to listen for messages on
	inChan := make(chan pb.UciRequest)
	go func(input chan pb.UciRequest) {
		for {
			in, err := stream.Recv()
			if err != nil {
				cancel()
				return
			}
			input <- *in
		}
	}(inChan)

	stream.Send(&pb.UciResponse{
		MessageType: pb.UciResponse_UCINEWGAME,
	})

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Connection closed")
			return fmt.Errorf("Context ended")
		case msg := <-inChan:
			switch msg.GetMessageType() {
			case pb.UciRequest_INFO:
				logger.Warn("Unimplemented uci info")
				break
			case pb.UciRequest_BESTMOVE:
				logger.Warn("Unimplemented uci best move")
				break
			default:
				logger.Errorf("Unknown uci message %v", msg.GetMessageType())
				break
			}
		}
	}
}
