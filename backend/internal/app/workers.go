package app

import (
	"context"
	"log"
	"time"

	learncommand "backend/internal/learn/application/command"
)

type LearnSyncWorker struct {
	sync       learncommand.SyncAll
	link       learncommand.LinkUsers
	apiEnabled bool
	interval   time.Duration
}

func NewLearnSyncWorker(sync learncommand.SyncAll, link learncommand.LinkUsers, apiEnabled bool, interval time.Duration) *LearnSyncWorker {
	return &LearnSyncWorker{sync: sync, link: link, apiEnabled: apiEnabled, interval: interval}
}

func (w *LearnSyncWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		w.tick(ctx)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w *LearnSyncWorker) tick(ctx context.Context) {
	if w.apiEnabled {
		if err := w.sync.SyncAll(ctx); err != nil {
			log.Printf("SyncAll ошибка: %v", err)
		}
	}
	if err := w.link.LinkUsers(ctx); err != nil {
		log.Printf("LinkUsers ошибка: %v", err)
	}
}