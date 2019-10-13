package client

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	pb "github.com/schafer14/grpc-chess/service"
)

type chessClient struct {
	l log.Entry
	c pb.ChessApplicationClient
	e Engine
}

// New creates a new client that can be used in grpc clients
func New(engine Engine, l log.Entry, client pb.ChessApplicationClient) Client {
	return chessClient{e: engine, l: l, c: client}
}

// Runs through the process of creating a chess game
func (c chessClient) NewGameRequest() {
	// Setup the logger
	requestLogger := c.l.WithField("request", "newGameRequest")
	requestLogger.Info("Requesting a new game")

	// Setup the context that will be used as a base context throughout
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start UCI stream
	stream, err := c.c.UCI(ctx)
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not create game request", err)
		return
	}

	// Recieves a UCI message
	uciMessage, err := stream.Recv()
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not read message", err)
		return
	}
	requestLogger.Info(uciMessage.GetMessageType().String())

	// Get engine ident and options
	engineIdent, options, err := c.e.Init()
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not init engine", err)
		return
	}
	err = stream.Send(&pb.UciRequest{
		MessageType: pb.UciRequest_ID,
		Id: &pb.UciRequest_Id{
			Name:   engineIdent.Name,
			Author: engineIdent.Author,
		},
	})
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not send uci id request", err)
		return
	}

	// Send Option messages
	for _, opt := range options {
		err = stream.Send(&pb.UciRequest{
			MessageType: pb.UciRequest_OPTION,
			Option: &pb.UciRequest_Option{
				Name:    opt.Name,
				Type:    opt.Type,
				Default: opt.Default,
				Min:     opt.Min,
				Max:     opt.Max,
				Var:     opt.Var,
			},
		})
		if err != nil {
			cancel()
			requestLogger.Errorln("Could not send option request", err)
			return
		}
	}

	// Send UCIOK
	err = stream.Send(&pb.UciRequest{
		MessageType: pb.UciRequest_UCIOK,
	})
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not send uciok request", err)
		return
	}

	// Listening for either setoption messages or isready messages
Loop:
	for {
		message, err := stream.Recv()
		if err != nil {
			cancel()
			requestLogger.Errorln("Could read set option/is ready message", err)
			return
		}

		switch message.GetMessageType() {
		case pb.UciResponse_SETOPTION:
			requestLogger.Infof("Setting option %v is not supported yet.", message.GetSetOption().GetName())
		case pb.UciResponse_ISREADY:
			break Loop
		}
	}

	// Send an ready okay message
	stream.Send(&pb.UciRequest{
		MessageType: pb.UciRequest_READYOK,
	})

	handleGameLogic(stream, requestLogger)
	stream.Recv()
}

// handleGameLogic is responsible for managing the relationship between the engine and the server
func handleGameLogic(stream pb.ChessApplication_UCIClient, logger *logrus.Entry) error {
	// Setup a new context
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	// Setup a channel to listen for messages on
	inChan := make(chan pb.UciResponse)
	go func(input chan pb.UciResponse) {
		for {
			in, err := stream.Recv()
			if err != nil {
				cancel()
				return
			}
			input <- *in
		}
	}(inChan)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Connection closed")
			return fmt.Errorf("Context ended")
		case msg := <-inChan:
			switch msg.GetMessageType() {
			case pb.UciResponse_POSITION:
				logger.Warn("Unimplemented uci position")
			case pb.UciResponse_GO:
				logger.Warn("Unimplemented uci go")
			case pb.UciResponse_UCINEWGAME:
				logger.Warn("Unimplemented uci ucinewgame")
			case pb.UciResponse_PONDERHIT:
				logger.Warn("Unimplemented uci ponderhit")
			case pb.UciResponse_STOP:
				logger.Warn("Unimplemented uci stop")
			case pb.UciResponse_QUIT:
				logger.Warn("Unimplemented uci quit")
				break
			default:
				logger.Errorf("Unknown uci message %v", msg.GetMessageType())
				break
			}
		}
	}
}
