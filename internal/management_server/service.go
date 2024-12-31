package managementserver

import "errors"

type RecordsService struct {
	recordsRepository RecordsRepository
}

func NewRecordsService(recordsRepository RecordsRepository) *RecordsService {
	return &RecordsService{
		recordsRepository: recordsRepository,
	}
}

func (s *RecordsService) GetRecords() ([]ManagedDNSResourceRecord, error) {

	records, err := s.recordsRepository.GetRecords()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (s *RecordsService) CreateRecord(record *ManagedDNSResourceRecord) error {
	if record == nil {
		return errors.New("record cannot be nil")
	}

	return s.recordsRepository.CreateRecord(record)
}

func (s *RecordsService) DeleteRecord(id int) error {
	if id <= 0 {
		return errors.New("invalid record ID")
	}

	return s.recordsRepository.DeleteRecord(id)
}
