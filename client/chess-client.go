package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	pb "github.com/schafer14/grpc-chess/service"
)

// Option is a object representing the UCI option object
type Option struct {
	// The option has the name id.
	// 	Certain options have a fixed value for , which means that the semantics of this option is fixed.
	// 	Usually those options should not be displayed in the normal engine options window of the GUI but
	// 	get a special treatment. "Pondering" for example should be set automatically when pondering is
	// 	enabled or disabled in the GUI options. The same for "UCI_AnalyseMode" which should also be set
	// 	automatically by the GUI. All those certain options have the prefix "UCI_" except for the
	// 	first 6 options below. If the GUI get an unknown Option with the prefix "UCI_", it should just
	// 	ignore it and not display it in the engine's options dialog.
	Name string
	// The option has type t.
	// There are 5 different types of options the engine can send
	// * check
	// 	a checkbox that can either be true or false
	// * spin
	// 	a spin wheel that can be an integer in a certain range
	// * combo
	// 	a combo box that can have different predefined strings as a value
	// * button
	// 	a button that can be pressed to send a command to the engine
	// * string
	// 	a text field that has a string as a value,
	// 	an empty string has the value ""
	// todo: make this a enum instead of a string
	Type string
	// the default value of this parameter is x
	Default string
	// the minimum value of this parameter is x
	Min int32
	// the maximum value of this parameter is x
	Max int32
	// a predefined value of this parameter is x
	Var string
}

// Engine defines the required specification for interfacing with the UCI over gRPC protocol
type Engine interface {
	// Id returns the engine name and the engine author
	Id() (string, string)
	// Options returns a list of options that the engine accepts
	Options() []Option
}

type chessClient struct {
	l log.Entry
	c pb.ChessApplicationClient
	e Engine
}

func newChessClient(l log.Entry, client pb.ChessApplicationClient) chessClient {
	return chessClient{l: l, c: client}
}

// Runs through the process of creating a chess game
// TODO: listen for quits
func (c chessClient) newGameRequest() {
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

	// Send ID message
	name, author := c.e.Id()
	err = stream.Send(&pb.UciRequest{
		MessageType: pb.UciRequest_ID,
		Id: &pb.UciRequest_Id{
			Name:   name,
			Author: author,
		},
	})
	if err != nil {
		cancel()
		requestLogger.Errorln("Could not send uci id request", err)
		return
	}

	// Send Option messages
	for _, opt := range c.e.Options() {
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
