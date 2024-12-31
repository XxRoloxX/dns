package record

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type ResourceRecord interface {
	Name() []string
	Class() ResourceRecordClass
	Type() ResourceRecordType
	Data() []byte
}

// RR pointing to IPv4 address
type ARecord struct {
	name    []string
	class   ResourceRecordClass
	address net.IP
}

func NewARecord(name []string, class ResourceRecordClass, address net.IP) *ARecord {
	return &ARecord{
		name:    name,
		class:   class,
		address: address,
	}
}

func (r *ARecord) Name() []string {
	return r.name
}

func (r *ARecord) Class() ResourceRecordClass {
	return r.class
}

func (r *ARecord) Type() ResourceRecordType {
	return ResourceRecordType__A
}

func (r *ARecord) Data() []byte {
	return r.address.To4()
}

// RR pointing to IPv6 address
type AAAARecord struct {
	name    []string
	class   ResourceRecordClass
	address net.IP
}

func NewAAAARecord(name []string, class ResourceRecordClass, address net.IP) *AAAARecord {
	return &AAAARecord{
		name:    name,
		class:   class,
		address: address,
	}
}

func (r *AAAARecord) Name() []string {
	return r.name
}

func (r *AAAARecord) Class() ResourceRecordClass {
	return r.class
}

func (r *AAAARecord) Type() ResourceRecordType {
	return ResourceRecordType__AAAA
}

func (r *AAAARecord) Data() []byte {
	return r.address.To4()
}

// RR pointing to another domain
type CNAMERecord struct {
	name   []string
	class  ResourceRecordClass
	domain []string
}

func (r *CNAMERecord) Name() []string {
	return r.name
}

func (r *CNAMERecord) Class() ResourceRecordClass {
	return r.class
}

func (r *CNAMERecord) Type() ResourceRecordType {
	return ResourceRecordType__CNAME
}

func (r *CNAMERecord) Data() []byte {
	return []byte(strings.Join(r.domain, ""))
}

func NewCNAMERecord(name []string, class ResourceRecordClass, domain []string) *CNAMERecord {
	return &CNAMERecord{
		name:   name,
		class:  class,
		domain: domain,
	}
}

// RR pointing to raw data
type TXTRecord struct {
	name  []string
	class ResourceRecordClass
	data  []byte
}

func (r *TXTRecord) Name() []string {
	return r.name
}

func (r *TXTRecord) Class() ResourceRecordClass {
	return r.class
}

func (r *TXTRecord) Type() ResourceRecordType {
	return ResourceRecordType__TXT
}

func (r *TXTRecord) Data() []byte {
	return r.data
}

func NewTXTRecord(name []string, class ResourceRecordClass, data []byte) *TXTRecord {
	return &TXTRecord{
		name:  name,
		class: class,
		data:  data,
	}
}

// RR pointing to raw data
type MXRecord struct {
	name  []string
	class ResourceRecordClass
	data  []byte
}

func (r *MXRecord) Name() []string {
	return r.name
}

func (r *MXRecord) Class() ResourceRecordClass {
	return r.class
}

func (r *MXRecord) Type() ResourceRecordType {
	return ResourceRecordType__TXT
}

func (r *MXRecord) Data() []byte {
	return r.data
}

func NewMXRecord(name []string, class ResourceRecordClass, data []byte) *MXRecord {
	return &MXRecord{
		name:  name,
		class: class,
		data:  data,
	}
}

type ResourceRecordType uint

const (
	ResourceRecordType__A     = 1
	ResourceRecordType__AAAA  = 2
	ResourceRecordType__MX    = 3
	ResourceRecordType__TXT   = 4
	ResourceRecordType__CNAME = 5
)

func NewResourceRecordType(code uint16) (ResourceRecordType, error) {
	switch code {
	case 1:
		return ResourceRecordType__A, nil
	case 2:
		return ResourceRecordType__AAAA, nil
	case 3:
		return ResourceRecordType__MX, nil
	case 4:
		return ResourceRecordType__TXT, nil
	default:
		return 0, errors.New("Invalid resource record type code")
	}
}

type ResourceRecordClass uint

const (
	ResourceRecordClass__In     = 1
	ResourceRecordClass__Ch     = 2
	ResourceRecordClass__Review = 256
)

func NewResourceRecordClass(code uint16) (ResourceRecordClass, error) {
	switch code {
	case 1:
		return ResourceRecordClass__In, nil
	case 256:
		return ResourceRecordClass__Review, nil
	default:
		return 0, errors.New(fmt.Sprintf("Invalid resource record class code: %d", code))
	}
}
