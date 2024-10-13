package proto

import (
	"strconv"
	"strings"

	"github.com/go-faster/errors"
)

// ColAuto is column that is initialized during decoding.
type ColAuto struct {
	Data     Column
	DataType ColumnType
}

// Infer and initialize Column from ColumnType.
func (c *ColAuto) Infer(t ColumnType) error {
	if c.Data != nil && !c.Type().Conflicts(t) {
		// Already ok.
		c.DataType = t // update subtype if needed
		return nil
	}
	if v := inferGenerated(t); v != nil {
		c.Data = v
		c.DataType = t
		return nil
	}
	if strings.HasPrefix(t.String(), ColumnTypeInterval.String()) {
		v := new(ColInterval)
		if err := v.Infer(t); err != nil {
			return errors.Wrap(err, "interval")
		}
		c.Data = v
		c.DataType = t
		return nil
	}
	switch t {
	case ColumnTypeNothing:
		c.Data = new(ColNothing)
	case ColumnTypeNullable.Sub(ColumnTypeNothing):
		c.Data = new(ColNothing).Nullable()
	case ColumnTypeArray.Sub(ColumnTypeNothing):
		c.Data = new(ColNothing).Array()
	case ColumnTypeString:
		c.Data = new(ColStr)
	case ColumnTypeArray.Sub(ColumnTypeString):
		c.Data = new(ColStr).Array()
	case ColumnTypeNullable.Sub(ColumnTypeString):
		c.Data = new(ColStr).Nullable()
	case ColumnTypeLowCardinality.Sub(ColumnTypeString):
		c.Data = new(ColStr).LowCardinality()
	case ColumnTypeArray.Sub(ColumnTypeLowCardinality.Sub(ColumnTypeString)):
		c.Data = new(ColStr).LowCardinality().Array()
	case ColumnTypeBool:
		c.Data = new(ColBool)
	case ColumnTypeDateTime:
		c.Data = new(ColDateTime)
	case ColumnTypeDate:
		c.Data = new(ColDate)
	case "Map(String,String)":
		c.Data = NewMap[string, string](new(ColStr), new(ColStr))
	case ColumnTypeUUID:
		c.Data = new(ColUUID)
	case ColumnTypeArray.Sub(ColumnTypeUUID):
		c.Data = new(ColUUID).Array()
	case ColumnTypeNullable.Sub(ColumnTypeUUID):
		c.Data = new(ColUUID).Nullable()
	default:
		switch t.Base() {
		case ColumnTypeDateTime:
			v := new(ColDateTime)
			if err := v.Infer(t); err != nil {
				return errors.Wrap(err, "datetime")
			}
			c.Data = v
			c.DataType = t
			return nil
		case ColumnTypeDecimal:
			var prec int
			precStr, _, _ := strings.Cut(string(t.Elem()), ",")
			if precStr != "" {
				var err error
				precStr = strings.TrimSpace(precStr)
				prec, err = strconv.Atoi(precStr)
				if err != nil {
					return errors.Wrap(err, "decimal")
				}
			} else {
				prec = 10
			}
			switch {
			case prec >= 1 && prec < 10:
				c.Data = new(ColDecimal32)
			case prec >= 10 && prec < 19:
				c.Data = new(ColDecimal64)
			case prec >= 19 && prec < 39:
				c.Data = new(ColDecimal128)
			case prec >= 39 && prec < 77:
				c.Data = new(ColDecimal256)
			default:
				return errors.Errorf("decimal precision %d out of range", prec)
			}
			c.DataType = t
			return nil
		case ColumnTypeDecimal32:
			c.Data = new(ColDecimal32)
			c.DataType = t
			return nil
		case ColumnTypeDecimal64:
			c.Data = new(ColDecimal64)
			c.DataType = t
			return nil
		case ColumnTypeDecimal128:
			c.Data = new(ColDecimal128)
			c.DataType = t
			return nil
		case ColumnTypeDecimal256:
			c.Data = new(ColDecimal256)
			c.DataType = t
			return nil
		case ColumnTypeEnum8, ColumnTypeEnum16:
			v := new(ColEnum)
			if err := v.Infer(t); err != nil {
				return errors.Wrap(err, "enum")
			}
			c.Data = v
			c.DataType = t
			return nil
		case ColumnTypeDateTime64:
			v := new(ColDateTime64)
			if err := v.Infer(t); err != nil {
				return errors.Wrap(err, "datetime64")
			}
			c.Data = v
			c.DataType = t
			return nil
		}
		return errors.Errorf("automatic column inference not supported for %q", t)
	}

	c.DataType = t
	return nil
}

var (
	_ Column    = &ColAuto{}
	_ Inferable = &ColAuto{}
)

func (c ColAuto) Type() ColumnType {
	return c.DataType
}

func (c ColAuto) Rows() int {
	return c.Data.Rows()
}

func (c ColAuto) DecodeColumn(r *Reader, rows int) error {
	return c.Data.DecodeColumn(r, rows)
}

func (c ColAuto) Reset() {
	c.Data.Reset()
}

func (c ColAuto) EncodeColumn(b *Buffer) {
	c.Data.EncodeColumn(b)
}

func (c ColAuto) WriteColumn(w *Writer) {
	c.Data.WriteColumn(w)
}
