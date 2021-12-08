package proto

import (
	"bytes"
	"testing"
)

func BenchmarkColumnUInt32_DecodeColumn(b *testing.B) {
	const rows = 50_000
	var data ColumnUInt32
	for i := 0; i < rows; i++ {
		data = append(data, uint32(i))
	}

	var buf Buffer
	data.EncodeColumn(&buf)

	br := bytes.NewReader(buf.Buf)
	r := NewReader(br)

	b.SetBytes(int64(len(buf.Buf)))
	b.ResetTimer()
	b.ReportAllocs()

	var dec ColumnUInt32
	for i := 0; i < b.N; i++ {
		br.Reset(buf.Buf)
		dec.Reset()
		if err := dec.DecodeColumn(r, rows); err != nil {
			b.Fatal(err)
		}
	}
}
