package cn.polarismesh.examples.eureka.consumer;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.client.discovery.EnableDiscoveryClient;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

@SpringBootApplication
@RestController
@EnableDiscoveryClient
public class EurekaConsumerApplication {

    private static final String PROVIDER_SERVICE_NAME = "eureka-provider-service";

    @Autowired
    private RestTemplate restTemplate;

    @RequestMapping("/echo")
    public String echo(@RequestParam(required = false, defaultValue = PROVIDER_SERVICE_NAME) String providerServiceName, @RequestParam String value) {
        return restTemplate.getForObject(
                String.format("http://%s/echo1?value=%s", providerServiceName, value), String.class);
    }

    public static void main(String[] args) {
        SpringApplication.run(EurekaConsumerApplication.class, args);
    }
}
