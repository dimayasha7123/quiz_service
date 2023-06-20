package queries

import (
	"context"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/models"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

func convertResultsFromApiToApp(ctx context.Context, sessions domain.Sessions, top *api.SingleTop) models.Results {
	if top == nil {
		return models.Results{}
	}

	ret := models.Results{}
	if top.UserResults != nil {
		ret.UserResults = models.UserResults{
			Place:  top.UserResults.Place,
			Points: int64(top.UserResults.PointCount),
		}
	}

	if len(top.QuizTop.Results) != 0 {
		topResults := make(models.TopResults, 0, len(top.QuizTop.Results))
		for _, result := range top.QuizTop.Results {
			topResults = append(topResults, models.ResultRow{
				Username: sessions.GetName(ctx, result.Name),
				Points:   int64(result.PointCount),
			})
		}
		ret.TopResults = topResults
	}

	return ret
}
