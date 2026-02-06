package com.foorbar.financialservices.domain.service;

import com.foorbar.financialservices.domain.model.EntityModel;
import com.foorbar.financialservices.domain.repository.EntityRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
@Transactional
public class EntityService {
    
    private final EntityRepository entityRepository;
    
    public EntityModel save(EntityModel entity) {
        return entityRepository.save(entity);
    }
    
    public Optional<EntityModel> findById(Long id) {
        return entityRepository.findById(id);
    }
    
    public Optional<EntityModel> findByRegistrationNumber(String registrationNumber) {
        return entityRepository.findByRegistrationNumber(registrationNumber);
    }
    
    public List<EntityModel> findAll() {
        return entityRepository.findAll();
    }
    
    public void deleteById(Long id) {
        entityRepository.deleteById(id);
    }
}
