package com.foorbar.financialservices.gateway;

import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

import jakarta.annotation.PostConstruct;
import java.util.HashMap;
import java.util.Map;

@Component
@Slf4j
public class ServiceRegistry {
    
    private final Map<String, ServiceInfo> services = new HashMap<>();
    
    @PostConstruct
    public void init() {
        // Register internal services
        registerService("domain-data-service", "http://localhost:8080", "Domain Data Service");
        registerService("data-acquisition-service", "http://localhost:8080", "Data Acquisition Service");
        registerService("notification-service", "http://localhost:8080", "Notification Service");
        
        log.info("Service registry initialized with {} services", services.size());
    }
    
    public void registerService(String serviceId, String url, String description) {
        ServiceInfo serviceInfo = new ServiceInfo(serviceId, url, description);
        services.put(serviceId, serviceInfo);
        log.info("Registered service: {} at {}", serviceId, url);
    }
    
    public ServiceInfo getService(String serviceId) {
        return services.get(serviceId);
    }
    
    public Map<String, ServiceInfo> getAllServices() {
        return new HashMap<>(services);
    }
    
    public void deregisterService(String serviceId) {
        services.remove(serviceId);
        log.info("Deregistered service: {}", serviceId);
    }
    
    public static class ServiceInfo {
        private final String serviceId;
        private final String url;
        private final String description;
        
        public ServiceInfo(String serviceId, String url, String description) {
            this.serviceId = serviceId;
            this.url = url;
            this.description = description;
        }
        
        public String getServiceId() { return serviceId; }
        public String getUrl() { return url; }
        public String getDescription() { return description; }
    }
}
