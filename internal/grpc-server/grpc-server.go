package grpcserver

import (
	"context"

	"github.com/GrishaSkurikhin/OzonTestTask/internal/config"
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

func New(cfg config.Server, log *zerolog.Logger, strg Storage, shortlinksService Shortlinks) *grpc.Server {
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

	Register(gRPCServer, shortlinksService)

	return gRPCServer
}

func InterceptorLogger(l *zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...interface{}) {
		l.Log().Fields(fields).Msg(msg)
	})
}

type Shortlinks interface {
	GetURL(ctx context.Context, shortURL string) (string, error)
	SaveURL(ctx context.Context, longURL string) (string, error)
}

type serverAPI struct {
	shortlinksv1.UnimplementedShortlinksServer
	shortlinks Shortlinks
}

func Register(gRPCServer *grpc.Server, shortlinks Shortlinks) {
	shortlinksv1.RegisterShortlinksServer(gRPCServer, &serverAPI{shortlinks: shortlinks})
}

func (s *serverAPI) GetURL(ctx context.Context, in *shortlinksv1.GetURLRequest) (*shortlinksv1.GetURLResponse, error) {
	panic("implement me")
}

func (s *serverAPI) SaveURL(ctx context.Context, in *shortlinksv1.SaveURLRequest) (*shortlinksv1.SaveURLResponse, error) {
	panic("implement me")
}
