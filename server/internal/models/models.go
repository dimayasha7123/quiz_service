package models

type UserAccount struct {
	ID   int64
	Name string
}

type Quiz struct {
	ID    int64
	Title string
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

type PartyResults struct {
	Name   string
	Points int32
	Place  int64
}

type SingleTop struct {
	UserResults PartyResults
	GlobalTop   GlobalTop
}

type GlobalTop struct {
	Results []PartyResults
}
