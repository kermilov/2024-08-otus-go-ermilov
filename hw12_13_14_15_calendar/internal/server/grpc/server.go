package internalgrpc

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	addr   string
	logger Logger
}

type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context,
		id string, title string, datetime time.Time, duration *time.Duration, userid int64) (*storage.Event, error)
	UpdateEvent(ctx context.Context,
		id string, title string, datetime time.Time, duration *time.Duration, userid int64) error
	DeleteEvent(ctx context.Context, id string) error
	FindEventByDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	FindEventByWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	FindEventByMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
	FindEventByID(ctx context.Context, id string) (storage.Event, error)
}

func NewServer(logger Logger, app Application, addr string) *Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingMiddleware(),
		),
	)
	pb.RegisterEventServiceServer(grpcServer, &Service{app: app})
	return &Server{
		Server: grpcServer,
		addr:   addr,
		logger: logger,
	}
}

func (s *Server) Start(_ context.Context) error {
	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	err = s.Server.Serve(lsn)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.Server.GracefulStop()
	return nil
}

type Service struct {
	pb.UnimplementedEventServiceServer
	app Application
}

// CreateEvent implements pb.EventServiceServer.
func (s *Service) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	// Указание формата времени
	layout := time.RFC3339
	// Преобразование строки в time.Time
	eventTime, err := time.Parse(layout, req.Event.Datetime)
	if err != nil {
		return nil, err
	}
	// Преобразование строки в time.Duration
	var eventDuration *time.Duration
	if req.Event.Duration != "" {
		parseDuration, err := time.ParseDuration(req.Event.Duration)
		if err != nil {
			return nil, err
		}
		eventDuration = &parseDuration
	}
	event, err := s.app.CreateEvent(ctx,
		req.Event.Id,
		req.Event.Title,
		eventTime,
		eventDuration,
		req.Event.Userid,
	)
	if err != nil {
		if errors.Is(err, storage.ErrBusiness) {
			resp := &pb.CreateEventResponse{
				Result: &pb.CreateEventResponse_Error{
					Error: err.Error(),
				},
			}
			return resp, nil
		}
		return nil, err
	}
	protoEvent := &pb.Event{
		Id:       event.ID,
		Title:    event.Title,
		Datetime: event.DateTime.String(),
		Duration: event.Duration.String(),
		Userid:   event.UserID,
	}
	resp := &pb.CreateEventResponse{
		Result: &pb.CreateEventResponse_Event{
			Event: protoEvent,
		},
	}
	return resp, nil
}

// UpdateEvent implements pb.EventServiceServer.
func (s *Service) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	// Указание формата времени
	layout := time.RFC3339
	// Преобразование строки в time.Time
	eventTime, err := time.Parse(layout, req.Event.Datetime)
	if err != nil {
		return nil, err
	}
	// Преобразование строки в time.Duration
	var eventDuration *time.Duration
	if req.Event.Duration != "" {
		parseDuration, err := time.ParseDuration(req.Event.Duration)
		if err != nil {
			return nil, err
		}
		eventDuration = &parseDuration
	}
	err = s.app.UpdateEvent(ctx,
		req.Event.Id,
		req.Event.Title,
		eventTime,
		eventDuration,
		req.Event.Userid,
	)
	if err != nil {
		if errors.Is(err, storage.ErrBusiness) {
			resp := &pb.UpdateEventResponse{
				Error: err.Error(),
			}
			return resp, nil
		}
		return nil, err
	}
	return &pb.UpdateEventResponse{}, nil
}

// DeleteEvent implements pb.EventServiceServer.
func (s *Service) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	err := s.app.DeleteEvent(ctx, req.Id)
	if err != nil {
		if errors.Is(err, storage.ErrBusiness) {
			resp := &pb.DeleteEventResponse{
				Error: err.Error(),
			}
			return resp, nil
		}
		return nil, err
	}
	return &pb.DeleteEventResponse{}, nil
}

// GetEventByID implements pb.EventServiceServer.
func (s *Service) GetEventByID(ctx context.Context, req *pb.GetEventByIDRequest) (*pb.GetEventByIDResponse, error) {
	event, err := s.app.FindEventByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, storage.ErrBusiness) {
			resp := &pb.GetEventByIDResponse{
				// Error: err.Error(),
			}
			return resp, nil
		}
		return nil, err
	}
	protoEvent := &pb.Event{
		Id:       event.ID,
		Title:    event.Title,
		Datetime: event.DateTime.String(),
		Duration: event.Duration.String(),
		Userid:   event.UserID,
	}
	resp := &pb.GetEventByIDResponse{
		Event: protoEvent,
	}
	return resp, nil
}

// GetEventsByDay implements pb.EventServiceServer.
func (s *Service) GetEventsByDay(ctx context.Context, req *pb.GetEventsByDateRequest) (*pb.GetEventsResponse, error) {
	// Указание формата времени
	layout := time.RFC3339
	// Преобразование строки в time.Time
	searchDate, err := time.Parse(layout, req.Date)
	if err != nil {
		return nil, err
	}
	events, err := s.app.FindEventByDay(ctx, searchDate)
	if err != nil {
		return nil, err
	}
	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvent := &pb.Event{
			Id:       event.ID,
			Title:    event.Title,
			Datetime: event.DateTime.String(),
			Duration: event.Duration.String(),
			Userid:   event.UserID,
		}
		protoEvents[i] = protoEvent
	}
	return &pb.GetEventsResponse{Events: protoEvents}, nil
}

// GetEventsByMonth implements pb.EventServiceServer.
func (s *Service) GetEventsByMonth(ctx context.Context, req *pb.GetEventsByDateRequest) (*pb.GetEventsResponse, error) {
	// Указание формата времени
	layout := time.RFC3339
	// Преобразование строки в time.Time
	searchDate, err := time.Parse(layout, req.Date)
	if err != nil {
		return nil, err
	}
	events, err := s.app.FindEventByMonth(ctx, searchDate)
	if err != nil {
		return nil, err
	}
	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvent := &pb.Event{
			Id:       event.ID,
			Title:    event.Title,
			Datetime: event.DateTime.String(),
			Duration: event.Duration.String(),
			Userid:   event.UserID,
		}
		protoEvents[i] = protoEvent
	}
	return &pb.GetEventsResponse{Events: protoEvents}, nil
}

// GetEventsByWeek implements pb.EventServiceServer.
func (s *Service) GetEventsByWeek(ctx context.Context, req *pb.GetEventsByDateRequest) (*pb.GetEventsResponse, error) {
	// Указание формата времени
	layout := time.RFC3339
	// Преобразование строки в time.Time
	searchDate, err := time.Parse(layout, req.Date)
	if err != nil {
		return nil, err
	}
	events, err := s.app.FindEventByWeek(ctx, searchDate)
	if err != nil {
		return nil, err
	}
	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvent := &pb.Event{
			Id:       event.ID,
			Title:    event.Title,
			Datetime: event.DateTime.String(),
			Duration: event.Duration.String(),
			Userid:   event.UserID,
		}
		protoEvents[i] = protoEvent
	}
	return &pb.GetEventsResponse{Events: protoEvents}, nil
}
