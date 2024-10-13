//go:build !(amd64 || arm64 || riscv64) || purego

// Code generated by ./cmd/ch-gen-col, DO NOT EDIT.

package proto

import (
	"encoding/binary"

	"github.com/go-faster/errors"
)

var _ = binary.LittleEndian // clickHouse uses LittleEndian

// DecodeColumn decodes Enum8 rows from *Reader.
func (c *ColEnum8) DecodeColumn(r *Reader, rows int) error {
	if rows == 0 {
		return nil
	}
	data, err := r.ReadRaw(rows)
	if err != nil {
		return errors.Wrap(err, "read")
	}
	v := *c
	v = append(v, make([]Enum8, rows)...)
	for i := range data {
		v[i] = Enum8(data[i])
	}
	*c = v
	return nil
}

// EncodeColumn encodes Enum8 rows to *Buffer.
func (c ColEnum8) EncodeColumn(b *Buffer) {
	v := c
	if len(v) == 0 {
		return
	}
	start := len(b.Buf)
	b.Buf = append(b.Buf, make([]byte, len(v))...)
	for i := range v {
		b.Buf[i+start] = uint8(v[i])
	}
}

func (c ColEnum8) WriteColumn(w *Writer) {
	w.ChainBuffer(c.EncodeColumn)
}
