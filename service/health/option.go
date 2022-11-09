package health

import (
	"gitlab.com/dipper-iot/shared/cli"
	"gitlab.com/dipper-iot/shared/service"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func HealthService() service.Option {
	srv := &Health{}
	return func(o *service.Options) {
		o.BeforeStart(func(o *service.Options, c *cli.Context) error {
			healthpb.RegisterHealthServer(o.Server, srv)
			return nil
		})
	}
}

func HealthServiceOption(srv healthpb.HealthServer) service.Option {
	if srv == nil {
		srv = &Health{}
	}
	return func(o *service.Options) {
		o.BeforeStart(func(o *service.Options, c *cli.Context) error {
			healthpb.RegisterHealthServer(o.Server, srv)
			return nil
		})
	}
}
