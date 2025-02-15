package main

import (
	"GuGoTik/src/constant/config"
	"GuGoTik/src/extra/profiling"
	"GuGoTik/src/extra/tracing"
	"GuGoTik/src/rpc/health"
	"GuGoTik/src/rpc/publish"
	healthImpl "GuGoTik/src/services/health"
	"GuGoTik/src/utils/consul"
	"GuGoTik/src/utils/logging"
	"context"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"net"
)

func main() {
	tp, err := tracing.SetTraceProvider(config.PublishRpcServerName)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panicf("Error to set the trace")
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error to set the trace")
		}
	}()

	// Configure Pyroscope
	profiling.InitPyroscope("GuGoTik.PublishService")

	maxSize := 100 * 1024 * 1024
	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxSize),
		grpc.MaxSendMsgSize(maxSize),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)

	log := logging.LogService(config.PublishRpcServerName)
	lis, err := net.Listen("tcp", config.EnvCfg.PodIpAddr+config.PublishRpcServerPort)

	if err != nil {
		log.Panicf("Rpc %s listen happens error: %v", config.PublishRpcServerName, err)
	}

	var srv PublishServiceImpl
	var probe healthImpl.ProbeImpl
	defer CloseMQConn()
	publish.RegisterPublishServiceServer(s, srv)
	health.RegisterHealthServer(s, &probe)
	if err := consul.RegisterConsul(config.PublishRpcServerName, config.PublishRpcServerPort); err != nil {
		log.Panicf("Rpc %s register consul happens error for: %v", config.PublishRpcServerName, err)
	}
	srv.New()
	log.Infof("Rpc %s is running at %s now", config.PublishRpcServerName, config.PublishRpcServerPort)
	if err := s.Serve(lis); err != nil {
		log.Panicf("Rpc %s listen happens error for: %v", config.PublishRpcServerName, err)
	}
}
