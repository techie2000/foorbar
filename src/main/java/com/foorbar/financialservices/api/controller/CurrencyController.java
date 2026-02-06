package com.foorbar.financialservices.api.controller;

import com.foorbar.financialservices.domain.model.Currency;
import com.foorbar.financialservices.domain.service.CurrencyService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/domain/currencies")
@RequiredArgsConstructor
public class CurrencyController {
    
    private final CurrencyService currencyService;
    
    @GetMapping
    public List<Currency> getAllCurrencies() {
        return currencyService.findAll();
    }
    
    @GetMapping("/{id}")
    public ResponseEntity<Currency> getCurrencyById(@PathVariable Long id) {
        return currencyService.findById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @GetMapping("/code/{code}")
    public ResponseEntity<Currency> getCurrencyByCode(@PathVariable String code) {
        return currencyService.findByCode(code)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @PostMapping
    public Currency createCurrency(@RequestBody Currency currency) {
        return currencyService.save(currency);
    }
    
    @PutMapping("/{id}")
    public ResponseEntity<Currency> updateCurrency(@PathVariable Long id, @RequestBody Currency currency) {
        return currencyService.findById(id)
                .map(existing -> {
                    currency.setId(id);
                    return ResponseEntity.ok(currencyService.save(currency));
                })
                .orElse(ResponseEntity.notFound().build());
    }
    
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteCurrency(@PathVariable Long id) {
        currencyService.deleteById(id);
        return ResponseEntity.noContent().build();
    }
}
