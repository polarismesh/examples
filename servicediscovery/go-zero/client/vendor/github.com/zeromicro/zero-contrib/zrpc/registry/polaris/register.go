package polaris

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/polarismesh/polaris-go/api"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/proc"

	"github.com/tal-tech/go-zero/core/netx"
)

var (
	provider api.ProviderAPI
	beatlock *sync.RWMutex                 = &sync.RWMutex{}
	hearts   map[string]context.CancelFunc = make(map[string]context.CancelFunc)
)

func init() {
	p, err := api.NewProviderAPI()
	if err != nil {
		panic(err)
	}

	provider = p
}

// RegitserService register service to polaris
func RegitserService(opts *Options) error {
	pubListenOn := figureOutListenOn(opts.ListenOn)

	host, ports, err := net.SplitHostPort(pubListenOn)
	if err != nil {
		return fmt.Errorf("failed parsing address error: %v", err)
	}
	port, _ := strconv.ParseInt(ports, 10, 64)

	if err != nil {
		log.Panic(err)
	}

	req := &api.InstanceRegisterRequest{}
	req.Service = opts.ServiceName
	req.Namespace = opts.Namespace
	req.Version = &opts.Version
	req.Protocol = &opts.Protocol
	req.Host = host
	req.Port = int(port)
	req.TTL = &opts.HeartbeatInervalSec

	resp, err := provider.Register(req)
	if err != nil {
		return err
	}
	if resp.Existed {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	beatlock.Lock()
	hearts[fmt.Sprintf("%s_%s_%s", opts.Namespace, opts.ServiceName, opts.ListenOn)] = cancel
	beatlock.Unlock()

	go doHeartbeat(ctx, req, opts)

	addShutdownListener(req, opts)
	return nil
}

func addShutdownListener(registerReq *api.InstanceRegisterRequest, opts *Options) {
	// service deregister
	proc.AddShutdownListener(func() {
		beatlock.Lock()
		cancel := hearts[fmt.Sprintf("%s_%s_%s", opts.Namespace, opts.ServiceName, opts.ListenOn)]
		beatlock.Unlock()
		cancel()

		req := &api.InstanceDeRegisterRequest{}
		req.Namespace = opts.Namespace
		req.Service = opts.ServiceName
		req.Host = registerReq.Host
		req.Port = registerReq.Port
		provider.Deregister(req)
	})
}

// doHeartbeat
func doHeartbeat(ctx context.Context, req *api.InstanceRegisterRequest, opts *Options) {
	ticker := time.NewTicker(time.Duration(opts.HeartbeatInervalSec * int(time.Second)))
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			beatreq := &api.InstanceHeartbeatRequest{}
			beatreq.Namespace = opts.Namespace
			beatreq.Service = opts.ServiceName
			beatreq.Host = req.Host
			beatreq.Port = req.Port

			if err := provider.Heartbeat(beatreq); err != nil {
				logx.Error("[Polaris provider] do heartbeat fail : %s, req : %#v", err.Error(), beatreq)
			}
		}
	}
}

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIP)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}
