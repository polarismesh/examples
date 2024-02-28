package com.tencent.polaris.dubbo.discovery.example.consumer;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import org.apache.dubbo.config.spring.context.annotation.EnableDubbo;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;

import java.io.IOException;
import java.io.OutputStream;
import java.io.UnsupportedEncodingException;
import java.net.InetSocketAddress;
import java.net.URI;
import java.net.URLDecoder;
import java.util.LinkedHashMap;
import java.util.Map;

public class Main {
    private static final int LISTEN_PORT = 15700;

    private static final String PATH = "/echo";
    public static void main(String[] args) throws Exception {
        AnnotationConfigApplicationContext context = new AnnotationConfigApplicationContext(ConsumerConfiguration.class);
        context.start();
        HelloConsumer greetingServiceConsumer = context.getBean(HelloConsumer.class);

        HttpServer server = HttpServer.create(new InetSocketAddress(LISTEN_PORT), 0);
        server.createContext(PATH, new EchoClientHandler(greetingServiceConsumer));
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            server.stop(1);
        }));
        server.start();
    }

    private static class EchoClientHandler implements HttpHandler {

        private final HelloConsumer consumer;

        public EchoClientHandler(HelloConsumer consumer) {
            this.consumer = consumer;
        }

        @Override
        public void handle(HttpExchange exchange) throws IOException {
            Map<String, String> parameters = splitQuery(exchange.getRequestURI());
            String echoValue = parameters.get("value");
            String response = consumer.doSay(echoValue);
            exchange.sendResponseHeaders(200, 0);
            OutputStream os = exchange.getResponseBody();
            os.write(response.getBytes());
            os.close();
        }

        private static Map<String, String> splitQuery(URI uri) throws UnsupportedEncodingException {
            Map<String, String> query_pairs = new LinkedHashMap<>();
            String query = uri.getQuery();
            String[] pairs = query.split("&");
            for (String pair : pairs) {
                int idx = pair.indexOf("=");
                query_pairs.put(URLDecoder.decode(pair.substring(0, idx), "UTF-8"),
                        URLDecoder.decode(pair.substring(idx + 1), "UTF-8"));
            }
            return query_pairs;
        }

    }


    @Configuration
    @EnableDubbo(scanBasePackages = "com.tencent.polaris.dubbo.discovery.example.consumer")
    @PropertySource("classpath:/spring/dubbo-consumer.properties")
    @ComponentScan(value = {"com.tencent.polaris.dubbo.discovery.example.consumer"})
    static class ConsumerConfiguration {

    }
}