package com.foorbar.financialservices.domain.service;

import com.foorbar.financialservices.domain.model.Account;
import com.foorbar.financialservices.domain.repository.AccountRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
@Transactional
public class AccountService {
    
    private final AccountRepository accountRepository;
    
    public Account save(Account account) {
        return accountRepository.save(account);
    }
    
    public Optional<Account> findById(Long id) {
        return accountRepository.findById(id);
    }
    
    public Optional<Account> findByAccountNumber(String accountNumber) {
        return accountRepository.findByAccountNumber(accountNumber);
    }
    
    public List<Account> findAll() {
        return accountRepository.findAll();
    }
    
    public void deleteById(Long id) {
        accountRepository.deleteById(id);
    }
}
