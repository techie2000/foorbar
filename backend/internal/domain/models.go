package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Country represents a country entity
type Country struct {
	BaseModel
	Code       string `gorm:"uniqueIndex;size:2;not null" json:"code" validate:"required,len=2"`
	Name       string `gorm:"not null" json:"name" validate:"required"`
	Alpha3Code string `gorm:"size:3" json:"alpha3_code" validate:"len=3"`
	Region     string `json:"region"`
	Active     bool   `gorm:"default:true" json:"active"`
}

// TableName overrides the table name
func (Country) TableName() string {
	return "countries"
}

// Currency represents a currency entity
type Currency struct {
	BaseModel
	Code          string `gorm:"uniqueIndex;size:3;not null" json:"code" validate:"required,len=3"`
	Name          string `gorm:"not null" json:"name" validate:"required"`
	Symbol        string `json:"symbol"`
	DecimalPlaces int    `gorm:"default:2" json:"decimal_places"`
	Active        bool   `gorm:"default:true" json:"active"`
}

// TableName overrides the table name
func (Currency) TableName() string {
	return "currencies"
}

// Address represents a physical address (ISO20022 compliant)
type Address struct {
	BaseModel
	// Structured address fields (ISO20022)
	AddressType        string     `gorm:"size:50" json:"address_type,omitempty"`         // Type of address (ADDR, PBOX, HOME, BIZZ)
	Department         string     `gorm:"size:70" json:"department,omitempty"`           // Department
	SubDepartment      string     `gorm:"size:70" json:"sub_department,omitempty"`       // Sub-department
	StreetName         string     `gorm:"size:70" json:"street_name,omitempty"`          // Street name
	BuildingNumber     string     `gorm:"size:16" json:"building_number,omitempty"`      // Building number
	BuildingName       string     `gorm:"size:35" json:"building_name,omitempty"`        // Building name
	Floor              string     `gorm:"size:70" json:"floor,omitempty"`                // Floor
	PostBox            string     `gorm:"size:16" json:"post_box,omitempty"`             // Post office box number
	Room               string     `gorm:"size:70" json:"room,omitempty"`                 // Room
	PostalCode         string     `gorm:"size:16" json:"postal_code,omitempty"`          // Postal code/ZIP code
	TownName           string     `gorm:"size:35" json:"town_name,omitempty"`            // Town/city name
	TownLocationName   string     `gorm:"size:35" json:"town_location_name,omitempty"`   // Town location name
	DistrictName       string     `gorm:"size:35" json:"district_name,omitempty"`        // District name
	CountrySubDivision string     `gorm:"size:35" json:"country_sub_division,omitempty"` // State/province/region
	CountryID          *uuid.UUID `gorm:"type:uuid" json:"country_id,omitempty"`
	Country            *Country   `gorm:"foreignKey:CountryID" json:"country,omitempty"`
	// Unstructured address lines (ISO20022 allows up to 7 lines)
	AddressLine1 string `gorm:"size:70" json:"address_line_1,omitempty"`
	AddressLine2 string `gorm:"size:70" json:"address_line_2,omitempty"`
	AddressLine3 string `gorm:"size:70" json:"address_line_3,omitempty"`
	AddressLine4 string `gorm:"size:70" json:"address_line_4,omitempty"`
	AddressLine5 string `gorm:"size:70" json:"address_line_5,omitempty"`
	AddressLine6 string `gorm:"size:70" json:"address_line_6,omitempty"`
	AddressLine7 string `gorm:"size:70" json:"address_line_7,omitempty"`
}

// TableName overrides the table name
func (Address) TableName() string {
	return "addresses"
}

// EntityType represents the type of an entity
type EntityType string

const (
	EntityTypeCompany      EntityType = "COMPANY"
	EntityTypeBusiness     EntityType = "BUSINESS"
	EntityTypeCorporation  EntityType = "CORPORATION"
	EntityTypePartnership  EntityType = "PARTNERSHIP"
	EntityTypeIndividual   EntityType = "INDIVIDUAL"
)

// Entity represents a business entity (company, individual, etc.)
type Entity struct {
	BaseModel
	Name               string           `gorm:"not null" json:"name" validate:"required"`
	RegistrationNumber string           `gorm:"uniqueIndex" json:"registration_number"`
	Type               EntityType       `gorm:"type:varchar(50)" json:"type"`
	Addresses          []EntityAddress  `gorm:"foreignKey:EntityID" json:"addresses,omitempty"`
	Active             bool             `gorm:"default:true" json:"active"`
}

// TableName overrides the table name
func (Entity) TableName() string {
	return "entities"
}

// EntityAddress represents the many-to-many relationship between entities and addresses
type EntityAddress struct {
	BaseModel
	EntityID    uuid.UUID  `gorm:"type:uuid;not null" json:"entity_id"`
	Entity      *Entity    `gorm:"foreignKey:EntityID" json:"entity,omitempty"`
	AddressID   uuid.UUID  `gorm:"type:uuid;not null" json:"address_id"`
	Address     *Address   `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	AddressType string     `gorm:"size:50" json:"address_type,omitempty"` // e.g., 'REGISTERED', 'TRADING', 'BILLING', 'CORRESPONDENCE'
	IsPrimary   bool       `gorm:"default:false" json:"is_primary"`
}

// TableName overrides the table name
func (EntityAddress) TableName() string {
	return "entity_addresses"
}

// InstrumentType represents the type of financial instrument
type InstrumentType string

const (
	InstrumentTypeEquity     InstrumentType = "EQUITY"
	InstrumentTypeBond       InstrumentType = "BOND"
	InstrumentTypeDerivative InstrumentType = "DERIVATIVE"
	InstrumentTypeCommodity  InstrumentType = "COMMODITY"
	InstrumentTypeFund       InstrumentType = "FUND"
	InstrumentTypeForex      InstrumentType = "FOREX"
)

// Instrument represents a financial instrument
type Instrument struct {
	BaseModel
	Name              string             `gorm:"not null" json:"name" validate:"required"`
	Type              InstrumentType     `gorm:"type:varchar(50)" json:"type"`
	IssueCurrencyID   *uuid.UUID         `gorm:"type:uuid;column:issue_currency_id" json:"issue_currency_id"`
	IssueCurrency     *Currency          `gorm:"foreignKey:IssueCurrencyID" json:"issue_currency,omitempty"`
	PrimaryExchange   string             `gorm:"column:primary_exchange" json:"primary_exchange"`
	Codes             []InstrumentCode   `gorm:"foreignKey:InstrumentID" json:"codes,omitempty"`
	Active            bool               `gorm:"default:true" json:"active"`
}

// TableName overrides the table name
func (Instrument) TableName() string {
	return "instruments"
}

// IdentifierLevel represents the level of instrument identifier
type IdentifierLevel string

const (
	IdentifierLevelInternational IdentifierLevel = "INTERNATIONAL" // ISIN, FIGI
	IdentifierLevelRegional      IdentifierLevel = "REGIONAL"      // CUSIP, WKN
	IdentifierLevelLocal         IdentifierLevel = "LOCAL"         // Exchange-specific codes
)

// CodeType represents the type of instrument code
type CodeType string

const (
	CodeTypeISIN      CodeType = "ISIN"
	CodeTypeFIGI      CodeType = "FIGI"
	CodeTypeCUSIP     CodeType = "CUSIP"
	CodeTypeWKN       CodeType = "WKN"
	CodeTypeSEDOL     CodeType = "SEDOL"
	CodeTypeRIC       CodeType = "RIC"
	CodeTypeTicker    CodeType = "TICKER"
	CodeTypeBloomberg CodeType = "BLOOMBERG"
)

// InstrumentCode represents an identifier code for an instrument
type InstrumentCode struct {
	BaseModel
	InstrumentID          uuid.UUID       `gorm:"type:uuid;not null" json:"instrument_id"`
	Instrument            *Instrument     `gorm:"foreignKey:InstrumentID" json:"instrument,omitempty"`
	CodeType              CodeType        `gorm:"type:varchar(50);not null" json:"code_type"`
	CodeValue             string          `gorm:"size:100;not null" json:"code_value"`
	IdentifierLevel       IdentifierLevel `gorm:"type:varchar(50)" json:"identifier_level,omitempty"`
	MarketIdentifierCode  string          `gorm:"size:10" json:"market_identifier_code,omitempty"` // MIC code (e.g., XNAS, XFRA)
	Region                string          `gorm:"size:50" json:"region,omitempty"`                 // For regional codes (e.g., US, DE)
	IsPrimary             bool            `gorm:"default:false" json:"is_primary"`
}

// TableName overrides the table name
func (InstrumentCode) TableName() string {
	return "instrument_codes"
}

// AccountType represents the type of account
type AccountType string

const (
	AccountTypeTrading    AccountType = "TRADING"
	AccountTypeSettlement AccountType = "SETTLEMENT"
	AccountTypeCustody    AccountType = "CUSTODY"
	AccountTypeMargin     AccountType = "MARGIN"
)

// Account represents a financial account
type Account struct {
	BaseModel
	AccountNumber      string       `gorm:"uniqueIndex;not null" json:"account_number" validate:"required"`
	EntityID           *uuid.UUID   `gorm:"type:uuid" json:"entity_id"`
	Entity             *Entity      `gorm:"foreignKey:EntityID" json:"entity,omitempty"`
	AccountCurrencyID  *uuid.UUID   `gorm:"type:uuid;column:account_currency_id" json:"account_currency_id"`
	AccountCurrency    *Currency    `gorm:"foreignKey:AccountCurrencyID" json:"account_currency,omitempty"`
	Type               AccountType  `gorm:"type:varchar(50)" json:"type"`
	Balance            float64      `gorm:"type:decimal(19,4);default:0" json:"balance"`
	OpenedAt           time.Time    `json:"opened_at"`
	Active             bool         `gorm:"default:true" json:"active"`
}

// TableName overrides the table name
func (Account) TableName() string {
	return "accounts"
}

// SettlementType represents the type of settlement
type SettlementType string

const (
	SettlementTypeDVP SettlementType = "DVP" // Delivery Versus Payment
	SettlementTypeFOP SettlementType = "FOP" // Free Of Payment
	SettlementTypeRVP SettlementType = "RVP" // Receive Versus Payment
	SettlementTypeDAP SettlementType = "DAP" // Delivery Against Payment
)

// SSI represents Standard Settlement Instructions
type SSI struct {
	BaseModel
	EntityID               *uuid.UUID      `gorm:"type:uuid" json:"entity_id"`
	Entity                 *Entity         `gorm:"foreignKey:EntityID" json:"entity,omitempty"`
	SettlementCurrencyID   *uuid.UUID      `gorm:"type:uuid;column:settlement_currency_id" json:"settlement_currency_id"`
	SettlementCurrency     *Currency       `gorm:"foreignKey:SettlementCurrencyID" json:"settlement_currency,omitempty"`
	InstrumentID           *uuid.UUID      `gorm:"type:uuid" json:"instrument_id"`
	Instrument             *Instrument     `gorm:"foreignKey:InstrumentID" json:"instrument,omitempty"`
	BeneficiaryName        string          `gorm:"not null" json:"beneficiary_name" validate:"required"`
	BeneficiaryAccount     string          `gorm:"not null" json:"beneficiary_account" validate:"required"`
	BeneficiaryBank        string          `gorm:"not null" json:"beneficiary_bank" validate:"required"`
	BeneficiaryBankBIC     string          `json:"beneficiary_bank_bic"`
	IntermediaryBank       string          `json:"intermediary_bank"`
	IntermediaryBankBIC    string          `json:"intermediary_bank_bic"`
	SettlementType         SettlementType  `gorm:"type:varchar(50)" json:"settlement_type"`
	ValidFrom              time.Time       `json:"valid_from"`
	ValidTo                *time.Time      `json:"valid_to"`
	Active                 bool            `gorm:"default:true" json:"active"`
}

// TableName overrides the table name
func (SSI) TableName() string {
	return "ssis"
}

// AuditLog represents an audit trail entry
type AuditLog struct {
	BaseModel
	EntityType   string    `gorm:"not null" json:"entity_type"`
	EntityID     uuid.UUID `gorm:"type:uuid;not null" json:"entity_id"`
	Action       string    `gorm:"not null" json:"action"` // CREATE, UPDATE, DELETE
	UserID       uuid.UUID `gorm:"type:uuid" json:"user_id"`
	ChangedData  string    `gorm:"type:jsonb" json:"changed_data"`
	PreviousData string    `gorm:"type:jsonb" json:"previous_data"`
}

// TableName overrides the table name
func (AuditLog) TableName() string {
	return "audit_logs"
}
