package com.foorbar.financialservices.dataacquisition.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class DataAcquisitionService {
    
    private final FileDataService fileDataService;
    private final QueueDataService queueDataService;
    
    public <T> List<T> acquireData(String source, String type, Class<T> clazz) throws IOException {
        log.info("Acquiring data from source: {}, type: {}", source, type);
        
        switch (type.toLowerCase()) {
            case "json":
                return fileDataService.readJsonFile(source, clazz);
            case "xml":
                return fileDataService.readXmlFile(source, clazz);
            case "csv":
                return fileDataService.readCsvFile(source, clazz);
            default:
                throw new IllegalArgumentException("Unsupported data type: " + type);
        }
    }
    
    public <T> void distributeData(String destination, String type, List<T> data, Class<T> clazz) throws IOException {
        log.info("Distributing data to destination: {}, type: {}", destination, type);
        
        if (destination.startsWith("queue:")) {
            String queueName = destination.substring(6);
            queueDataService.sendToQueue(queueName, data);
        } else {
            switch (type.toLowerCase()) {
                case "json":
                    fileDataService.writeJsonFile(destination, data);
                    break;
                case "xml":
                    fileDataService.writeXmlFile(destination, data);
                    break;
                case "csv":
                    fileDataService.writeCsvFile(destination, data, clazz);
                    break;
                default:
                    throw new IllegalArgumentException("Unsupported data type: " + type);
            }
        }
    }
}
