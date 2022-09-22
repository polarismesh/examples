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


package cn.polarismesh.user.service;

import cn.polarismesh.common.core.Consts;
import cn.polarismesh.common.core.HttpResult;
import cn.polarismesh.common.core.InitOptions;
import cn.polarismesh.common.core.Utils;
import com.sun.net.httpserver.Headers;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import java.io.IOException;
import java.io.OutputStream;

public class UserServerHandler implements HttpHandler {

    private final InitOptions initOptions;

    private final String nextAddress;

    public UserServerHandler(InitOptions initOptions) {
        this.initOptions = initOptions;
        nextAddress = Consts.DOMAIN_CREDIT;
        //nextAddress = "127.0.0.1";
    }

    public void handle(HttpExchange httpExchange) throws IOException {
        Headers requestHeaders = httpExchange.getRequestHeaders();
        System.out.println("http request received: headers " + requestHeaders.entrySet() + ", method " + httpExchange.getRequestMethod());
        String content = String.format("%s[%s] -> ", Consts.DOMAIN_USER, initOptions.getVersion());
        String urlStr = String.format("http://%s:%d%s", nextAddress, Consts.PORT_CREDIT, Consts.PATH);
        HttpResult httpResult = Utils.httpGet(urlStr, requestHeaders);
        System.out.println("http result: " + httpResult);
        content += httpResult.getMessage();
        System.out.println("start to send response " + content);
        httpExchange.sendResponseHeaders(200, 0);
        try (OutputStream output = httpExchange.getResponseBody()) {
            output.write(content.getBytes());
        }
    }
}
