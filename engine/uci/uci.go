package uci

import (
	"bufio"
	"io"
	"os/exec"
	"strings"

	cli "github.com/schafer14/grpc-chess/client"
)

type uci struct {
	in  bufio.Reader
	out io.Writer
	cmd exec.Cmd
}

// New returns a new UCI instance
func New(path string) (cli.Engine, error) {
	command := exec.Command(path)

	out, err := command.StdinPipe()
	if err != nil {
		return nil, err
	}
	in, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = command.Start()

	return uci{*bufio.NewReader(in), out, *command}, err
}

// Init initializes the engine returning engine options and engine ident
func (uci uci) Init() (ident cli.EngineIdent, options []cli.Option, err error) {
	uci.out.Write([]byte("uci\n"))

Loop:
	for {
		// read string
		msg, err := uci.in.ReadString('\n')
		if err != nil {
			return cli.EngineIdent{}, nil, err
		}

		// match command
		cmd := strings.Split(msg, " ")
		switch cmd[0] {
		case "uciok\n":
			break Loop
		case "option":
			option := parseOptions(msg)
			options = append(options, option)
		case "id":
			parseIdent(cmd, &ident)
		}
	}

	return ident, options, nil
}
