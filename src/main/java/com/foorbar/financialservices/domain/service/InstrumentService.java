package com.foorbar.financialservices.domain.service;

import com.foorbar.financialservices.domain.model.Instrument;
import com.foorbar.financialservices.domain.repository.InstrumentRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
@Transactional
public class InstrumentService {
    
    private final InstrumentRepository instrumentRepository;
    
    public Instrument save(Instrument instrument) {
        return instrumentRepository.save(instrument);
    }
    
    public Optional<Instrument> findById(Long id) {
        return instrumentRepository.findById(id);
    }
    
    public Optional<Instrument> findByIsin(String isin) {
        return instrumentRepository.findByIsin(isin);
    }
    
    public List<Instrument> findAll() {
        return instrumentRepository.findAll();
    }
    
    public void deleteById(Long id) {
        instrumentRepository.deleteById(id);
    }
}
