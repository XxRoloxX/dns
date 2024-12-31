package record

import (
	"errors"
	"fmt"
	"net"
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

// RR pointing to another domain
type CNAMERecord struct {
	name   []string
	class  ResourceRecordClass
	domain []string
}

// RR pointing to raw data
type TXTRecord struct {
	name  []string
	class ResourceRecordClass
	data  []string
}

// RR pointing to raw data
type MXRecord struct {
	name  []string
	class ResourceRecordClass
	data  []string
}

type ResourceRecordType uint

const (
	ResourceRecordType__A    = 1
	ResourceRecordType__AAAA = 2
	ResourceRecordType__MX   = 3
	ResourceRecordType__TXT  = 4
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
