package com.tencent.polaris.dubbo.discovery.example.provider;

import com.tencent.polaris.dubbo.api.HelloService;
import org.apache.dubbo.config.annotation.DubboService;

@DubboService(version = "1.0.0")
public class HelloServiceImpl implements HelloService {
    @Override
    public String say(String name) {
        return "hello " + name;
    }
}
