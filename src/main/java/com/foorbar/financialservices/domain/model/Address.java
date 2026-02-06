package com.foorbar.financialservices.domain.model;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Entity
@Table(name = "addresses")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Address {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column
    private String street;
    
    @Column
    private String city;
    
    @Column
    private String state;
    
    @Column
    private String postalCode;
    
    @ManyToOne
    @JoinColumn(name = "country_id")
    private Country country;
}
