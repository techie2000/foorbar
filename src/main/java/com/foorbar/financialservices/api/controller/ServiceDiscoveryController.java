package com.foorbar.financialservices.api.controller;

import com.foorbar.financialservices.gateway.ServiceRegistry;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/api/gateway")
@RequiredArgsConstructor
public class ServiceDiscoveryController {
    
    private final ServiceRegistry serviceRegistry;
    
    @GetMapping("/services")
    public ResponseEntity<Map<String, ServiceRegistry.ServiceInfo>> getAllServices() {
        return ResponseEntity.ok(serviceRegistry.getAllServices());
    }
    
    @GetMapping("/services/{serviceId}")
    public ResponseEntity<ServiceRegistry.ServiceInfo> getService(@PathVariable String serviceId) {
        ServiceRegistry.ServiceInfo service = serviceRegistry.getService(serviceId);
        if (service != null) {
            return ResponseEntity.ok(service);
        }
        return ResponseEntity.notFound().build();
    }
    
    @PostMapping("/services")
    public ResponseEntity<String> registerService(@RequestBody Map<String, String> request) {
        String serviceId = request.get("serviceId");
        String url = request.get("url");
        String description = request.get("description");
        
        serviceRegistry.registerService(serviceId, url, description);
        return ResponseEntity.ok("Service registered successfully");
    }
    
    @DeleteMapping("/services/{serviceId}")
    public ResponseEntity<String> deregisterService(@PathVariable String serviceId) {
        serviceRegistry.deregisterService(serviceId);
        return ResponseEntity.ok("Service deregistered successfully");
    }
}
