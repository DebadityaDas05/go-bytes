package bytesutil_test

import (
	"testing"

	bytesutil "github.com/DebadityaDas05/go-bytes/src"
)

func TestParseNative(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"1KB", int64(1024)},
		{"1.5MB", int64(1572864)},
		{"100", int64(100)},
		{1024, int64(1024)},
		{1024.5, 1024.5},
		{"invalid", nil},
		{nil, nil},
	}

	for _, tt := range tests {
		res := bytesutil.Parse(tt.input)
		if res != tt.expected {
			t.Errorf("Parse(%v) = %v (%T); want %v (%T)", tt.input, res, res, tt.expected, tt.expected)
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
	if res := bytesutil.Bytes("1KB", nil); res != int64(1024) {
		t.Errorf("Bytes(\"1KB\") = %v; want 1024", res)
	}
	if res := bytesutil.Bytes(1024.0, nil); res != "1KB" {
		t.Errorf("Bytes(1024.0) = %v; want \"1KB\"", res)
	}
}
