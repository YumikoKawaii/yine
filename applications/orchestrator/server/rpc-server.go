package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/YumikoKawaii/rpc.com/protobuf/orchestrator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// DefaultConfig return a default server config
func DefaultConfig() Config {
	return Config{
		GRPC: ServerListen{
			Host: "0.0.0.0",
			Port: 10443,
		},
		HTTP: ServerListen{
			Host: "0.0.0.0",
			Port: 10080,
		},
	}
}

// Config hold http/grpc server config
type Config struct {
	GRPC ServerListen `json:"grpc" mapstructure:"grpc" yaml:"grpc"`
	HTTP ServerListen `json:"http" mapstructure:"http" yaml:"http"`
}

// ServerListen config for host/port socket listener
type ServerListen struct {
	Host string `json:"host" mapstructure:"host" yaml:"host"`
	Port int    `json:"port" mapstructure:"port" yaml:"port"`
}

// String return socket listen DSN
func (l ServerListen) String() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

// Server structure
type Server struct {
	gRPC *grpc.Server
	mux  *runtime.ServeMux
	cfg  Config
}

func NewServer(cfg Config, opt ...grpc.ServerOption) *Server {
	return &Server{
		gRPC: grpc.NewServer(opt...),
		mux: runtime.NewServeMux(
			//gatewayopt.ProtoJSONMarshaler(),
		),
		cfg: cfg,
	}
}

func (s *Server) Register(grpcServer ...interface{}) error {
	for _, srv := range grpcServer {
		isNotMatched := true
		if orchestratorSv, ok := srv.(api.OrchestratorServer); ok {
			isNotMatched = false
			api.RegisterOrchestratorServer(s.gRPC, orchestratorSv)
			if err := api.RegisterOrchestratorHandlerFromEndpoint(
				context.Background(),
				s.mux,
				s.cfg.GRPC.String(),
				[]grpc.DialOption{grpc.WithInsecure()},
			); err != nil {
				return err
			}
		}

		if isNotMatched {
			return fmt.Errorf("unknown GRPC Service to register %#v", srv)
		}
	}
	return nil
}

func isRunningInDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

// Serve server listen for HTTP and GRPC
func (s *Server) Serve() error {
	stop := make(chan os.Signal, 1)
	errch := make(chan error)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.Handler())
	httpMux.Handle("/", s.mux)
	httpServer := http.Server{
		Addr:    s.cfg.HTTP.String(),
		Handler: httpMux,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errch <- err
		}
	}()
	go func() {
		listener, err := net.Listen("tcp", s.cfg.GRPC.String())
		if err != nil {
			errch <- err
			return
		}
		if err := s.gRPC.Serve(listener); err != nil {
			errch <- err
		}
	}()
	for {
		select {
		case <-stop:
			ctx := context.Background()
			httpServer.Shutdown(ctx)
			s.gRPC.GracefulStop()
			if !isRunningInDockerContainer() {
				fmt.Println("Shutting down. Wait for 5 seconds")
				time.Sleep(5 * time.Second)
			}
			return nil
		case err := <-errch:
			return err
		}
	}
}
