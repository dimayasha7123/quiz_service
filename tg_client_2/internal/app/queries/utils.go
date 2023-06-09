package queries

import (
	"context"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

func convertResultsFromApiToApp(ctx context.Context, sessions domain.Sessions, top *api.SingleTop) app.Results {
	ret := app.Results{}
	if top.UserResults != nil {
		ret.UserResults = app.UserResults{
			Place:  top.UserResults.Place,
			Points: int64(top.UserResults.PointCount),
		}
	}

	if len(top.QuizTop.Results) != 0 {
		topResults := make(app.TopResults, 0, len(top.QuizTop.Results))
		for _, result := range top.QuizTop.Results {
			topResults = append(topResults, app.ResultRow{
				Username: sessions.GetName(ctx, result.Name),
				Points:   result.Place,
			})
		}
		ret.TopResults = topResults
	}

	return ret
}
