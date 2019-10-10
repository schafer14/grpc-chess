package main

import (
	"context"
	"flag"

	log "github.com/sirupsen/logrus"

	pb "github.com/schafer14/grpc-chess/service"
	"google.golang.org/grpc"
)

func main() {
	clientLogger := log.WithField("from", "client")
	clientLogger.Info("Starting")

	host := flag.String("host", ":8080", "The server host")

	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChessApplicationClient(conn)
	client := newChessClient(*clientLogger, c)

	client.newGameRequest()
}

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

	stream, err := c.c.GameRequest(ctx, &pb.GameControls{})
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not create game request", err)
		return
	}

	prop, err := stream.Recv()
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not read message", err)
		return
	}

	requestLogger.WithField("proposal", prop).Info("Got a game request")
}
