package com.example.decorate.service;

import com.tencent.cloud.polaris.router.spi.SpringWebRouterLabelResolver;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import org.springframework.http.HttpRequest;
import org.springframework.stereotype.Component;

@Component
public class GrayRouteLabelResolver implements SpringWebRouterLabelResolver {

    @Override
    public Map<String, String> resolve(HttpRequest request, byte[] body, Set<String> expressionLabelKeys) {
        List<String> gray = request.getHeaders().get("gray");
        Map<String, String> values = new HashMap<>();
        if (gray != null && !gray.isEmpty()) {
            values.put("${http.header.gray}", gray.get(0));
        }
        return values;
    }

    @Override
    public int getOrder() {
        return 0;
    }
}
