package internalgrpc

import (
	grpcgenerated "github.com/oleglarionov/otusgolang_finalproject/internal/grpc/generated"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	port       string
	grpcServer *grpc.Server
	service    grpcgenerated.BannerRotationServiceServer
}

func NewServer(port string, service grpcgenerated.BannerRotationServiceServer) *Server {
	return &Server{
		port:    port,
		service: service,
	}
}

func (s *Server) Serve() error {
	addr := net.JoinHostPort("", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer()
	grpcgenerated.RegisterBannerRotationServiceServer(s.grpcServer, s.service)

	if err := s.grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
