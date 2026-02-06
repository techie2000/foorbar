package com.foorbar.financialservices.dataacquisition.service;

import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
@Slf4j
public class QueueDataService {
    
    private final RabbitTemplate rabbitTemplate;
    private final ObjectMapper objectMapper = new ObjectMapper();
    
    public void sendToQueue(String queueName, Object data) {
        log.info("Sending data to queue: {}", queueName);
        rabbitTemplate.convertAndSend(queueName, data);
    }
    
    @RabbitListener(queues = "${app.queue.data-input:data-input-queue}")
    public void receiveFromQueue(String message) {
        log.info("Received message from queue: {}", message);
        // Process the message
    }
    
    public <T> void sendJsonToQueue(String queueName, T data) {
        try {
            String jsonData = objectMapper.writeValueAsString(data);
            sendToQueue(queueName, jsonData);
        } catch (Exception e) {
            log.error("Error sending JSON to queue", e);
        }
    }
}
