package managementserver

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	record "github.com/XxRoloxX/dns/pkg/dns_record"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DB_HOST_KEY     = "DB_HOST"
	DB_USER_KEY     = "DB_USER"
	DB_PASSWORD_KEY = "DB_PASSWORD"
	DB_NAME_KEY     = "DB_NAME"
	DB_PORT_KEY     = "DB_PORT"
)

// DNS Record Type Enums
type ManagedDNSRecordType string

const (
	ManagedDNSRecordType_A     ManagedDNSRecordType = "A"
	ManagedDNSRecordType_AAAA  ManagedDNSRecordType = "AAAA"
	ManagedDNSRecordType_CNAME ManagedDNSRecordType = "CNAME"
	ManagedDNSRecordType_MX    ManagedDNSRecordType = "MX"
	ManagedDNSRecordType_TXT   ManagedDNSRecordType = "TXT"
	ManagedDNSRecordType_NS    ManagedDNSRecordType = "NS"
	ManagedDNSRecordType_SOA   ManagedDNSRecordType = "SOA"
)

func ConvertRecordTypeToCode(recordType ManagedDNSRecordType) (uint16, error) {
	switch recordType {
	case ManagedDNSRecordType_A:
		return record.ResourceRecordType__A, nil
	case ManagedDNSRecordType_AAAA:
		return record.ResourceRecordType__AAAA, nil
	case ManagedDNSRecordType_MX:
		return record.ResourceRecordType__MX, nil
	case ManagedDNSRecordType_TXT:
		return record.ResourceRecordType__TXT, nil
	default:
		return 0, errors.New("invalid RecordType")
	}
}

// DNS Record Class Enums
type ManagedDNSRecordClass string

func (c *ManagedDNSRecordClass) ConvertToResourceRecordClass() (record.ResourceRecordClass, error) {
	switch *c {
	case ManagedDNSRecordClass_IN:
		return record.ResourceRecordClass__In, nil

	case ManagedDNSRecordClass_CH:
		return record.ResourceRecordClass__Ch, nil
	case ManagedDNSRecordClass_REVIEW:
		return record.ResourceRecordClass__Review, nil
	default:
		return 0, errors.New(fmt.Sprintf("Invalid managed resource record class code: %d", c))
	}
}

const (
	ManagedDNSRecordClass_IN     ManagedDNSRecordClass = "IN"     // Internet
	ManagedDNSRecordClass_CS     ManagedDNSRecordClass = "CS"     // CSNET
	ManagedDNSRecordClass_CH     ManagedDNSRecordClass = "CH"     // CHAOS
	ManagedDNSRecordClass_HS     ManagedDNSRecordClass = "HS"     // Hesiod
	ManagedDNSRecordClass_REVIEW ManagedDNSRecordClass = "REVIEW" // REVIEW
)

type ManagedDNSResourceRecord struct {
	ID    int                   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string                `gorm:"not null" json:"name"`
	Type  ManagedDNSRecordType  `gorm:"not null" json:"type"`
	Class ManagedDNSRecordClass `gorm:"not null" json:"class"`
	Data  string                `gorm:"not null" json:"data"`
}

func (r *ManagedDNSResourceRecord) ConvertToResourceRecord() (record.ResourceRecord, error) {
	names := strings.Split(r.Name, ".")

	class, err := r.Class.ConvertToResourceRecordClass()
	if err != nil {
		return nil, err
	}

	switch r.Type {
	case ManagedDNSRecordType_A:
		address := net.ParseIP(r.Data)
		if address == nil || address.To4() == nil {
			return nil, errors.New("invalid IPv4 address in Data")
		}
		return record.NewARecord(names, class, address), nil

	case ManagedDNSRecordType_AAAA:
		address := net.ParseIP(r.Data)
		if address == nil || address.To16() == nil {
			return nil, errors.New("invalid IPv6 address in Data")
		}

		return record.NewAAAARecord(names, class, address), nil

	case ManagedDNSRecordType_CNAME:
		domain := strings.Split(r.Data, ".") // Assuming the domain is comma-separated
		return record.NewCNAMERecord(names, class, domain), nil

	case ManagedDNSRecordType_TXT:
		return record.NewTXTRecord(names, class, []byte(r.Data)), nil

	case ManagedDNSRecordType_MX:
		return record.NewMXRecord(names, class, []byte(r.Data)), nil

	default:
		return nil, errors.New("unsupported record type")
	}
}

type RecordsRepository interface {
	GetRecords() ([]ManagedDNSResourceRecord, error)
	GetRecordsByName(name string) ([]ManagedDNSResourceRecord, error)
	CreateRecord(record *ManagedDNSResourceRecord) error
	DeleteRecord(id int) error
}

type PostgresRecordsRepository struct {
	db *gorm.DB
}

func (r *PostgresRecordsRepository) GetRecords() ([]ManagedDNSResourceRecord, error) {
	var records []ManagedDNSResourceRecord
	if err := r.db.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (r *PostgresRecordsRepository) GetRecordsByName(name string) ([]ManagedDNSResourceRecord, error) {
	var records []ManagedDNSResourceRecord
	if err := r.db.Find(&records, ManagedDNSResourceRecord{
		Name: name,
	}).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (r *PostgresRecordsRepository) CreateRecord(record *ManagedDNSResourceRecord) error {
	if err := r.db.Create(record).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgresRecordsRepository) DeleteRecord(id int) error {
	if err := r.db.Delete(&ManagedDNSResourceRecord{}, id).Error; err != nil {
		return err
	}
	return nil
}

func createConnectionString() string {
	host := os.Getenv(DB_HOST_KEY)
	user := os.Getenv(DB_USER_KEY)
	password := os.Getenv(DB_PASSWORD_KEY)
	dbname := os.Getenv(DB_NAME_KEY)
	port := os.Getenv(DB_PORT_KEY)

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
}

func NewPostgresRecordsRepository() *PostgresRecordsRepository {
	db, err := gorm.Open(
		postgres.Open(createConnectionString()),
		&gorm.Config{},
	)

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %s", err))
	}

	if err := db.AutoMigrate(&ManagedDNSResourceRecord{}); err != nil {
		panic(fmt.Sprintf("Failed to migrate database schema: %s", err))
	}

	return &PostgresRecordsRepository{
		db: db,
	}
}
