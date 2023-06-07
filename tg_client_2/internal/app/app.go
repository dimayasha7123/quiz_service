package app

type Queries struct {
}

type Commands struct {
}

type App struct {
	Queries  Queries
	Commands Commands
}

func New() App {
	return App{}
}
