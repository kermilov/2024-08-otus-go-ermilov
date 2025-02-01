package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/server"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/server/http/dto"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	http.Server
	logger  server.Logger
	service *Service
}

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total number of API requests",
		},
		[]string{"method", "endpoint", "status"},
	)
	responseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_response_duration_seconds",
			Help:    "Response duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	eventsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_created_total",
			Help: "Total number of created events",
		},
	)
	eventsUpdated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_updated_total",
			Help: "Total number of updated events",
		},
	)
	eventsDeleted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_deleted_total",
			Help: "Total number of deleted events",
		},
	)
	backgroundTasksStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "background_tasks_status",
			Help: "Status of background tasks",
		},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal, responseDuration, eventsCreated, eventsUpdated, eventsDeleted, backgroundTasksStatus)
}

func NewServer(logger server.Logger, app server.Application, addr string) *Server {
	service := &Service{
		app: app,
	}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /hello", service.Hello)
	serveMux.HandleFunc("POST /event", service.CreateEvent)
	serveMux.HandleFunc("PUT /event", service.UpdateEvent)
	serveMux.HandleFunc("DELETE /event/{id}", service.DeleteEvent)
	serveMux.HandleFunc("GET /event/{id}", service.GetEventByID)
	serveMux.HandleFunc("GET /events", service.GetEvents)
	serveMux.Handle("/metrics", promhttp.Handler())

	return &Server{
		Server: http.Server{
			Addr:         addr,
			Handler:      loggingMiddleware(serveMux),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		logger:  logger,
		service: service,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.service.ctx = ctx
	err := s.Server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.Server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

type Service struct {
	app server.Application
	ctx context.Context
}

func (s *Service) Hello(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	name := args.Get("name")
	results, err := s.hello(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

func (s *Service) hello(name string) (string, error) {
	if name != "" {
		if name == "world" {
			return "", errors.New("error")
		}
		return `{"hello":"` + name + `"}`, nil
	}
	return `{"hello":"world"}`, nil
}

func (s *Service) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event dto.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()
	// Преобразование строки в time.Time
	eventTime, err := time.Parse(server.Layout, event.DateTime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// Преобразование строки в time.Duration
	var eventDuration *time.Duration
	if event.Duration != "" {
		parseDuration, err := time.ParseDuration(event.Duration)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		eventDuration = &parseDuration
	}
	// Преобразование строки в time.Duration
	var eventNotificationDuration *time.Duration
	if event.NotificationDuration != "" {
		parseDuration, err := time.ParseDuration(event.NotificationDuration)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		eventNotificationDuration = &parseDuration
	}
	result, err := s.app.CreateEvent(s.ctx,
		event.ID,
		event.Title,
		eventTime,
		eventDuration,
		event.UserID,
		eventNotificationDuration,
	)
	if err != nil {
		if errors.Is(err, storage.ErrBusiness) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}
	eventsCreated.Inc()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp := s.mapToDtoEvent(*result)
	json.NewEncoder(w).Encode(resp)
}

func (s *Service) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	var event dto.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()
	event.ID = idString
	// Преобразование строки в time.Time
	eventTime, err := time.Parse(server.Layout, event.DateTime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// Преобразование строки в time.Duration
	var eventDuration *time.Duration
	if event.Duration != "" {
		parseDuration, err := time.ParseDuration(event.Duration)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		eventDuration = &parseDuration
	}
	// Преобразование строки в time.Duration
	var eventNotificationDuration *time.Duration
	if event.NotificationDuration != "" {
		parseDuration, err := time.ParseDuration(event.NotificationDuration)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		eventNotificationDuration = &parseDuration
	}
	err = s.app.UpdateEvent(s.ctx,
		event.ID,
		event.Title,
		eventTime,
		eventDuration,
		event.UserID,
		eventNotificationDuration,
	)
	if err != nil {
		if errors.Is(err, storage.ErrBusiness) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}
	eventsUpdated.Inc()
	w.WriteHeader(http.StatusOK)
}

func (s *Service) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	err := s.app.DeleteEvent(s.ctx, idString)
	if err != nil {
		if errors.Is(err, storage.ErrBusiness) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}
	eventsDeleted.Inc()
	w.WriteHeader(http.StatusOK)
}

func (s *Service) GetEventByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	result, err := s.app.FindEventByID(s.ctx, idString)
	if err != nil {
		if errors.Is(err, storage.ErrBusiness) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp := s.mapToDtoEvent(result)
	json.NewEncoder(w).Encode(resp)
}

func (s *Service) GetEvents(w http.ResponseWriter, r *http.Request) {
	results, err := s.getEvents(r)
	if err != nil {
		if errors.Is(err, storage.ErrBusiness) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp := make([]dto.Event, len(results))
	for i, result := range results {
		resp[i] = s.mapToDtoEvent(result)
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *Service) getEvents(r *http.Request) ([]storage.Event, error) {
	args := r.URL.Query()
	day := args.Get("day")
	if day != "" {
		// Преобразование строки в time.Time
		date, err := time.Parse(server.Layout, day)
		if err != nil {
			return nil, err
		}
		return s.app.FindEventByDay(s.ctx, date)
	}
	week := args.Get("week")
	if week != "" {
		// Преобразование строки в time.Time
		date, err := time.Parse(server.Layout, week)
		if err != nil {
			return nil, err
		}
		return s.app.FindEventByWeek(s.ctx, date)
	}
	month := args.Get("month")
	if month != "" {
		// Преобразование строки в time.Time
		date, err := time.Parse(server.Layout, month)
		if err != nil {
			return nil, err
		}
		return s.app.FindEventByMonth(s.ctx, date)
	}
	return nil, errors.ErrUnsupported
}

func (*Service) mapToDtoEvent(result storage.Event) dto.Event {
	return dto.Event{
		ID:                   result.ID,
		Title:                result.Title,
		DateTime:             result.DateTime.Format(server.Layout),
		Duration:             result.Duration.String(),
		UserID:               result.UserID,
		NotificationDuration: result.NotificationDuration.String(),
	}
}
