package com.foorbar.financialservices.domain.service;

import com.foorbar.financialservices.domain.model.Country;
import com.foorbar.financialservices.domain.repository.CountryRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Arrays;
import java.util.List;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CountryServiceTest {

    @Mock
    private CountryRepository countryRepository;

    @InjectMocks
    private CountryService countryService;

    private Country testCountry;

    @BeforeEach
    void setUp() {
        testCountry = new Country();
        testCountry.setId(1L);
        testCountry.setCode("US");
        testCountry.setName("United States");
        testCountry.setAlpha3Code("USA");
        testCountry.setRegion("North America");
        testCountry.setActive(true);
    }

    @Test
    void testSaveCountry() {
        when(countryRepository.save(any(Country.class))).thenReturn(testCountry);

        Country saved = countryService.save(testCountry);

        assertNotNull(saved);
        assertEquals("US", saved.getCode());
        verify(countryRepository, times(1)).save(testCountry);
    }

    @Test
    void testFindByCode() {
        when(countryRepository.findByCode("US")).thenReturn(Optional.of(testCountry));

        Optional<Country> found = countryService.findByCode("US");

        assertTrue(found.isPresent());
        assertEquals("United States", found.get().getName());
    }

    @Test
    void testFindAll() {
        Country country2 = new Country();
        country2.setId(2L);
        country2.setCode("GB");
        country2.setName("United Kingdom");

        when(countryRepository.findAll()).thenReturn(Arrays.asList(testCountry, country2));

        List<Country> countries = countryService.findAll();

        assertEquals(2, countries.size());
        verify(countryRepository, times(1)).findAll();
    }
}
