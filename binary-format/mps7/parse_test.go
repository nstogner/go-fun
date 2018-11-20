package mps7_test

import (
	"os"
	"testing"
	"time"

	"github.com/nstogner/go-fun/binary-format/mps7"
)

func TestReadFile(t *testing.T) {
	f, err := os.Open("./txnlog.dat")
	if err != nil {
		t.Fatalf("reading file from disk: %s", err)
	}
	defer f.Close()

	mps7f, err := mps7.ReadFile(f)
	if err != nil {
		t.Fatalf("reading file as mps7 format: %s", err)
	}

	if err := mps7f.Validate(); err != nil {
		t.Fatalf("unxpected validation error: %s", err)
	}

	if exp, got := int(mps7f.Header.RecordCount), len(mps7f.Records); exp != got {
		t.Fatalf("expected header record count %v to match length of records %v", exp, got)
	}

	if n := len(mps7f.Records); n < 1 {
		t.Fatalf("expected at least one record in test file, got %v", n)
	}

	// Expected:
	// | 'Debit'     | 1393108945     | 4136353673894269217 | 604.274335557087  |
	expRec0 := mps7.Record{
		Type:      mps7.RecordTypeDebit,
		Timestamp: time.Unix(1393108945, 0),
		UserID:    4136353673894269217,
		Amount:    604.274335557087,
	}
	if exp, got := expRec0, mps7f.Records[0]; exp != got {
		t.Fatalf("expected record[0] = %+v, got %+v", exp, got)
	}
	if exp, got := int64(1393108945), mps7f.Records[0].Timestamp.Unix(); exp != got {
		t.Fatalf("expected unix record[0].timestamp = %v, got %v", exp, got)
	}
}
