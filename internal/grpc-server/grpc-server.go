package grpcserver

import (
	"context"
	"fmt"
	"net"

	customerrors "github.com/GrishaSkurikhin/OzonTestTask/internal/custom-errors"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/service/shortlinks"
	shortlinksv1 "github.com/GrishaSkurikhin/OzonTestTask/protos/gen/go/shortlinks"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Storage interface {
	shortlinks.URLSaver
	shortlinks.URLGetter
}

type Shortlinks interface {
	GetURL(ctx context.Context, shortURL string, getter shortlinks.URLGetter) (string, error)
	SaveURL(ctx context.Context, longURL string, host string, saver shortlinks.URLSaver) (string, error)
}

type serverAPI struct {
	shortlinksv1.UnimplementedShortlinksServer
	shortlinks   Shortlinks
	storage      Storage
	shortURLHost string
}

func New(log *zerolog.Logger, shortURLHost string, strg Storage) *grpc.Server {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Fatal().Msg("Recovered from panic")
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))

	shortlinksService := &shortlinks.ShortlinksService{}
	Register(gRPCServer, shortlinksService, strg, shortURLHost)

	return gRPCServer
}

func InterceptorLogger(l *zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...interface{}) {
		l.Log().Fields(fields).Msg(msg)
	})
}

func Run(gRPCServer *grpc.Server, port int) error {
	const op = "grpcserver.Run"
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("%s: failed to listen: %v", op, err)
	}

	if err := gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: failed to serve: %v", op, err)
	}

	return nil
}

func Register(gRPCServer *grpc.Server, shortlinks Shortlinks, strg Storage, shortURLHost string) {
	shortlinksv1.RegisterShortlinksServer(gRPCServer, &serverAPI{
		shortlinks: shortlinks, 
		storage: strg, 
		shortURLHost: shortURLHost,
	})
}

func (s *serverAPI) GetURL(ctx context.Context, in *shortlinksv1.GetURLRequest) (*shortlinksv1.GetURLResponse, error) {
	if in.ShortURL == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}

	longURL, err := s.shortlinks.GetURL(context.Background(), in.ShortURL, s.storage)
	if err != nil {
		switch err.(type) {
		case customerrors.URLNotFound:
			return nil, status.Error(codes.InvalidArgument, "url not found")
		case customerrors.WrongURL:
			return nil, status.Error(codes.InvalidArgument, "wrong url")
		default:
			return nil, status.Error(codes.InvalidArgument, "internal error")
		}
	}

	return &shortlinksv1.GetURLResponse{LongURL: longURL}, nil
}

func (s *serverAPI) SaveURL(ctx context.Context, in *shortlinksv1.SaveURLRequest) (*shortlinksv1.SaveURLResponse, error) {
	if in.LongURL == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}

	shortURL, err := s.shortlinks.SaveURL(context.Background(), in.LongURL, s.shortURLHost, s.storage)
	if err != nil {
		switch err.(type) {
		case customerrors.WrongURL:
			return nil, status.Error(codes.InvalidArgument, "wrong url")
		default:
			return nil, status.Error(codes.InvalidArgument, "internal error")
		}
	}

	return &shortlinksv1.SaveURLResponse{ShortURL: shortURL}, nil
}
