
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
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"

	_ "dubbo.apache.org/dubbo-go/v3/imports"

	_ "dubbo.apache.org/dubbo-go/v3/registry/polaris"

	hessian "github.com/apache/dubbo-go-hessian2"
)

var (
	survivalTimeout = int(3e9)
)

func init() {
	config.SetProviderService(&UserProvider{})
	config.SetProviderService(&UserProviderWithCustomGroupAndVersion{})
	// ------for hessian2------
	hessian.RegisterPOJO(&User{})
}

type User struct {
	ID   string
	Name string
	Age  int32
	Time time.Time
}

type UserProvider struct {
}

func (u *UserProvider) GetUser(ctx context.Context, req *User) (*User, error) {
	logger.Infof("req:%#v", req)
	rsp := User{"A001", "Alex Stocks", 18, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

type UserProviderWithCustomGroupAndVersion struct {
}

func (u *UserProviderWithCustomGroupAndVersion) GetUser(ctx context.Context, req *User) (*User, error) {
	logger.Infof("req:%#v", req)
	rsp := User{"A001", "Alex Stocks from UserProviderWithCustomGroupAndVersion", 18, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

// need to setup environment variable "CONF_PROVIDER_FILE_PATH" to "conf/server.yml" before run
func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			time.AfterFunc(time.Duration(survivalTimeout), func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})

			// The program exits normally or timeout forcibly exits.
			fmt.Println("provider app exit now...")
			return
		}
	}
}
