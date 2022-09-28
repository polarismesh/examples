/*
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


package cn.polarismesh.decorate.service;

import cn.polarismesh.common.core.Consts;
import cn.polarismesh.common.core.InitOptions;
import cn.polarismesh.common.core.Utils;
import com.sun.net.httpserver.HttpServer;
import java.io.IOException;
import java.net.InetSocketAddress;
import org.apache.commons.cli.ParseException;

public class DecorateService {

    public static void main(String[] args) {
        InitOptions initOptions;
        try {
            initOptions = Utils.initOptions(args);
        } catch (ParseException e) {
            throw new RuntimeException(e);
        }
        System.out.println("InitOption from arguments is " + initOptions);
        HttpServer server;
        try {
            server = HttpServer.create(new InetSocketAddress(Consts.PORT_DECORATE), 0);
        } catch (IOException e) {
            throw new RuntimeException("fail to listen on port " + Consts.PORT_DECORATE, e);
        }
        server.createContext(Consts.PATH, new DecorateHandler(initOptions));
        System.out.println("promotion service start successfully, listen on port " + Consts.PORT_DECORATE);
        server.start();
    }
}
