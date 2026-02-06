package com.foorbar.financialservices.domain.service;

import com.foorbar.financialservices.domain.model.Currency;
import com.foorbar.financialservices.domain.repository.CurrencyRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
@Transactional
public class CurrencyService {
    
    private final CurrencyRepository currencyRepository;
    
    public Currency save(Currency currency) {
        return currencyRepository.save(currency);
    }
    
    public Optional<Currency> findById(Long id) {
        return currencyRepository.findById(id);
    }
    
    public Optional<Currency> findByCode(String code) {
        return currencyRepository.findByCode(code);
    }
    
    public List<Currency> findAll() {
        return currencyRepository.findAll();
    }
    
    public void deleteById(Long id) {
        currencyRepository.deleteById(id);
    }
}
