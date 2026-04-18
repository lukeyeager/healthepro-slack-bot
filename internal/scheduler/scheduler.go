package scheduler

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/lukeyeager/school-lunch-schedule/internal/config"
	"github.com/lukeyeager/school-lunch-schedule/internal/healthepro"
	"github.com/lukeyeager/school-lunch-schedule/internal/metrics"
	"github.com/lukeyeager/school-lunch-schedule/internal/store"
	"github.com/lukeyeager/school-lunch-schedule/internal/week"
)

// Scheduler polls the Health-e Pro API on a cron schedule and persists
// the current display week's menu to the store.
type Scheduler struct {
	cfg  *config.Config
	hep  *healthepro.Client
	db   *store.Store
	m    *metrics.Metrics
	cron *cron.Cron
	loc  *time.Location
}

// New creates a Scheduler. The timezone in cfg is used for cron scheduling
// and week boundary calculations.
func New(cfg *config.Config, hep *healthepro.Client, db *store.Store, m *metrics.Metrics) (*Scheduler, error) {
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone %q: %w", cfg.Timezone, err)
	}
	return &Scheduler{
		cfg:  cfg,
		hep:  hep,
		db:   db,
		m:    m,
		cron: cron.New(cron.WithLocation(loc)),
		loc:  loc,
	}, nil
}

// Start registers the fetch cron job, triggers an immediate first fetch,
// and starts the scheduler.
func (s *Scheduler) Start() error {
	if _, err := s.cron.AddFunc(s.cfg.FetchCron, s.fetchWeek); err != nil {
		return fmt.Errorf("adding fetch cron %q: %w", s.cfg.FetchCron, err)
	}
	s.cron.Start()
	go s.fetchWeek() // populate data immediately on startup
	return nil
}

// Stop gracefully stops the scheduler.
func (s *Scheduler) Stop() {
	s.cron.Stop()
}

// fetchWeek fetches Mon–Fri of the current display week and upserts each day
// into the store. Fetch failures per day are logged and counted; the loop
// continues so a single bad day does not block the others.
func (s *Scheduler) fetchWeek() {
	now := time.Now().In(s.loc)
	monday := week.DisplayMonday(now)
	slog.Info("fetch run", "week_of", monday.Format("2006-01-02"))

	for i := 0; i < 5; i++ {
		date := monday.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")

		entry, err := s.hep.FetchMenu(date)
		if err != nil {
			slog.Error("fetch failed", "date", dateStr, "err", err)
			s.m.FetchFailures.Inc()
			continue
		}
		if entry == nil {
			slog.Info("no menu data", "date", dateStr)
			continue
		}

		if err := s.db.Upsert(dateStr, entry); err != nil {
			slog.Error("store upsert failed", "date", dateStr, "err", err)
		}
	}
}
