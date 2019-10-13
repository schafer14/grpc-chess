package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/schafer14/grpc-chess/client"
	engine "github.com/schafer14/grpc-chess/engine/uci"
	pb "github.com/schafer14/grpc-chess/service"
	"google.golang.org/grpc"
)

func main() {
	clientLogger := log.WithField("from", "client")
	clientLogger.Info("Starting")

	host := flag.String("host", ":8080", "The server host")
	executable := flag.String("executable", "/home/banner/Documents/proj/Stockfish/stockfish-10-linux/Linux/stockfish_10_x64", "Path to the uci engine executable")

	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChessApplicationClient(conn)
	agent, err := engine.New(*executable)
	if err != nil {
		clientLogger.Errorln(err)
	}
	stockfish := client.New(agent, *clientLogger, c)

	stockfish.NewGameRequest()
}
