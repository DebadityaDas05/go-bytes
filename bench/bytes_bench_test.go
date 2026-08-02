package main

import (
	"testing"

	bytesutil "github.com/DebadityaDas05/go-bytes/src"
)

func BenchmarkParseString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bytesutil.Parse("1.5MB")
	}
}

func BenchmarkParseInt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bytesutil.Parse(1572864)
	}
}

func BenchmarkFormat(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bytesutil.Format(1572864.0, nil)
	}
}

func BenchmarkFormatWithOptions(b *testing.B) {
	sep := " "
	opts := &bytesutil.Options{
		ThousandsSeparator: &sep,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bytesutil.Format(1000.0, opts)
	}
}

func BenchmarkBytes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bytesutil.Bytes("1000", nil)
	}
}
