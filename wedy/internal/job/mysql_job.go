package job

import (
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/service"
	"GoProj/wedy/pkg/logger"
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"time"
)

type Scheduler struct {
	dbTimeout time.Duration
	svc       service.CronJobService
	executors map[string]Executor
	l         logger.LoggerV1
	limiter   *semaphore.Weighted
}

func NewScheduler(l logger.LoggerV1, svc service.CronJobService) *Scheduler {
	return &Scheduler{l: l, svc: svc, dbTimeout: time.Second, limiter: semaphore.NewWeighted(100), executors: map[string]Executor{}}
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, j domian.Job) error
}

func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{
		funcs: map[string]func(ctx context.Context, j domian.Job) error{},
	}
}

func (l *LocalFuncExecutor) Name() string {
	return "LocalFuncExecutor"
}
func (l *LocalFuncExecutor) RegisterFunc(name string, fn func(ctx context.Context, j domian.Job) error) {
	l.funcs[name] = fn
}
func (l *LocalFuncExecutor) Exec(ctx context.Context, j domian.Job) error {
	fn, ok := l.funcs[j.Name]
	if !ok {
		return fmt.Errorf("unknown job name: %s", j.Name)
	}
	return fn(ctx, j)
}

// 执行任务器
type Executor interface {
	Name() string
	//全局控制 Exec的实现需要正确处理ctx超时或者取消
	Exec(ctx context.Context, j domian.Job) error
}

func (s *Scheduler) RegisterExecutor(executor Executor) {
	s.executors[executor.Name()] = executor
}
func (s *Scheduler) Schedule(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			return
		}
		dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		j, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			continue
		}
		//调度执行j
		exec, ok := s.executors[j.Executor]
		if !ok {
			s.l.Error("cannot find executor", logger.Int64("job id", j.Id), logger.String("executor", j.Executor))
			continue
		}
		go func() {
			defer func() {
				s.limiter.Release(1)
				//执行完毕释放
				j.Cancel()
			}()
			err1 := exec.Exec(ctx, j)
			if err1 != nil {
				s.l.Error("exec job failed", logger.Int64("job id", j.Id), logger.String("executor", j.Executor))
				return
			}
			err1 = s.svc.ResetNextTime(ctx, j)
			if err1 != nil {
				s.l.Error("Reset Next Time failed", logger.Int64("job id", j.Id), logger.String("executor", j.Executor))
				return
			}
		}()
	}
}
