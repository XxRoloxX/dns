package record

import (
	"net"
	"strings"
)

type RRStore map[string][]ResourceRecord

func NewRRStore() *RRStore {
	return &RRStore{
		"dns.test.com": []ResourceRecord{
			&ARecord{
				name:    []string{"dns", "test", "com"},
				class:   ResourceRecordClass__In,
				address: net.IPv4(8, 8, 8, 8),
			},
		},
	}
}

func (s RRStore) ResourceRecords(name []string) []ResourceRecord {
	if records, ok := s[strings.Join(name, ".")]; ok == true {
		return records
	}

	return make([]ResourceRecord, 0)
}
