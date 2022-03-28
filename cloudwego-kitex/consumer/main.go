/**
 * Tencent is pleased to support the open source community by making Polaris available.
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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api"
	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api/hello"
	"github.com/cloudwego/kitex/client"
	polaris "github.com/kitex-contrib/registry-polaris"
	"github.com/polarismesh/polaris-go/pkg/config"
)

const (
	confPath  = "polaris.yaml"
	Namespace = "default"
	// At present,polaris server tag is v1.4.0，can't support auto create namespace,
	// if you want to use a namespace other than default,Polaris ,before you register an instance,
	// you should create the namespace at polaris console first.
)

type KitexConsumer struct {
	client hello.Client
}

func (svr *KitexConsumer) Run() {
	svr.runWebServer()
}

func (svr *KitexConsumer) runWebServer() {
	http.HandleFunc("/echo", func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("\n\n\nstart to test cloudwego/kitex")
		resp, err := svr.client.Echo(context.TODO(), &api.Request{Message: "Hi,polaris!"})
		if err != nil {
			log.Printf("error: %v\n", err)
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte("Hello, Occur error : " + err.Error()))
			return
		}
		data, _ := json.Marshal(resp)
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(data)
	})

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 18080), nil); err != nil {
		log.Fatalf("[ERROR]fail to run webServer, err is %v", err)
	}
}

func main() {
	Conf, err := config.LoadConfigurationByFile(confPath)
	if err != nil {
		log.Fatal(err)
	}
	polarisAddresses := Conf.Global.ServerConnector.Addresses

	r, err := polaris.NewPolarisResolver(polarisAddresses)
	if err != nil {
		log.Fatal(err)
	}

	// client.WithTag sets the namespace tag for service discovery
	newClient := hello.MustNewClient("EchoServerKitex", client.WithTag("namespace", Namespace),
		client.WithResolver(r), client.WithRPCTimeout(time.Second*60))

	consumer := &KitexConsumer{
		client: newClient,
	}

	consumer.Run()
}
