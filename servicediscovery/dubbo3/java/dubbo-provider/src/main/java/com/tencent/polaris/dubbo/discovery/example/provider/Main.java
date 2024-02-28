package com.tencent.polaris.dubbo.discovery.example.provider;

import org.apache.dubbo.config.spring.context.annotation.EnableDubbo;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;

import java.util.concurrent.CountDownLatch;

public class Main {
    public static void main(String[] args) throws Exception {
        AnnotationConfigApplicationContext context = new AnnotationConfigApplicationContext(ProviderConfiguration.class);
        context.start();
        System.out.println("dubbo service started");
        new CountDownLatch(1).await();
    }

    @Configuration
    @EnableDubbo(scanBasePackages = "com.tencent.polaris.dubbo.discovery.example.provider")
    @PropertySource("classpath:/spring/dubbo-provider.properties")
    @ComponentScan(value = {"com.tencent.polaris.dubbo.discovery.example.provider"})
    static class ProviderConfiguration {

    }
}