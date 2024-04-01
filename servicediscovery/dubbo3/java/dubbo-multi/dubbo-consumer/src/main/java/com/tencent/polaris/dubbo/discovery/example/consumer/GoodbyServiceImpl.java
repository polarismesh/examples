package com.tencent.polaris.dubbo.discovery.example.consumer;

import com.tencent.polaris.dubbo.api.GoodbyService;
import org.apache.dubbo.config.annotation.DubboService;

@DubboService(version = "1.0.0", group = "")
public class GoodbyServiceImpl implements GoodbyService {
    @Override
    public String print(String name) {
        return "print " + name;
    }
}
