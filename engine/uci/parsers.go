package uci

import (
	"bufio"
	"strconv"
	"strings"

	cli "github.com/schafer14/grpc-chess/client"
)

func parseIdent(tokens []string, ident *cli.EngineIdent) {
	switch tokens[1] {
	case "name":
		ident.Name = strings.Trim(strings.Join(tokens, " "), "\n")
	case "author":
		ident.Author = strings.Trim(strings.Join(tokens, " "), "\n")
	}
	return
}

func parseOptions(msg string) (option cli.Option) {
	scanner := bufio.NewReader(strings.NewReader(msg))
	scanner.Discard(6)

	for {
		str, err := scanner.ReadString(' ')
		if err != nil {
			break
		}
		str = strings.Trim(str, " ")

		switch str {
		case "name":
			option.Name, err = scanner.ReadString(' ')
			if err != nil {
				break
			}
			for {
				nextChar, _ := scanner.Peek(1)
				if nextChar[0] < 'A' || nextChar[0] > 'Z' {
					break
				}
				nextWord, err := scanner.ReadString(' ')
				if err != nil {
					break
				}
				option.Name += nextWord
			}
		case "type":
			option.Type, err = scanner.ReadString(' ')
			if err != nil {
				break
			}
		case "default":
			option.Type, err = scanner.ReadString(' ')
			if err != nil {
				break
			}
		case "var":
			v, err := scanner.ReadString(' ')
			if err != nil {
				break
			}
			option.Var = append(option.Var, v)
		case "min":
			minStr, err := scanner.ReadString(' ')
			if err != nil {
				break
			}
			minInt, err := strconv.Atoi(strings.Trim(minStr, " \n"))
			if err != nil {
				break
			}
			option.Min = int32(minInt)
		case "max":
			maxStr, _ := scanner.ReadString(' ')
			if err != nil {
				break
			}
			maxInt, err := strconv.Atoi(strings.Trim(maxStr, " \n"))
			if err != nil {
				break
			}
			option.Max = int32(maxInt)
		}

	}

	return option
}
