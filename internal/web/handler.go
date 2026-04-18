package web

import (
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/lukeyeager/school-lunch-schedule/internal/store"
	"github.com/lukeyeager/school-lunch-schedule/internal/week"
)

// Handler serves the weekly menu page.
type Handler struct {
	store *store.Store
	loc   *time.Location
	tmpl  *template.Template
}

// dayData holds the data for one column of the weekly grid.
type dayData struct {
	Date   time.Time
	Record *store.MenuRecord // nil when no data has been fetched yet
}

// pageData is passed to the HTML template.
type pageData struct {
	WeekLabel string
	Days      []dayData
}

// New creates a Handler.
func New(db *store.Store, loc *time.Location) *Handler {
	return &Handler{
		store: db,
		loc:   loc,
		tmpl:  template.Must(template.New("menu").Parse(htmlTemplate)),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now().In(h.loc)
	monday := week.DisplayMonday(now)
	friday := monday.AddDate(0, 0, 4)

	days := make([]dayData, 5)
	for i := range days {
		date := monday.AddDate(0, 0, i)
		rec, err := h.store.Get(date.Format("2006-01-02"))
		if err != nil {
			slog.Warn("store get failed", "date", date.Format("2006-01-02"), "err", err)
		}
		days[i] = dayData{Date: date, Record: rec}
	}

	data := pageData{
		WeekLabel: monday.Format("Jan 2") + "–" + friday.Format("Jan 2, 2006"),
		Days:      days,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.Execute(w, data); err != nil {
		slog.Error("template execute failed", "err", err)
	}
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>School Lunch · {{.WeekLabel}}</title>
  <style>
    *, *::before, *::after { box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      background: #f4f4f4;
      color: #333;
      margin: 0;
      padding: 2rem 1rem 3rem;
    }
    header { text-align: center; margin-bottom: 2rem; }
    h1 { margin: 0 0 .3rem; font-size: 1.6rem; color: #222; }
    .week-label { color: #888; font-size: .95rem; }
    .grid {
      display: grid;
      grid-template-columns: repeat(5, 1fr);
      gap: 1rem;
      max-width: 960px;
      margin: 0 auto;
    }
    .card {
      background: #fff;
      border-radius: 10px;
      padding: 1rem 1rem 1.25rem;
      box-shadow: 0 1px 4px rgba(0,0,0,.08);
    }
    .day-name {
      font-size: .75rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: .06em;
      color: #aaa;
    }
    .day-date {
      font-size: 1.1rem;
      font-weight: 600;
      margin-bottom: .75rem;
    }
    .category {
      margin-top: .75rem;
      font-size: .75rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: .04em;
      color: #2e7d32;
    }
    .recipe { font-size: .85rem; line-height: 1.45; color: #555; }
    .no-data { font-size: .85rem; color: #ccc; font-style: italic; margin-top: .5rem; }
    footer { text-align: center; margin-top: 2.5rem; font-size: .75rem; color: #bbb; }
    @media (max-width: 640px) { .grid { grid-template-columns: 1fr; } }
  </style>
</head>
<body>
<header>
  <h1>School Lunch</h1>
  <div class="week-label">Week of {{.WeekLabel}}</div>
</header>
<div class="grid">
  {{- range .Days}}
  <div class="card">
    <div class="day-name">{{.Date.Format "Mon"}}</div>
    <div class="day-date">{{.Date.Format "Jan 2"}}</div>
    {{- if .Record}}
      {{- range .Record.Items}}
        {{- if eq .Type "category"}}
    <div class="category">{{.Name}}</div>
        {{- else}}
    <div class="recipe">{{.Name}}</div>
        {{- end}}
      {{- end}}
    {{- else}}
    <div class="no-data">No data yet</div>
    {{- end}}
  </div>
  {{- end}}
</div>
<footer>Data from Health-e Pro &middot; Updated hourly</footer>
</body>
</html>`
