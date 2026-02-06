package com.foorbar.financialservices.domain.repository;

import com.foorbar.financialservices.domain.model.SSI;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface SSIRepository extends JpaRepository<SSI, Long> {
}
