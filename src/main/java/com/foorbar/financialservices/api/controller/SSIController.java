package com.foorbar.financialservices.api.controller;

import com.foorbar.financialservices.domain.model.SSI;
import com.foorbar.financialservices.domain.service.SSIService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/domain/ssis")
@RequiredArgsConstructor
public class SSIController {
    
    private final SSIService ssiService;
    
    @GetMapping
    public List<SSI> getAllSSIs() {
        return ssiService.findAll();
    }
    
    @GetMapping("/{id}")
    public ResponseEntity<SSI> getSSIById(@PathVariable Long id) {
        return ssiService.findById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @PostMapping
    public SSI createSSI(@RequestBody SSI ssi) {
        return ssiService.save(ssi);
    }
    
    @PutMapping("/{id}")
    public ResponseEntity<SSI> updateSSI(@PathVariable Long id, @RequestBody SSI ssi) {
        return ssiService.findById(id)
                .map(existing -> {
                    ssi.setId(id);
                    return ResponseEntity.ok(ssiService.save(ssi));
                })
                .orElse(ResponseEntity.notFound().build());
    }
    
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteSSI(@PathVariable Long id) {
        ssiService.deleteById(id);
        return ResponseEntity.noContent().build();
    }
}
