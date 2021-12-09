{{- /*gotype: github.com/go-faster/ch/internal/proto/cmd/ch-gen-int.Variant*/ -}}
// Code generated by ch-gen-int, DO NOT EDIT.

package proto

import "github.com/go-faster/errors"

type {{ .Type }} []{{ .ElemType }}

func ({{ .Type }}) Type() ColumnType { return {{ .ColumnType }} }
func (c {{ .Type }}) Rows() int      { return len(c) }
func (c *{{ .Type }}) Reset()        { *c = (*c)[:0] }

func (c {{ .Type }}) EncodeColumn(b *Buffer) {
  for _, v := range c {
    b.Put{{ .Name }}(v)
  }
}

func (c *{{ .Type }}) DecodeColumn(r *Reader, rows int) error {
  const size = {{ .Bits }} / 8
  data, err := r.ReadRaw(rows * size)
  if err != nil {
    return errors.Wrap(err, "read")
  }
  {{ if .Byte }}
  *c = append(*c, data...)
  {{ else }}
  v := *c
  for i := 0; i < len(data); i += size {
    v = append(v,
    {{- if .Signed }}
      {{- if eq .Bits 8 }}
       {{ .ElemType }}(data[i]),
      {{- else }}
        {{ .ElemType }}(bin.{{ .BinFunc }}(data[i:i+size])),
      {{- end }}
    {{- else }}
      bin.{{ .BinFunc }}(data[i:i+size]),
    {{- end }}
    )
  }
  *c = v
  {{ end }}
  return nil
}
