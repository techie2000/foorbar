package com.foorbar.financialservices.domain.service;

import com.foorbar.financialservices.domain.model.SSI;
import com.foorbar.financialservices.domain.repository.SSIRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
@Transactional
public class SSIService {
    
    private final SSIRepository ssiRepository;
    
    public SSI save(SSI ssi) {
        return ssiRepository.save(ssi);
    }
    
    public Optional<SSI> findById(Long id) {
        return ssiRepository.findById(id);
    }
    
    public List<SSI> findAll() {
        return ssiRepository.findAll();
    }
    
    public void deleteById(Long id) {
        ssiRepository.deleteById(id);
    }
}
