package com.foorbar.financialservices.api.controller;

import com.foorbar.financialservices.domain.model.Country;
import com.foorbar.financialservices.domain.service.CountryService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/domain/countries")
@RequiredArgsConstructor
public class CountryController {
    
    private final CountryService countryService;
    
    @GetMapping
    public List<Country> getAllCountries() {
        return countryService.findAll();
    }
    
    @GetMapping("/{id}")
    public ResponseEntity<Country> getCountryById(@PathVariable Long id) {
        return countryService.findById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @GetMapping("/code/{code}")
    public ResponseEntity<Country> getCountryByCode(@PathVariable String code) {
        return countryService.findByCode(code)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @PostMapping
    public Country createCountry(@RequestBody Country country) {
        return countryService.save(country);
    }
    
    @PutMapping("/{id}")
    public ResponseEntity<Country> updateCountry(@PathVariable Long id, @RequestBody Country country) {
        return countryService.findById(id)
                .map(existing -> {
                    country.setId(id);
                    return ResponseEntity.ok(countryService.save(country));
                })
                .orElse(ResponseEntity.notFound().build());
    }
    
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteCountry(@PathVariable Long id) {
        countryService.deleteById(id);
        return ResponseEntity.noContent().build();
    }
}
