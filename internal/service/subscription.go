package service

import (
	"context"
	"time"

	"github.com/MichaelShi-san/subscription-test/internal/repository"

	"github.com/MichaelShi-san/subscription-test/internal/model"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(r *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: r}
}

func (s *SubscriptionService) Create(ctx context.Context, sub *model.Subscription) error {
	return s.repo.Create(ctx, sub)
}

func (s *SubscriptionService) List(ctx context.Context) ([]model.Subscription, error) {
	return s.repo.List(ctx)
}

func (s *SubscriptionService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *SubscriptionService) Update(ctx context.Context, sub *model.Subscription) error {
	return s.repo.Update(ctx, sub)
}

func (s *SubscriptionService) TotalCost(
	ctx context.Context,
	userID, service string,
	from, to time.Time,
) (int, error) {
	return s.repo.TotalCost(ctx, userID, service, from, to)
}
