package models

type UserAccount struct {
	ID   int
	Name string
}

type Quiz struct {
	ID   int
	Name string
}

type Participation struct {
	ID            int
	UserAccountID int
	QuizID        int
}

type ResponseReport struct {
	ID              int
	ParticipationID int
	Correct         bool
	PenaltyTime     int
}
