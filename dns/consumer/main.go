/**
 * Tencent is pleased to support the open source community by making polaris-go available.
 *
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
	namespace    string
	service      string
	providerPort int64
	port         int64
)

func initArgs() {
	flag.Int64Var(&port, "port", 20000, "self port")
	flag.StringVar(&namespace, "calleeNamespace", "default", "callee namespace")
	flag.StringVar(&service, "calleeService", "echoserver", "callee service")
	flag.Int64Var(&providerPort, "calleePort", 10000, "callee port")
}

// PolarisConsumer is a consumer of polaris
type PolarisConsumer struct {
	namespace string
	service   string
}

// Run starts the consumer
func (svr *PolarisConsumer) Run() {
	svr.runWebServer()
}

func (svr *PolarisConsumer) runWebServer() {
	http.HandleFunc("/echo", func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("start to invoke by dns operation")
		url := fmt.Sprintf("http://%s.%s:%d/echo", service, namespace, providerPort)
		gReq, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte(fmt.Sprintf("[errot] send request to %s fail : %s", url, err)))
			return
		}
		for k, v := range r.Header {
			for i := range v {
				r.Header.Add(k, v[i])
			}
		}
		resp, err := http.DefaultClient.Do(gReq)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte(fmt.Sprintf("[errot] send request to %s fail : %s", url, err)))
			return
		}

		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[error] read resp from %s fail : %s", url, err)
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte(fmt.Sprintf("[error] read resp from %s fail : %s", url, err)))
			return
		}
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(data)
	})

	http.HandleFunc("/echoWithParams", func(rw http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		ns := queryParams.Get("ns")
		svc := queryParams.Get("svc")
		portStr := queryParams.Get("port")
		if ns == "" || svc == "" || portStr == "" {
			rw.WriteHeader(http.StatusBadRequest)
			_, _ = rw.Write([]byte("Missing required query parameters: ns, svc, port"))
			return
		}
		port, err := strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			_, _ = rw.Write([]byte("Invalid port parameter"))
			return
		}
		url := fmt.Sprintf("http://%s.%s:%d/echo", svc, ns, port)
		gReq, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte(fmt.Sprintf("[errot] send request to %s fail : %s", url, err)))
			return
		}
		for k, v := range r.Header {
			for i := range v {
				r.Header.Add(k, v[i])
			}
		}
		resp, err := http.DefaultClient.Do(gReq)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte(fmt.Sprintf("[errot] send request to %s fail : %s", url, err)))
			return
		}

		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[error] read resp from %s fail : %s", url, err)
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte(fmt.Sprintf("[error] read resp from %s fail : %s", url, err)))
			return
		}
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(data)
	})

	log.Printf("start run web server, port : %d", port)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil); err != nil {
		log.Fatalf("[ERROR]fail to run webServer, err is %v", err)
	}
}

func main() {
	initArgs()
	flag.Parse()
	if len(namespace) == 0 || len(service) == 0 {
		log.Print("namespace and service are required")
		return
	}

	svr := &PolarisConsumer{
		namespace: namespace,
		service:   service,
	}

	svr.Run()

}
