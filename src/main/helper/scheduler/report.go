package scheduler

import (
	"context"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/observer"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/exam_service"
	"log"
	"time"
)

func init() {
	go func() {
		StartSessionReportWatcher(context.Background(), time.Minute*30)
	}()
}

func StartSessionReportWatcher(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("StartSessionReportWatcher")
			sessionService := exam_service.NewExamSessionService()
			sessionService.GenerateReportSession()
			observer.Trigger(model.EventExamSessionReportChanged)
			log.Println("EndSessionReportWatcher")
		case <-ctx.Done():
			log.Println("Watcher stopped")
			return
		}
	}
}
