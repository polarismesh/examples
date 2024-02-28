package com.tencent.polaris.dubbo.discovery.example.consumer;

import com.tencent.polaris.dubbo.api.HelloService;
import org.apache.dubbo.config.annotation.DubboReference;
import org.springframework.stereotype.Component;

@Component
public class HelloConsumer {

    @DubboReference(version = "1.0.0")
    private HelloService helloService;

    public String doSay(String name) {
        return helloService.say(name);
    }
}
