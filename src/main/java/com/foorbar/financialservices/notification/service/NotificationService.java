package com.foorbar.financialservices.notification.service;

import com.foorbar.financialservices.notification.model.Notification;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

@Service
@RequiredArgsConstructor
@Slf4j
public class NotificationService {
    
    private final NotificationRepository notificationRepository;
    
    public Notification createNotification(String recipient, String subject, String message, Notification.NotificationType type) {
        Notification notification = new Notification();
        notification.setRecipient(recipient);
        notification.setSubject(subject);
        notification.setMessage(message);
        notification.setType(type);
        notification.setStatus(Notification.NotificationStatus.PENDING);
        
        return notificationRepository.save(notification);
    }
    
    public void sendNotification(Long notificationId) {
        notificationRepository.findById(notificationId).ifPresent(notification -> {
            log.info("Sending notification {} to {} via {}", 
                notification.getId(), notification.getRecipient(), notification.getType());
            
            // Simulate sending notification
            notification.setStatus(Notification.NotificationStatus.SENT);
            notification.setSentAt(LocalDateTime.now());
            notificationRepository.save(notification);
            
            log.info("Notification {} sent successfully", notification.getId());
        });
    }
    
    public void sendEmailNotification(String recipient, String subject, String message) {
        Notification notification = createNotification(recipient, subject, message, Notification.NotificationType.EMAIL);
        sendNotification(notification.getId());
    }
    
    public void sendSmsNotification(String recipient, String message) {
        Notification notification = createNotification(recipient, "SMS", message, Notification.NotificationType.SMS);
        sendNotification(notification.getId());
    }
    
    public void sendSystemNotification(String recipient, String subject, String message) {
        Notification notification = createNotification(recipient, subject, message, Notification.NotificationType.SYSTEM);
        sendNotification(notification.getId());
    }
}
