package com.foorbar.financialservices.domain.service;

import com.foorbar.financialservices.domain.model.Country;
import com.foorbar.financialservices.domain.repository.CountryRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
@Transactional
public class CountryService {
    
    private final CountryRepository countryRepository;
    
    public Country save(Country country) {
        return countryRepository.save(country);
    }
    
    public Optional<Country> findById(Long id) {
        return countryRepository.findById(id);
    }
    
    public Optional<Country> findByCode(String code) {
        return countryRepository.findByCode(code);
    }
    
    public List<Country> findAll() {
        return countryRepository.findAll();
    }
    
    public void deleteById(Long id) {
        countryRepository.deleteById(id);
    }
}
