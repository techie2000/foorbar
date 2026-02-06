package com.foorbar.financialservices.domain.model;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import java.time.LocalDateTime;

@Entity
@Table(name = "ssis")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class SSI {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @ManyToOne
    @JoinColumn(name = "entity_id")
    private EntityModel entity;
    
    @ManyToOne
    @JoinColumn(name = "currency_id")
    private Currency currency;
    
    @ManyToOne
    @JoinColumn(name = "instrument_id")
    private Instrument instrument;
    
    @Column(nullable = false)
    private String beneficiaryName;
    
    @Column(nullable = false)
    private String beneficiaryAccount;
    
    @Column(nullable = false)
    private String beneficiaryBank;
    
    @Column
    private String beneficiaryBankBIC;
    
    @Column
    private String intermediaryBank;
    
    @Column
    private String intermediaryBankBIC;
    
    @Column
    @Enumerated(EnumType.STRING)
    private SettlementType settlementType;
    
    @Column
    private LocalDateTime validFrom = LocalDateTime.now();
    
    @Column
    private LocalDateTime validTo;
    
    @Column
    private Boolean active = true;
    
    public enum SettlementType {
        DVP, FOP, RVP, DAP
    }
}
