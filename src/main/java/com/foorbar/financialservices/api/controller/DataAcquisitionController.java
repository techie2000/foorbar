package com.foorbar.financialservices.api.controller;

import com.foorbar.financialservices.dataacquisition.service.DataAcquisitionService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/api/data")
@RequiredArgsConstructor
@Slf4j
public class DataAcquisitionController {
    
    private final DataAcquisitionService dataAcquisitionService;
    
    @PostMapping("/acquire")
    public ResponseEntity<String> acquireData(@RequestBody Map<String, String> request) {
        try {
            String source = request.get("source");
            String type = request.get("type");
            String targetClass = request.get("targetClass");
            
            log.info("Acquiring data from source: {}, type: {}", source, type);
            // This is a simplified example - in production, you'd need proper class resolution
            
            return ResponseEntity.ok("Data acquisition initiated");
        } catch (Exception e) {
            log.error("Error acquiring data", e);
            return ResponseEntity.badRequest().body("Error: " + e.getMessage());
        }
    }
    
    @PostMapping("/distribute")
    public ResponseEntity<String> distributeData(@RequestBody Map<String, Object> request) {
        try {
            String destination = (String) request.get("destination");
            String type = (String) request.get("type");
            
            log.info("Distributing data to destination: {}, type: {}", destination, type);
            
            return ResponseEntity.ok("Data distribution initiated");
        } catch (Exception e) {
            log.error("Error distributing data", e);
            return ResponseEntity.badRequest().body("Error: " + e.getMessage());
        }
    }
}
