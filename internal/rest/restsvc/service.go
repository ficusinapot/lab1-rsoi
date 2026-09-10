package restsvc

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/ficusinapot/ds/internal/observability/metrics"
	restmetrics "github.com/ficusinapot/ds/internal/rest/metrics"
	"github.com/ficusinapot/ds/internal/rest/restsvc/resources/persons"
	"github.com/ficusinapot/ds/internal/rest/restsvc/resources/status"

	"github.com/joomcode/errorx"
)

type Service struct {
	config                  ServiceConfig
	apiServer               *http.Server
	metricsServer           *http.Server
	apiListener             net.Listener
	metricsListener         net.Listener
	terminationNotification chan struct{}
	stopOnce                sync.Once
	logger                  *slog.Logger
}

func NewService(
	appConfig AppInfo,
	serviceConfig ServiceConfig,
	metricsConfig MetricsConfig,
	personHandler *persons.Handler,
	statusHandler *status.Handler,
	logger *slog.Logger,
	metricsRegistry *metrics.Registry,
) *Service {
	ch := make(chan struct{})
	close(ch)

	return &Service{
		config:                  serviceConfig,
		apiServer:               newHTTPServer(appConfig, serviceConfig, personHandler, statusHandler, logger, metricsRegistry),
		metricsServer:           newMetricsServer(metricsConfig, metricsRegistry),
		apiListener:             nil,
		metricsListener:         nil,
		terminationNotification: ch,
		stopOnce:                sync.Once{},
		logger:                  logger.With("service", "rest"),
	}
}

func (s *Service) Run(ctx context.Context) error {
	_ = ctx

	if s.apiListener != nil {
		return errorx.IllegalState.New("REST API listener already exists")
	}

	apiListener, err := net.Listen("tcp", s.apiServer.Addr)
	if err != nil {
		return errorx.ExternalError.Wrap(err, "create REST API listener")
	}
	s.apiListener = apiListener
	s.terminationNotification = make(chan struct{})

	if s.metricsServer != nil {
		metricsListener, err := net.Listen("tcp", s.metricsServer.Addr)
		if err != nil {
			_ = apiListener.Close()
			s.apiListener = nil
			return errorx.ExternalError.Wrap(err, "create metrics listener")
		}
		s.metricsListener = metricsListener
	}

	if s.metricsServer != nil {
		go s.serve("metrics", s.metricsServer, s.metricsListener, nil)
		s.logger.Info("metrics service started", "addr", s.metricsServer.Addr)
	}

	go s.serve("REST API", s.apiServer, s.apiListener, s.terminationNotification)
	s.logger.Info("REST API service started", "addr", s.apiServer.Addr)
	restmetrics.ServiceStartsTotal.Inc()

	return nil
}

func (s *Service) TerminationNotification() <-chan struct{} {
	return s.terminationNotification
}

func (s *Service) Stop(ctx context.Context) error {
	var result error

	s.stopOnce.Do(func() {
		shutdownCtx, cancel := context.WithTimeout(ctx, s.config.Timeout.Shutdown)
		defer cancel()

		if s.metricsServer != nil {
			if err := s.metricsServer.Shutdown(shutdownCtx); err != nil {
				result = errorx.DecorateMany("stop REST API service", result, err)
			}
			if s.metricsListener != nil {
				if err := s.metricsListener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
					result = errorx.DecorateMany("stop REST API service", result, err)
				}
			}
			s.metricsListener = nil
		}

		if err := s.apiServer.Shutdown(shutdownCtx); err != nil {
			result = errorx.DecorateMany("stop REST API service", result, err)
		}
		if s.apiListener != nil {
			if err := s.apiListener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
				result = errorx.DecorateMany("stop REST API service", result, err)
			}
		}
		s.apiListener = nil

		select {
		case <-s.terminationNotification:
		case <-shutdownCtx.Done():
			result = errorx.DecorateMany("stop REST API service", result, shutdownCtx.Err())
		}
	})

	return result //nolint:wrapcheck // Result is already accumulated with contextual shutdown errors.
}

func (s *Service) Dispose() {
	if s.metricsServer != nil {
		if err := s.metricsServer.Close(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Warn("dispose metrics service", "error", err)
		}
	}
	if s.metricsListener != nil {
		if err := s.metricsListener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			s.logger.Warn("dispose metrics listener", "error", err)
		}
		s.metricsListener = nil
	}

	if err := s.apiServer.Close(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Warn("dispose REST API service", "error", err)
	}
	if s.apiListener != nil {
		if err := s.apiListener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			s.logger.Warn("dispose REST API listener", "error", err)
		}
		s.apiListener = nil
	}
}

func (s *Service) serve(name string, server *http.Server, listener net.Listener, done chan struct{}) {
	if done != nil {
		defer close(done)
	}

	if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, net.ErrClosed) {
		s.logger.Error(name+" service error", "error", err)
	}
	s.logger.Info(name + " service stopped")
}
