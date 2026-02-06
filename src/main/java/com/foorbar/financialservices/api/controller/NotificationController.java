package com.foorbar.financialservices.api.controller;

import com.foorbar.financialservices.notification.model.Notification;
import com.foorbar.financialservices.notification.service.NotificationService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/api/notifications")
@RequiredArgsConstructor
public class NotificationController {
    
    private final NotificationService notificationService;
    
    @PostMapping("/email")
    public ResponseEntity<String> sendEmailNotification(@RequestBody Map<String, String> request) {
        String recipient = request.get("recipient");
        String subject = request.get("subject");
        String message = request.get("message");
        
        notificationService.sendEmailNotification(recipient, subject, message);
        return ResponseEntity.ok("Email notification sent");
    }
    
    @PostMapping("/sms")
    public ResponseEntity<String> sendSmsNotification(@RequestBody Map<String, String> request) {
        String recipient = request.get("recipient");
        String message = request.get("message");
        
        notificationService.sendSmsNotification(recipient, message);
        return ResponseEntity.ok("SMS notification sent");
    }
    
    @PostMapping("/system")
    public ResponseEntity<String> sendSystemNotification(@RequestBody Map<String, String> request) {
        String recipient = request.get("recipient");
        String subject = request.get("subject");
        String message = request.get("message");
        
        notificationService.sendSystemNotification(recipient, subject, message);
        return ResponseEntity.ok("System notification sent");
    }
}
