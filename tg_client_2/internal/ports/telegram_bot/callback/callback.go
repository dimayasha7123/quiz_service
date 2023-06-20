package callback

import (
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/commands"
	"strconv"
	"strings"
)

const reqParams = 1

type Data struct {
	NextCommand commands.Command
	Args        []string
}

func NewData(nextCommand commands.Command, args ...string) Data {
	return Data{
		NextCommand: nextCommand,
		Args:        args,
	}
}

func NewDataFromString(s string) (Data, error) {
	parts := strings.Split(s, "_")
	if len(parts) < reqParams {
		return Data{}, fmt.Errorf("too few parts, get %d, must be %d or more", len(parts), reqParams)
	}

	nextCommandString := parts[0]
	nextCommandInt, err := strconv.ParseInt(nextCommandString, 10, 64)
	if err != nil {
		return Data{}, fmt.Errorf("can't parse NextCommand = %v: %v", nextCommandString, err)
	}
	nextCommand := commands.Command(nextCommandInt)
	if !nextCommand.IsValid() {
		return Data{}, fmt.Errorf("NextCommand is not valid command")
	}

	var args []string
	if len(parts) > reqParams {
		args = make([]string, len(parts)-reqParams)
		copy(args, parts[reqParams:])
	}

	return Data{
		NextCommand: nextCommand,
		Args:        args,
	}, nil
}

// originalMessageID_nextCommand_arg1_arg2_...
func (d Data) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d", int64(d.NextCommand)))
	for _, arg := range d.Args {
		sb.WriteString(fmt.Sprintf("_%s", arg))
	}
	return sb.String()
}
