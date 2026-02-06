package com.foorbar.financialservices.domain.repository;

import com.foorbar.financialservices.domain.model.Instrument;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.util.Optional;

@Repository
public interface InstrumentRepository extends JpaRepository<Instrument, Long> {
    Optional<Instrument> findByIsin(String isin);
}
