/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"dubbo.apache.org/dubbo-go/v3/config"
	
	_ "dubbo.apache.org/dubbo-go/v3/registry/polaris"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"
)

type UserProviderWithCustomGroupAndVersion struct {
	GetUser func(ctx context.Context, req *User) (rsp *User, err error)
}

type UserProvider struct {
	GetUser func(ctx context.Context, req *User) (rsp *User, err error)
}

type User struct {
	ID   string    `json:"id" xml:"id"`
	Name string    `json:"name" xml:"name"`
	Age  int32     `json:"age" xml:"age"`
	Time time.Time `json:"time" xml:"time"`
}

func (u *User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

type PolarisDubboGoConsumer struct {
	userProvider                                  *UserProvider
	userProviderWithCustomRegistryGroupAndVersion *UserProviderWithCustomGroupAndVersion
}

func (svr *PolarisDubboGoConsumer) Run() {
	svr.runWebServer()
}

func (svr *PolarisDubboGoConsumer) runWebServer() {
	http.HandleFunc("/echo", func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("\n\n\nstart to test dubbo")
		resp := make(map[string]interface{}, 2)
		user, err := svr.userProvider.GetUser(context.TODO(), &User{Name: "PolarisMesh"})
		if err != nil {
			log.Printf("error: %v\n", err)
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte("Hello, Occur error : " + err.Error()))
			return
		}
		resp["UserProvider"] = user

		user, err = svr.userProviderWithCustomRegistryGroupAndVersion.GetUser(context.TODO(), &User{Name: "PolarisMesh"})
		if err != nil {
			log.Printf("error: %v\n", err)
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte("Hello, Occur error : " + err.Error()))
			return
		}

		resp["UserProviderWithCustomGroupAndVersion"] = user
		data, _ := json.Marshal(resp)
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(data)
	})

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 18080), nil); err != nil {
		log.Fatalf("[ERROR]fail to run webServer, err is %v", err)
	}
}

func main() {

	consumer := &PolarisDubboGoConsumer{
		userProvider: &UserProvider{},
		userProviderWithCustomRegistryGroupAndVersion: &UserProviderWithCustomGroupAndVersion{},
	}

	config.SetConsumerService(consumer.userProvider)
	config.SetConsumerService(consumer.userProviderWithCustomRegistryGroupAndVersion)
	hessian.RegisterPOJO(&User{})
	err := config.Load()
	if err != nil {
		panic(err)
	}

	consumer.Run()
}
