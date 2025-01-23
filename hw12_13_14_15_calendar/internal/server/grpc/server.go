package internalgrpc

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/server"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	addr   string
	logger server.Logger
}

func NewServer(logger server.Logger, app server.Application, addr string) *Server {
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
	app server.Application
}

// CreateEvent implements pb.EventServiceServer.
func (s *Service) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	// Преобразование строки в time.Time
	eventTime, err := time.Parse(server.Layout, req.Event.Datetime)
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
	// Преобразование строки в time.Duration
	var eventNotificationDuration *time.Duration
	if req.Event.Notificationduration != "" {
		parseDuration, err := time.ParseDuration(req.Event.Notificationduration)
		if err != nil {
			return nil, err
		}
		eventNotificationDuration = &parseDuration
	}
	event, err := s.app.CreateEvent(ctx,
		req.Event.Id,
		req.Event.Title,
		eventTime,
		eventDuration,
		req.Event.Userid,
		eventNotificationDuration,
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
	resp := &pb.CreateEventResponse{
		Result: &pb.CreateEventResponse_Event{
			Event: s.mapToProtobufEvent(*event),
		},
	}
	return resp, nil
}

// UpdateEvent implements pb.EventServiceServer.
func (s *Service) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	// Преобразование строки в time.Time
	eventTime, err := time.Parse(server.Layout, req.Event.Datetime)
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
	// Преобразование строки в time.Duration
	var eventNotificationDuration *time.Duration
	if req.Event.Notificationduration != "" {
		parseDuration, err := time.ParseDuration(req.Event.Notificationduration)
		if err != nil {
			return nil, err
		}
		eventNotificationDuration = &parseDuration
	}
	err = s.app.UpdateEvent(ctx,
		req.Event.Id,
		req.Event.Title,
		eventTime,
		eventDuration,
		req.Event.Userid,
		eventNotificationDuration,
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
	resp := &pb.GetEventByIDResponse{
		Event: s.mapToProtobufEvent(event),
	}
	return resp, nil
}

type findEventFunc func(ctx context.Context, date time.Time) ([]storage.Event, error)

// GetEventsByDay implements pb.EventServiceServer.
func (s *Service) GetEventsByDay(ctx context.Context, req *pb.GetEventsByDateRequest) (*pb.GetEventsResponse, error) {
	return s.getEvents(ctx, req, s.app.FindEventByDay)
}

// GetEventsByMonth implements pb.EventServiceServer.
func (s *Service) GetEventsByMonth(ctx context.Context, req *pb.GetEventsByDateRequest) (*pb.GetEventsResponse, error) {
	return s.getEvents(ctx, req, s.app.FindEventByMonth)
}

// GetEventsByWeek implements pb.EventServiceServer.
func (s *Service) GetEventsByWeek(ctx context.Context, req *pb.GetEventsByDateRequest) (*pb.GetEventsResponse, error) {
	return s.getEvents(ctx, req, s.app.FindEventByWeek)
}

func (s *Service) getEvents(
	ctx context.Context, req *pb.GetEventsByDateRequest, method findEventFunc,
) (
	*pb.GetEventsResponse, error,
) {
	// Преобразование строки в time.Time
	searchDate, err := time.Parse(server.Layout, req.Date)
	if err != nil {
		return nil, err
	}
	events, err := method(ctx, searchDate)
	if err != nil {
		return nil, err
	}
	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvents[i] = s.mapToProtobufEvent(event)
	}
	return &pb.GetEventsResponse{Events: protoEvents}, nil
}

func (*Service) mapToProtobufEvent(event storage.Event) *pb.Event {
	return &pb.Event{
		Id:                   event.ID,
		Title:                event.Title,
		Datetime:             event.DateTime.Format(server.Layout),
		Duration:             event.Duration.String(),
		Userid:               event.UserID,
		Notificationduration: event.NotificationDuration.String(),
	}
}
