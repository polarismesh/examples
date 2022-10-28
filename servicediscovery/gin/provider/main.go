package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polarismesh/polaris-go"
)

var (
	namespace string
	service   string
	token     string
	port      int64
)

func initArgs() {
	flag.StringVar(&namespace, "namespace", "default", "namespace")
	flag.StringVar(&service, "service", "DiscoverEchoServer", "service")
	// 当北极星开启鉴权时，需要配置此参数完成相关的权限检查
	flag.StringVar(&token, "token", "", "token")
	flag.Int64Var(&port, "port", 0, "port")
}

// PolarisProvider is an example of provider
type PolarisProvider struct {
	provider   polaris.ProviderAPI
	namespace  string
	service    string
	token      string
	host       string
	port       int
	isShutdown bool
}

// Run starts the provider
func (svr *PolarisProvider) Run() {
	tmpHost, err := getLocalHost(svr.provider.SDKContext().GetConfig().GetGlobal().GetServerConnector().GetAddresses()[0])
	if err != nil {
		panic(fmt.Errorf("error occur while fetching localhost: %v", err))
	}

	svr.host = tmpHost
	svr.runWebServer()
	svr.registerService()
	svr.doHeartbeat()
}

func (svr *PolarisProvider) runWebServer() {
	engine := gin.Default()
	engine.GET("/echo", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		_, _ = c.Writer.Write([]byte(fmt.Sprintf("Hello, I'm DiscoverEchoServer Provider, My host : %s:%d", svr.host, svr.port)))
	})

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("[ERROR]fail to listen tcp, err is %v", err)
	}

	svr.port = ln.Addr().(*net.TCPAddr).Port

	go func() {
		log.Printf("[INFO] start gin http server, listen port is %v", svr.port)
		if err := engine.RunListener(ln); err != nil {
			svr.isShutdown = false
			log.Fatalf("[ERROR]fail to run webServer, err is %v", err)
		}
	}()
}

func (svr *PolarisProvider) registerService() {
	registerRequest := &polaris.InstanceRegisterRequest{}
	registerRequest.Service = svr.service
	registerRequest.Namespace = svr.namespace
	registerRequest.Host = svr.host
	registerRequest.Port = svr.port
	registerRequest.ServiceToken = svr.token
	registerRequest.SetTTL(10)
	resp, err := svr.provider.Register(registerRequest)
	if err != nil {
		log.Fatalf("fail to register instance, err is %v", err)
	}
	log.Printf("register response: instanceId %s", resp.InstanceID)
}

func (svr *PolarisProvider) deregisterService() {
	log.Printf("start to invoke deregister operation")
	deregisterRequest := &polaris.InstanceDeRegisterRequest{}
	deregisterRequest.Service = svr.service
	deregisterRequest.Namespace = svr.namespace
	deregisterRequest.Host = svr.host
	deregisterRequest.Port = svr.port
	deregisterRequest.ServiceToken = svr.token
	if err := svr.provider.Deregister(deregisterRequest); err != nil {
		log.Fatalf("fail to deregister instance, err is %v", err)
	}
	log.Printf("deregister successfully.")
}

func (svr *PolarisProvider) doHeartbeat() {
	log.Printf("start to invoke heartbeat operation")
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		if !svr.isShutdown {
			heartbeatRequest := &polaris.InstanceHeartbeatRequest{}
			heartbeatRequest.Namespace = svr.namespace
			heartbeatRequest.Service = svr.service
			heartbeatRequest.Host = svr.host
			heartbeatRequest.Port = svr.port
			heartbeatRequest.ServiceToken = svr.token
			svr.provider.Heartbeat(heartbeatRequest)
		}
	}
}

func (svr *PolarisProvider) runMainLoop() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, []os.Signal{
		syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGSEGV,
	}...)

	for s := range ch {
		log.Printf("catch signal(%+v), stop servers", s)
		svr.isShutdown = true
		svr.deregisterService()
		return
	}
}

func main() {
	initArgs()
	flag.Parse()
	if len(namespace) == 0 || len(service) == 0 {
		log.Print("namespace and service are required")
		return
	}
	provider, err := polaris.NewProviderAPI()
	// 或者使用以下方法,则不需要创建配置文件
	// provider, err = api.NewProviderAPIByAddress("127.0.0.1:8091")

	if err != nil {
		log.Fatalf("fail to create providerAPI, err is %v", err)
	}
	defer provider.Destroy()

	svr := PolarisProvider{
		provider:  provider,
		namespace: namespace,
		service:   service,
		token:     token,
	}

	svr.Run()

	svr.runMainLoop()
}

func getLocalHost(serverAddr string) (string, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return "", err
	}
	localAddr := conn.LocalAddr().String()
	colonIdx := strings.LastIndex(localAddr, ":")
	if colonIdx > 0 {
		return localAddr[:colonIdx], nil
	}
	return localAddr, nil
}
