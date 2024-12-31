package record

import (
	"net"
	"strings"
)

type RRStore map[string][]ResourceRecord

func NewRRStore() *RRStore {
	return &RRStore{
		"a.subdomain.rolo-labs.xyz": []ResourceRecord{
			&ARecord{
				name:    []string{"a", "subdomain", "rolo-labs", "xyz"},
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
