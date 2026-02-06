package com.foorbar.financialservices.domain.model;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Entity
@Table(name = "currencies")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Currency {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(unique = true, nullable = false, length = 3)
    private String code;
    
    @Column(nullable = false)
    private String name;
    
    @Column
    private String symbol;
    
    @Column
    private Integer decimalPlaces = 2;
    
    @Column
    private Boolean active = true;
}
