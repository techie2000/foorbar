package com.foorbar.financialservices.dataacquisition.service;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.csv.CsvMapper;
import com.fasterxml.jackson.dataformat.csv.CsvSchema;
import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class FileDataService {
    
    private final ObjectMapper jsonMapper = new ObjectMapper();
    private final XmlMapper xmlMapper = new XmlMapper();
    private final CsvMapper csvMapper = new CsvMapper();
    
    public <T> List<T> readJsonFile(String filePath, Class<T> clazz) throws IOException {
        log.info("Reading JSON file: {}", filePath);
        return jsonMapper.readValue(
            new File(filePath),
            jsonMapper.getTypeFactory().constructCollectionType(List.class, clazz)
        );
    }
    
    public <T> void writeJsonFile(String filePath, List<T> data) throws IOException {
        log.info("Writing JSON file: {}", filePath);
        jsonMapper.writerWithDefaultPrettyPrinter().writeValue(new File(filePath), data);
    }
    
    public <T> List<T> readXmlFile(String filePath, Class<T> clazz) throws IOException {
        log.info("Reading XML file: {}", filePath);
        String xmlContent = new String(Files.readAllBytes(Paths.get(filePath)));
        return xmlMapper.readValue(
            xmlContent,
            xmlMapper.getTypeFactory().constructCollectionType(List.class, clazz)
        );
    }
    
    public <T> void writeXmlFile(String filePath, List<T> data) throws IOException {
        log.info("Writing XML file: {}", filePath);
        xmlMapper.writerWithDefaultPrettyPrinter().writeValue(new File(filePath), data);
    }
    
    public <T> List<T> readCsvFile(String filePath, Class<T> clazz) throws IOException {
        log.info("Reading CSV file: {}", filePath);
        CsvSchema schema = CsvSchema.emptySchema().withHeader();
        return csvMapper.readerFor(clazz)
            .with(schema)
            .<T>readValues(new File(filePath))
            .readAll();
    }
    
    public <T> void writeCsvFile(String filePath, List<T> data, Class<T> clazz) throws IOException {
        log.info("Writing CSV file: {}", filePath);
        CsvSchema schema = csvMapper.schemaFor(clazz).withHeader();
        csvMapper.writerFor(clazz)
            .with(schema)
            .writeValues(new File(filePath))
            .writeAll(data);
    }
}
