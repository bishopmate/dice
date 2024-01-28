package core

import (
	"errors"
	"net"
)

func evalPING(args []string, c net.Conn) error {
	argumentsLength := len(args)
	var b []byte

	if argumentsLength >= 2 {
		return errors.New("wrong number of arguments for 'ping' command")
	}

	if argumentsLength == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	_, err := c.Write(b)

	return err
}

func EvalAndRespond(cmd *RedisCmd, c net.Conn) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, c)
	default:
		return evalPING(cmd.Args, c)
	}
}
