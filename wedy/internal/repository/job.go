package repository

import (
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/repository/dao"
	"context"
	"time"
)

type JobRepository interface {
	Preempt(ctx context.Context) (domian.Job, error)
	Release(ctx context.Context, id int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, time time.Time) error
}

type PreemptJobRepository struct {
	dao dao.JobDAO
}

func (p *PreemptJobRepository) Preempt(ctx context.Context) (domian.Job, error) {
	j, err := p.dao.Preempt(ctx)
	return domian.Job{
		Id:         j.Id,
		Expression: j.Expression,
		Executor:   j.Executor,
		Name:       j.Name,
	}, err
}
func (p *PreemptJobRepository) UpdateNextTime(ctx context.Context, id int64, time time.Time) error {
	return p.dao.UpdateNextTime(ctx, id, time)
}
func (p *PreemptJobRepository) UpdateUtime(ctx context.Context, id int64) error {
	return p.dao.UpdateUtime(ctx, id)
}
func (p *PreemptJobRepository) Release(ctx context.Context, id int64) error {
	return p.dao.Release(ctx, id)
}
