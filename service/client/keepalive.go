package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

func Keepalive() grpc.DialOption {
	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 20 seconds if there is no activity
		Timeout:             5 * time.Second,  // wait 5 second for ping back
		PermitWithoutStream: true,             // send pings even without active streams
	}
	return grpc.WithKeepaliveParams(kacp)
}
