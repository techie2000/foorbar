package com.foorbar.financialservices.domain.model;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import java.time.LocalDateTime;

@Entity
@Table(name = "entities")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class EntityModel {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(nullable = false)
    private String name;
    
    @Column(unique = true)
    private String registrationNumber;
    
    @Column
    @Enumerated(EnumType.STRING)
    private EntityType type;
    
    @ManyToOne(cascade = CascadeType.ALL)
    @JoinColumn(name = "address_id")
    private Address address;
    
    @Column
    private LocalDateTime createdAt = LocalDateTime.now();
    
    @Column
    private Boolean active = true;
    
    public enum EntityType {
        COMPANY, BUSINESS, CORPORATION, PARTNERSHIP, INDIVIDUAL
    }
}
