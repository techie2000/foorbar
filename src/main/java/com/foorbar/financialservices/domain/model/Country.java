package com.foorbar.financialservices.domain.model;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Entity
@Table(name = "countries")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Country {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(unique = true, nullable = false, length = 2)
    private String code;
    
    @Column(nullable = false)
    private String name;
    
    @Column(length = 3)
    private String alpha3Code;
    
    @Column
    private String region;
    
    @Column
    private Boolean active = true;
}
