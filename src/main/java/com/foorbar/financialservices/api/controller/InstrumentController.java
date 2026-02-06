package com.foorbar.financialservices.api.controller;

import com.foorbar.financialservices.domain.model.Instrument;
import com.foorbar.financialservices.domain.service.InstrumentService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/domain/instruments")
@RequiredArgsConstructor
public class InstrumentController {
    
    private final InstrumentService instrumentService;
    
    @GetMapping
    public List<Instrument> getAllInstruments() {
        return instrumentService.findAll();
    }
    
    @GetMapping("/{id}")
    public ResponseEntity<Instrument> getInstrumentById(@PathVariable Long id) {
        return instrumentService.findById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @GetMapping("/isin/{isin}")
    public ResponseEntity<Instrument> getInstrumentByIsin(@PathVariable String isin) {
        return instrumentService.findByIsin(isin)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @PostMapping
    public Instrument createInstrument(@RequestBody Instrument instrument) {
        return instrumentService.save(instrument);
    }
    
    @PutMapping("/{id}")
    public ResponseEntity<Instrument> updateInstrument(@PathVariable Long id, @RequestBody Instrument instrument) {
        return instrumentService.findById(id)
                .map(existing -> {
                    instrument.setId(id);
                    return ResponseEntity.ok(instrumentService.save(instrument));
                })
                .orElse(ResponseEntity.notFound().build());
    }
    
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteInstrument(@PathVariable Long id) {
        instrumentService.deleteById(id);
        return ResponseEntity.noContent().build();
    }
}
