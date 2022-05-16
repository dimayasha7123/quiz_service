package models

type UserAccount struct {
	ID   int64
	Name string
}

type Quiz struct {
	ID   int64
	Name string
}

type Participation struct {
	ID            int64
	UserAccountID int64
	QuizID        int64
}

type ResponseReport struct {
	ID              int64
	ParticipationID int64
	Correct         bool
	PenaltyTime     int
}
