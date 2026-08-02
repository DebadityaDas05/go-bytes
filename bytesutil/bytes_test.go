package bytesutil_test

import (
	"testing"

	"github.com/DebadityaDas05/go-bytes/bytesutil"
)

func TestParseNative(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"1KB", 1024.0},
		{"1.5MB", 1572864.0},
		{"100", 100.0},
		{1024, 1024.0},
		{1024.5, 1024.5},
		{"invalid", nil},
		{nil, nil},
	}

	for _, tt := range tests {
		res := bytesutil.Parse(tt.input)
		if res != tt.expected {
			t.Errorf("Parse(%v) = %v; want %v", tt.input, res, tt.expected)
		}
	}
}

func TestFormatNative(t *testing.T) {
	sep := " "
	unitSep := " "
	opts := &bytesutil.Options{
		ThousandsSeparator: &sep,
		UnitSeparator:      &unitSep,
	}

	tests := []struct {
		value    float64
		opts     *bytesutil.Options
		expected interface{}
	}{
		{1024, nil, "1KB"},
		{1000, opts, "1 000 B"},
		{1572864, nil, "1.5MB"},
	}

	for _, tt := range tests {
		res := bytesutil.Format(tt.value, tt.opts)
		if res != tt.expected {
			t.Errorf("Format(%v) = %v; want %v", tt.value, res, tt.expected)
		}
	}
}

func TestBytesNative(t *testing.T) {
	if res := bytesutil.Bytes("1KB", nil); res != 1024.0 {
		t.Errorf("Bytes(\"1KB\") = %v; want 1024.0", res)
	}
	if res := bytesutil.Bytes(1024.0, nil); res != "1KB" {
		t.Errorf("Bytes(1024.0) = %v; want \"1KB\"", res)
	}
}
