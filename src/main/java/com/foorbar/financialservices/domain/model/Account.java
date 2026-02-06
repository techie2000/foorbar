package com.foorbar.financialservices.domain.model;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import java.math.BigDecimal;
import java.time.LocalDateTime;

@Entity
@Table(name = "accounts")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Account {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(unique = true, nullable = false)
    private String accountNumber;
    
    @ManyToOne
    @JoinColumn(name = "entity_id")
    private EntityModel entity;
    
    @ManyToOne
    @JoinColumn(name = "currency_id")
    private Currency currency;
    
    @Column
    @Enumerated(EnumType.STRING)
    private AccountType type;
    
    @Column(precision = 19, scale = 4)
    private BigDecimal balance = BigDecimal.ZERO;
    
    @Column
    private LocalDateTime openedAt = LocalDateTime.now();
    
    @Column
    private Boolean active = true;
    
    public enum AccountType {
        TRADING, SETTLEMENT, CUSTODY, MARGIN
    }
}
