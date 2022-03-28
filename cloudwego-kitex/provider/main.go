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
	"log"
	"net"
	"strings"

	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api"
	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api/hello"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	polaris "github.com/kitex-contrib/registry-polaris"
	"github.com/polarismesh/polaris-go/pkg/config"
)

const (
	confPath  = "./polaris.yaml"
	Namespace = "default"
	// At present,polaris server tag is v1.4.0，can't support auto create namespace,
	// If you want to use a namespace other than default,Polaris ,before you register an instance,
	// you should create the namespace at polaris console first.
)

type HelloImpl struct{}

func (h *HelloImpl) Echo(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	resp = &api.Response{
		Message: req.Message + "Hi,Kitex!",
	}
	return
}

//  // https://www.cloudwego.io/docs/kitex/tutorials/framework-exten/registry/#integrate-into-kitex
func main() {
	Conf, err := config.LoadConfigurationByFile(confPath)
	if err != nil {
		log.Fatal(err)
	}
	polarisAddresses := Conf.Global.ServerConnector.Addresses

	r, err := polaris.NewPolarisRegistry(polarisAddresses)
	if err != nil {
		log.Fatal(err)
	}

	addr, err := getLocalHost(polarisAddresses[0])
	if err != nil {
		log.Fatal(err)
	}

	Info := &registry.Info{
		ServiceName: "EchoServerKitex",
		Tags: map[string]string{
			"namespace": Namespace,
		},
	}
	newServer := hello.NewServer(new(HelloImpl), server.WithRegistry(r), server.WithRegistryInfo(Info),
		server.WithServiceAddr(&net.TCPAddr{IP: net.ParseIP(addr), Port: 8888}))

	err = newServer.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func getLocalHost(serverAddr string) (string, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if nil != err {
		return "", err
	}
	localAddr := conn.LocalAddr().String()
	colonIdx := strings.LastIndex(localAddr, ":")
	if colonIdx > 0 {
		return localAddr[:colonIdx], nil
	}
	return localAddr, nil
}
