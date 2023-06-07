package states

type State int64

const (
	Unknown State = iota
	Lobby
	Quiz
)

func (s State) String() string {
	switch s {
	case Lobby:
		return "In lobby"
	case Quiz:
		return "In quiz"
	default:
		return "Unknown"
	}
}
