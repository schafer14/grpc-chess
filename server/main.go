package main

import (
	"flag"
	"net"

	pb "github.com/schafer14/grpc-chess/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	serverLogger := log.WithField("from", "server")
	serverLogger.Info("Starting")

	err := run(*serverLogger)
	if err != nil {
		serverLogger.Fatalf("Could not start server %v", err)
	}
}

func run(logger log.Entry) error {
	host := flag.String("host", ":8080", "The server host")

	flag.Parse()

	lis, err := net.Listen("tcp", *host)

	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	pb.RegisterChessApplicationServer(grpcServer, NewChessService(logger))

	logger.WithField("port", *host).Info("Listening")
	logger.Fatal(grpcServer.Serve(lis))

	return nil
}
