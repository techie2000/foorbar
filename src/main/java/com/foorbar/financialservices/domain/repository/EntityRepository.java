package com.foorbar.financialservices.domain.repository;

import com.foorbar.financialservices.domain.model.EntityModel;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.util.Optional;

@Repository
public interface EntityRepository extends JpaRepository<EntityModel, Long> {
    Optional<EntityModel> findByRegistrationNumber(String registrationNumber);
}
