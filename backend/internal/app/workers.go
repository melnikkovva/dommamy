package app

import (
	"context"
	"log"
	"time"

	"backend/internal/learn"
)

type LearnSyncWorker struct {
	sync       *learn.SyncService
	apiEnabled bool
	interval   time.Duration
}

func NewLearnSyncWorker(s *learn.SyncService, apiEnabled bool, interval time.Duration) *LearnSyncWorker {
	return &LearnSyncWorker{sync: s, apiEnabled: apiEnabled, interval: interval}
}

func (w *LearnSyncWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		w.tick()
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w *LearnSyncWorker) tick() {
	if w.apiEnabled {
		if err := w.sync.SyncAll(); err != nil {
			log.Printf("SyncAll ошибка: %v", err)
		}
	}
	if err := w.sync.SyncUsers(); err != nil {
		log.Printf("SyncUsers ошибка: %v", err)
	}
}