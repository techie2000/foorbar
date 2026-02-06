package com.foorbar.financialservices.gateway;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.function.RouterFunction;
import org.springframework.web.servlet.function.ServerResponse;

import static org.springframework.web.servlet.function.RouterFunctions.route;

@Configuration
public class GatewayConfig {
    
    @Bean
    public RouterFunction<ServerResponse> gatewayRoutes() {
        return route()
            .GET("/gateway/health", request -> 
                ServerResponse.ok().body("API Gateway is running"))
            .build();
    }
}
