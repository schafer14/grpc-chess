package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	pb "github.com/schafer14/grpc-chess/service"
)

type chessClient struct {
	l log.Entry
	c pb.ChessApplicationClient
}

func newChessClient(l log.Entry, client pb.ChessApplicationClient) chessClient {
	return chessClient{l: l, c: client}
}

func (c chessClient) newGameRequest() {
	requestLogger := c.l.WithField("request", "newGameRequest")
	requestLogger.Info("Requesting a new game")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.c.UCI(ctx)
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not create game request", err)
		return
	}

	// Recieves a UCI response
	uciMessage, err := stream.Recv()
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not read message", err)
		return
	}
	requestLogger.Info(uciMessage.GetMessageType().String())

	err = stream.Send(&pb.UciRequest{
		MessageType: pb.UciRequest_ID,
		Id: &pb.UciRequest_Id{
			Name:   "MtM",
			Author: "Banner",
		},
	})
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not send uci id request", err)
		return
	}
	err = stream.Send(&pb.UciRequest{
		MessageType: pb.UciRequest_OPTION,
		Option: &pb.UciRequest_Option{
			Name: "Hash",
			Type: "Check",
		},
	})
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not send option request", err)
		return
	}
	err = stream.Send(&pb.UciRequest{
		MessageType: pb.UciRequest_UCIOK,
	})
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not send uciok request", err)
		return
	}

	stream.Recv()

}
