package main

import (
	"context"

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
	return chessService{l}
}

// 3
func (cs chessService) GameStream(*chess.GameRequestMessage, chess.ChessApplication_GameStreamServer) error {
	panic("not implemented")
}

// 3
func (cs chessService) GameAction(chess.ChessApplication_GameActionServer) error {
	panic("not implemented")
}

// not important
func (cs chessService) MainChatRoom(*chess.RoomRequest, chess.ChessApplication_MainChatRoomServer) error {
	panic("not implemented")
}

// 1
func (cs chessService) GameRequest(gameControls *chess.GameControls, res chess.ChessApplication_GameRequestServer) error {
	logger := cs.l.WithField("request", "GameRequest")

	logger.Info("Got a request for a new game.", "Game Controls", gameControls)

	return res.Send(&pb.GameProposals{
		Opponent: &pb.Person{
			Name: "Magnus",
			Id:   "1",
		},
	})
}

// 2
func (cs chessService) GameConfirmation(context.Context, *chess.GameRequestMessage) (*chess.Confimation, error) {
	panic("not implemented")
}
