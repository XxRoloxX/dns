package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecoder_decodeBody(t *testing.T) {

	testCases := []struct {
		name         string
		rawQuery     []byte
		expectedBody []Query
		expectedErr  error
	}{
		{
			name: "A message for example.com should be decoded into single query",
			rawQuery: []byte{
				0b00000111, 0b01100101, 0b01111000, 0b01100001, 0b01101101, 0b01110000, 0b01101100, 0b01100101, // "example"
				0b00000011, 0b01100011, 0b01101111, 0b01101101, // "com"
				0b00000000,             // Terminating null byte
				0b00000000, 0b00000001, // Query Type: AAAA
				0b00000000, 0b00000001, // Query Class: IN
			},
			expectedBody: []Query{
				{
					groups: []string{"example", "com"},
					t:      ResourceRecordType__AAAA,
					class:  ResourceRecordClass__In,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			header := Header{
				numberOfQuestions: 1,
			}

			query, err := NewDecoder(tc.rawQuery).decodeBody(&header)
			if err != nil && tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tc.expectedBody, query)
		})
	}
}

func TestDecoder_decodeHeader(t *testing.T) {

	testCases := []struct {
		name           string
		rawQuery       []byte
		expectedHeader Header
		expectedErr    error
	}{
		{
			name: "A message for example.com should be decoded into single query",
			rawQuery: []byte{
				0b00011010, 0b00101011, // Transaction ID
				0b00000001, 0b00000000, // Flags
				0b00000000, 0b00000001, // Questions
				0b00000000, 0b00000000, // Answer RRs
				0b00000000, 0b00000000, // Authority RRs
				0b00000000, 0b00000000, // Additional RRs
			},
			expectedHeader: Header{
				transactionId: 0b0001101000101011,
				flags: HeaderFlags{
					query:              true,
					operationCode:      OpCode__Query,
					authorativeAnswer:  false,
					truncation:         false,
					recursionDesired:   true,
					recursionAvailable: false,
				},
				numberOfQuestions:    1,
				numberOfAnswers:      0,
				numberOfAuthorityRR:  0,
				numberOfAdditionalRR: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			header, err := NewDecoder(tc.rawQuery).decodeHeader()
			if err != nil && tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tc.expectedHeader, *header)
		})
	}
}
