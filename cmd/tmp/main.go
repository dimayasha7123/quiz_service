package main

import (
	quizApi "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/quiz_party_api_client"
	"log"
)

func main() {
	QACl := quizApi.New("cCP7GHWCoe2d4VBysvmaAT1x1hToB54mm1rdEiYK")
	ans, err := QACl.GetParty("Linux")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v\n", ans)
}
