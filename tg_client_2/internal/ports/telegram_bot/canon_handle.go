package telegram_bot

import (
	"context"
	"fmt"
)

func (s *service) canonHandle(ctx context.Context, req message) error {
	err := s.sendCanon(ctx, req.From.ID)
	if err != nil {
		return fmt.Errorf("can't send canon to chat with id = %v: %v", req.From.ID, err)
	}

	return nil
}
