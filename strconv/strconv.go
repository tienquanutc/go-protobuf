package strconv

import (
	"strconv"
	"unicode/utf8"
)

func AppendCardinal[T Ordered](b []byte, value T) []byte {
	units := []unit{
		{1, ""},
		{1e-3, " thousand"},
		{1e-6, " million"},
		{1e-9, " billion"},
		{1e-12, " trillion"},
	}
	return scale(b, value, units)
}

func AppendInt[T Signed](b []byte, value T, base int) []byte {
	return strconv.AppendInt(b, int64(value), base)
}

func AppendQuote[T String](b []byte, value T) []byte {
	return strconv.AppendQuote(b, string(value))
}

func AppendSize[T Integer](b []byte, value T) []byte {
	units := []unit{
		{1, " byte"},
		{1e-3, " kilobyte"},
		{1e-6, " megabyte"},
		{1e-9, " gigabyte"},
		{1e-12, " terabyte"},
	}
	return scale(b, value, units)
}

func AppendUint[T Unsigned](b []byte, value T, base int) []byte {
	return strconv.AppendUint(b, uint64(value), base)
}

// mimesniff.spec.whatwg.org#binary-data-byte
func Valid(b []byte) bool {
	for _, c := range b {
		if c <= 0x08 {
			return false
		}
		if c == 0x0B {
			return false
		}
		if c >= 0x0E && c <= 0x1A {
			return false
		}
		if c >= 0x1C && c <= 0x1F {
			return false
		}
	}
	return utf8.Valid(b)
}

func label[T Ordered](b []byte, value T, u unit) []byte {
	var prec int
	if u.factor != 1 {
		prec = 2
	}
	u.factor *= float64(value)
	b = strconv.AppendFloat(b, u.factor, 'f', prec, 64)
	return append(b, u.name...)
}

func scale[T Ordered](b []byte, value T, units []unit) []byte {
	var u unit
	for _, u = range units {
		if u.factor*float64(value) < 1000 {
			break
		}
	}
	return label(b, value, u)
}

type Integer interface {
	Signed | Unsigned
}

type Ordered interface {
	Integer | ~float32 | ~float64
}

type Ratio float64

func NewRatio[T, U Ordered](value T, total U) Ratio {
	var r float64
	if total != 0 {
		r = float64(value) / float64(total)
	}
	return Ratio(r)
}

func (r Ratio) AppendPercent(b []byte) []byte {
	return label(b, r, unit{100, "%"})
}

func (r Ratio) AppendRate(b []byte) []byte {
	units := []unit{
		{1, " byte/s"},
		{1e-3, " kilobyte/s"},
		{1e-6, " megabyte/s"},
		{1e-9, " gigabyte/s"},
		{1e-12, " terabyte/s"},
	}
	return scale(b, r, units)
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type String interface {
	~[]byte | ~[]rune | ~byte | ~rune | ~string
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type unit struct {
	factor float64
	name   string
}
