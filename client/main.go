package main

import (
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
