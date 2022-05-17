package models

type UserAccount struct {
	ID   int64
	Name string
}

type Quiz struct {
	ID   int64
	Name string
}

type Party struct {
	ID            int64
	UserAccountID int64
	QuizID        int64
	Completed     bool
	Points        int32
	Questions     []Question
}

type Question struct {
	ID      int64
	Title   string
	Tags    []string
	Answers []Answer
}

type Answer struct {
	ID      int64
	Title   string
	Correct bool
}
