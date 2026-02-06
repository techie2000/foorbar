# Financial Services Static Data System

A modular monolith system for handling financial services static data with comprehensive support for data acquisition, distribution, and domain management.

## Overview

This system provides a robust platform for managing financial services static data including:

- **Domain Data Management**: Countries, Currencies, Entities, Instruments, Accounts, and SSI's
- **Data Acquisition/Distribution**: Support for CSV, XML, and JSON formats from files and queues
- **Notifications**: Multi-channel notification system (Email, SMS, System)
- **API Gateway**: Centralized API gateway with request routing
- **Service Discovery**: Internal service registry for service location and communication

## Architecture

The application follows a modular monolith architecture with clear service boundaries:

```
financial-services/
├── api/              # REST API Controllers
├── domain/           # Domain Models, Repositories, Services
├── dataacquisition/  # Data Import/Export Services
├── notification/     # Notification Services
└── gateway/          # API Gateway & Service Discovery
```

## Key Features

### 1. Domain Data Services

Manage core financial data entities:

- **Countries**: ISO country codes and related data
- **Currencies**: ISO currency codes with decimal places and symbols
- **Entities**: Companies, businesses, and legal entities with addresses
- **Instruments**: Financial instruments (equities, bonds, derivatives, etc.)
- **Accounts**: Trading, settlement, custody, and margin accounts
- **SSI's**: Standard Settlement Instructions for settlements

### 2. Data Acquisition/Distribution

Support for multiple data formats and sources:

- **Formats**: CSV, XML, JSON
- **Sources**: File systems, message queues (RabbitMQ)
- **Bidirectional**: Both import and export capabilities

### 3. Notification System

Multi-channel notification support:

- Email notifications
- SMS notifications  
- System notifications
- Notification tracking and status management

### 4. API Gateway

- Request routing to appropriate services
- Service discovery and registration
- Request filtering and header management

## Technology Stack

- **Framework**: Spring Boot 3.2.1
- **Language**: Java 17
- **Database**: H2 (in-memory, can be configured for production databases)
- **Message Queue**: RabbitMQ
- **API Gateway**: Spring Cloud Gateway MVC
- **Data Processing**: Jackson (JSON, XML, CSV)
- **ORM**: Spring Data JPA with Hibernate

## Getting Started

### Prerequisites

- Java 17 or higher
- Maven 3.6 or higher
- RabbitMQ (optional, for queue-based data acquisition)

### Building the Application

```bash
mvn clean install
```

### Running the Application

```bash
mvn spring-boot:run
```

The application will start on `http://localhost:8080`

### Running Tests

```bash
mvn test
```

## API Endpoints

### Domain Data APIs

#### Countries
- `GET /api/domain/countries` - Get all countries
- `GET /api/domain/countries/{id}` - Get country by ID
- `GET /api/domain/countries/code/{code}` - Get country by code
- `POST /api/domain/countries` - Create country
- `PUT /api/domain/countries/{id}` - Update country
- `DELETE /api/domain/countries/{id}` - Delete country

#### Currencies
- `GET /api/domain/currencies` - Get all currencies
- `GET /api/domain/currencies/{id}` - Get currency by ID
- `GET /api/domain/currencies/code/{code}` - Get currency by code
- `POST /api/domain/currencies` - Create currency
- `PUT /api/domain/currencies/{id}` - Update currency
- `DELETE /api/domain/currencies/{id}` - Delete currency

#### Entities
- `GET /api/domain/entities` - Get all entities
- `GET /api/domain/entities/{id}` - Get entity by ID
- `GET /api/domain/entities/registration/{registrationNumber}` - Get entity by registration number
- `POST /api/domain/entities` - Create entity
- `PUT /api/domain/entities/{id}` - Update entity
- `DELETE /api/domain/entities/{id}` - Delete entity

#### Instruments
- `GET /api/domain/instruments` - Get all instruments
- `GET /api/domain/instruments/{id}` - Get instrument by ID
- `GET /api/domain/instruments/isin/{isin}` - Get instrument by ISIN
- `POST /api/domain/instruments` - Create instrument
- `PUT /api/domain/instruments/{id}` - Update instrument
- `DELETE /api/domain/instruments/{id}` - Delete instrument

#### Accounts
- `GET /api/domain/accounts` - Get all accounts
- `GET /api/domain/accounts/{id}` - Get account by ID
- `GET /api/domain/accounts/number/{accountNumber}` - Get account by number
- `POST /api/domain/accounts` - Create account
- `PUT /api/domain/accounts/{id}` - Update account
- `DELETE /api/domain/accounts/{id}` - Delete account

#### SSI's
- `GET /api/domain/ssis` - Get all SSI's
- `GET /api/domain/ssis/{id}` - Get SSI by ID
- `POST /api/domain/ssis` - Create SSI
- `PUT /api/domain/ssis/{id}` - Update SSI
- `DELETE /api/domain/ssis/{id}` - Delete SSI

### Data Acquisition APIs

- `POST /api/data/acquire` - Acquire data from external sources
- `POST /api/data/distribute` - Distribute data to external destinations

### Notification APIs

- `POST /api/notifications/email` - Send email notification
- `POST /api/notifications/sms` - Send SMS notification
- `POST /api/notifications/system` - Send system notification

### Service Discovery APIs

- `GET /api/gateway/services` - Get all registered services
- `GET /api/gateway/services/{serviceId}` - Get service by ID
- `POST /api/gateway/services` - Register new service
- `DELETE /api/gateway/services/{serviceId}` - Deregister service

## Configuration

Key configuration in `application.yml`:

```yaml
spring:
  application:
    name: financial-services
  datasource:
    url: jdbc:h2:mem:financialdb
  rabbitmq:
    host: localhost
    port: 5672

server:
  port: 8080
```

## Database Console

H2 Console is available at: `http://localhost:8080/h2-console`

- JDBC URL: `jdbc:h2:mem:financialdb`
- Username: `sa`
- Password: (empty)

## Example Usage

### Creating a Country

```bash
curl -X POST http://localhost:8080/api/domain/countries \
  -H "Content-Type: application/json" \
  -d '{
    "code": "US",
    "name": "United States",
    "alpha3Code": "USA",
    "region": "North America",
    "active": true
  }'
```

### Creating a Currency

```bash
curl -X POST http://localhost:8080/api/domain/currencies \
  -H "Content-Type: application/json" \
  -d '{
    "code": "USD",
    "name": "US Dollar",
    "symbol": "$",
    "decimalPlaces": 2,
    "active": true
  }'
```

### Sending a Notification

```bash
curl -X POST http://localhost:8080/api/notifications/email \
  -H "Content-Type: application/json" \
  -d '{
    "recipient": "user@example.com",
    "subject": "Test Notification",
    "message": "This is a test notification"
  }'
```

## Development

### Project Structure

- `src/main/java` - Application source code
- `src/main/resources` - Configuration files
- `src/test/java` - Test source code

### Adding New Domain Entities

1. Create entity model in `domain.model` package
2. Create repository interface in `domain.repository` package
3. Create service in `domain.service` package
4. Create controller in `api.controller` package

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
