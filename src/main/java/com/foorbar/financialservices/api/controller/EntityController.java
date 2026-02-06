package com.foorbar.financialservices.api.controller;

import com.foorbar.financialservices.domain.model.EntityModel;
import com.foorbar.financialservices.domain.service.EntityService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/domain/entities")
@RequiredArgsConstructor
public class EntityController {
    
    private final EntityService entityService;
    
    @GetMapping
    public List<EntityModel> getAllEntities() {
        return entityService.findAll();
    }
    
    @GetMapping("/{id}")
    public ResponseEntity<EntityModel> getEntityById(@PathVariable Long id) {
        return entityService.findById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @GetMapping("/registration/{registrationNumber}")
    public ResponseEntity<EntityModel> getEntityByRegistrationNumber(@PathVariable String registrationNumber) {
        return entityService.findByRegistrationNumber(registrationNumber)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @PostMapping
    public EntityModel createEntity(@RequestBody EntityModel entity) {
        return entityService.save(entity);
    }
    
    @PutMapping("/{id}")
    public ResponseEntity<EntityModel> updateEntity(@PathVariable Long id, @RequestBody EntityModel entity) {
        return entityService.findById(id)
                .map(existing -> {
                    entity.setId(id);
                    return ResponseEntity.ok(entityService.save(entity));
                })
                .orElse(ResponseEntity.notFound().build());
    }
    
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteEntity(@PathVariable Long id) {
        entityService.deleteById(id);
        return ResponseEntity.noContent().build();
    }
}
