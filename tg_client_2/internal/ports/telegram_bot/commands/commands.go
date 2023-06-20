package commands

type Command int64

const (
	Unknown Command = iota
	Quiz
	TakeParty
	LookTop
	BackToQuizList
	SwitchAnswer
	Confirm
)

func (c Command) String() string {
	switch c {
	case Quiz:
		return "Quiz"
	case TakeParty:
		return "TakeParty"
	case LookTop:
		return "LookTop"
	case BackToQuizList:
		return "BackToQuizList"
	case SwitchAnswer:
		return "SwitchAnswer"
	case Confirm:
		return "Confirm"
	default:
		return "Unknown"
	}
}

func (c Command) IsValid() bool {
	switch c {
	case Unknown, Quiz, TakeParty, LookTop, BackToQuizList, SwitchAnswer, Confirm:
		return true
	}
	return false
}
