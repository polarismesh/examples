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

package com.example.decorate.service;

import com.tencent.cloud.polaris.PolarisDiscoveryProperties;
import java.net.URI;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.client.loadbalancer.LoadBalanced;
import org.springframework.context.annotation.Bean;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.RequestEntity;
import org.springframework.http.ResponseEntity;
import org.springframework.util.MultiValueMap;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

@SpringBootApplication
public class DecorateServiceApplication {

    public static void main(String[] args) {
        SpringApplication.run(DecorateServiceApplication.class, args);
    }

    @LoadBalanced
    @Bean
    public RestTemplate restTemplate() {
        return new RestTemplate();
    }

    @RestController
    static class DecorateController {

        private final PolarisDiscoveryProperties properties;

        private final RestTemplate template;

        private final String gray;

        DecorateController(PolarisDiscoveryProperties properties, RestTemplate template) {
            this.properties = properties;
            this.template = template;
            this.gray = System.getenv("GRAY");
        }

        @GetMapping(value = "/echo")
        public String echo() {
            RequestEntity<String> entity;
            if (null != gray && gray.length() > 0) {
                MultiValueMap<String, String> headers = new HttpHeaders();
                headers.add("gray", "true");
                System.out.println("headers to send is is " + headers);
                entity = new RequestEntity<String>(headers, HttpMethod.GET, URI.create("http://user/echo"));
            } else {
                entity = new RequestEntity<String>(HttpMethod.GET,
                        URI.create("http://user/echo"));
            }
            ResponseEntity<String> exchange = template.exchange(entity, String.class);
            System.out.println("response received is " + exchange);
            return exchange.getBody();
        }

    }
}
