package com.foorbar.financialservices.domain.model;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Entity
@Table(name = "instruments")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Instrument {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(unique = true, nullable = false)
    private String isin;
    
    @Column(nullable = false)
    private String name;
    
    @Column
    @Enumerated(EnumType.STRING)
    private InstrumentType type;
    
    @ManyToOne
    @JoinColumn(name = "currency_id")
    private Currency currency;
    
    @Column
    private String exchange;
    
    @Column
    private Boolean active = true;
    
    public enum InstrumentType {
        EQUITY, BOND, DERIVATIVE, COMMODITY, FUND, FOREX
    }
}
