package app

type UserInfo struct {
	UserID   int64
	UserName string
}

type QuizList []Quiz

type Quiz struct {
	ID    int64
	Title string
}

type UserQuizIDs struct {
	UserID int64
	QuizID int64
}

type Results struct {
	UserResults UserResults
	TopResults  TopResults
}

type UserResults struct {
	Place  int64
	Points int64
}

type TopResults []ResultRow

type ResultRow struct {
	Username string
	Points   int64
}

type Question struct {
	Title   string
	Answers Answers
}

type QuestionInfo struct {
	Exist    bool
	Question Question
}

type Answers []Answer

type Answer struct {
	Title  string
	Picked bool
}
