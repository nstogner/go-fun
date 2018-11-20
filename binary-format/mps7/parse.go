package mps7

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
)

// File is a mps7 file.
type File struct {
	Header  Header
	Records []Record
}

func (f *File) Validate() error {
	if err := f.Header.Validate(); err != nil {
		return errors.Wrap(err, "header")
	}

	return nil
}

// ReadFile reads an entire mps7 file.
func ReadFile(r io.Reader) (*File, error) {
	hdr, err := ReadHeader(r)
	if err != nil {
		return nil, errors.Wrap(err, "reading header")
	}

	f := File{
		Header:  hdr,
		Records: make([]Record, hdr.RecordCount),
	}

	for i := uint32(0); i < hdr.RecordCount; i++ {
		f.Records[i], err = ReadRecord(r)
		if err != nil {
			return nil, errors.Wrap(err, "reading record")
		}
		fmt.Println(f.Records[i].Timestamp.Unix())
	}

	return &f, nil
}

// Header gives a summary of the file.
type Header struct {
	Label string
	// Version is probably an integer, but the specification did not explicitly say.
	Version     byte
	RecordCount uint32
}

// Validate the header.
func (h Header) Validate() error {
	if exp, got := "MPS7", h.Label; exp != got {
		return fmt.Errorf("expected header label %q, got %q", exp, got)
	}
	return nil
}

// ReadHeader reads the header from an mps7 file. This should be the first read
// performed on a file.
func ReadHeader(r io.Reader) (Header, error) {
	// Format:
	// | 4 byte magic string "MPS7" | 1 byte version | 4 byte (uint32) # of records |

	var raw struct {
		Label       [4]byte
		Version     byte
		RecordCount uint32
	}

	if err := binary.Read(r, binary.BigEndian, &raw); err != nil {
		return Header{}, err
	}

	return Header{
		Label:       string(raw.Label[:]),
		Version:     raw.Version,
		RecordCount: raw.RecordCount,
	}, nil
}

const (
	RecordTypeDebit        = RecordType(0)
	RecordTypeCredit       = RecordType(1)
	RecordTypeStartAutopay = RecordType(2)
	RecordTypeEndAutopay   = RecordType(3)
)

type RecordType uint8

func (t RecordType) String() string {
	switch t {
	case RecordTypeDebit:
		return "DEBIT"
	case RecordTypeCredit:
		return "DEBIT"
	case RecordTypeStartAutopay:
		return "START_AUTOPAY"
	case RecordTypeEndAutopay:
		return "END_AUTOPAY"
	default:
		return "undefined"
	}
}

// Record represents a single line item in a mps7 file.
type Record struct {
	Type      RecordType
	Timestamp time.Time
	UserID    uint64
	// Amount is only present for Type = RecordTypeDebit || RecordTypeCredit
	Amount float64
}

// ReadRecord reads a single record.
func ReadRecord(r io.Reader) (Record, error) {
	// Format:
	// | 1 byte record type enum | 4 byte (uint32) Unix timestamp | 8 byte (uint64) user ID |
	// Example:
	// | Record type | Unix timestamp | user ID             | amount in dollars |
	// | 'Debit'     | 1393108945     | 4136353673894269217 | 604.274335557087  |

	//  Record Types:
	// * 0x00: Debit
	// * 0x01: Credit
	// * 0x02: StartAutopay
	// * 0x03: EndAutopay

	// Amounts:
	// For Debit and Credit record types, there is an additional field, an 8 byte
	// (float64) amount in dollars, at the end of the record.

	var raw struct {
		Type   uint8
		UnixTS uint32
		UserID uint64
	}

	if err := binary.Read(r, binary.BigEndian, &raw); err != nil {
		return Record{}, errors.Wrap(err, "reading main record")
	}

	var amt float64
	switch RecordType(raw.Type) {
	case RecordTypeDebit, RecordTypeCredit:
		if err := binary.Read(r, binary.BigEndian, &amt); err != nil {
			return Record{}, errors.Wrap(err, "reading amount")
		}
	}

	return Record{
		Type:      RecordType(raw.Type),
		Timestamp: time.Unix(int64(raw.UnixTS), 0),
		UserID:    raw.UserID,
		Amount:    amt,
	}, nil
}
