package bytesutil

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var parseRegExp = regexp.MustCompile(`^((-|\+)?(\d+(?:\.\d+)?)) *(kb|mb|gb|tb|pb)$`)

var unitMap = map[string]float64{
	"b":  1,
	"kb": 1 << 10,
	"mb": 1 << 20,
	"gb": 1 << 30,
	"tb": math.Pow(1024, 4),
	"pb": math.Pow(1024, 5),
}

type Options struct {
	DecimalPlaces      *int    `json:"decimalPlaces"`
	FixedDecimals      *bool   `json:"fixedDecimals"`
	ThousandsSeparator *string `json:"thousandsSeparator"`
	Unit               *string `json:"unit"`
	UnitSeparator      *string `json:"unitSeparator"`
}

// Bytes converts a string (e.g. "1KB") to bytes or a number (e.g. 1024) to a formatted byte string.
// Accepts native Go types (string, float64, int, etc.) and returns interface{} (float64, string, or nil).
func Bytes(val interface{}, opts *Options) interface{} {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case string:
		return Parse(v)
	case float64:
		return Format(v, opts)
	case float32:
		return Format(float64(v), opts)
	case int:
		return Format(float64(v), opts)
	case int64:
		return Format(float64(v), opts)
	case int32:
		return Format(float64(v), opts)
	case int8:
		return Format(float64(v), opts)
	case int16:
		return Format(float64(v), opts)
	case uint:
		return Format(float64(v), opts)
	case uint64:
		return Format(float64(v), opts)
	case uint32:
		return Format(float64(v), opts)
	case uint8:
		return Format(float64(v), opts)
	case uint16:
		return Format(float64(v), opts)
	default:
		return nil
	}
}

// Parse parses a string or number into integer byte count.
// Accepts native Go types (string, float64, int, etc.) and returns interface{} (int64, float64, or nil).
func Parse(val interface{}) interface{} {
	switch v := val.(type) {
	case float64:
		return toBestNumericType(v)
	case float32:
		return toBestNumericType(float64(v))
	case int:
		return int64(v)
	case int64:
		return v
	case int32:
		return int64(v)
	case uint:
		return int64(v)
	case uint64:
		return toBestNumericType(float64(v))
	case string:
		vTrim := strings.TrimSpace(v)
		matches := parseRegExp.FindStringSubmatch(strings.ToLower(vTrim))

		var floatValue float64
		unit := "b"

		if matches == nil {
			floatValue = parseJSInt10(vTrim)
			if math.IsNaN(floatValue) {
				return nil
			}
			unit = "b"
		} else {
			var err error
			floatValue, err = strconv.ParseFloat(matches[1], 64)
			if err != nil {
				return nil
			}
			unit = matches[4]
		}

		if math.IsNaN(floatValue) || math.IsInf(floatValue, 0) {
			return nil
		}

		return toBestNumericType(math.Floor(unitMap[unit] * floatValue))
	default:
		return nil
	}
}

func toBestNumericType(f float64) interface{} {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return nil
	}
	if f == math.Trunc(f) && f >= float64(math.MinInt64) && f <= float64(math.MaxInt64) {
		return int64(f)
	}
	return f
}

// Format converts a float64 byte count into a human-readable string (e.g. "1KB").
// Returns interface{} (string or nil if value is invalid).
func Format(value float64, options *Options) interface{} {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return nil
	}

	mag := math.Abs(value)
	thousandsSeparator := ""
	unitSeparator := ""
	decimalPlaces := 2
	fixedDecimals := false
	unit := ""

	if options != nil {
		if options.ThousandsSeparator != nil {
			thousandsSeparator = *options.ThousandsSeparator
		}
		if options.UnitSeparator != nil {
			unitSeparator = *options.UnitSeparator
		}
		if options.DecimalPlaces != nil {
			decimalPlaces = *options.DecimalPlaces
		}
		if options.FixedDecimals != nil {
			fixedDecimals = *options.FixedDecimals
		}
		if options.Unit != nil {
			unit = *options.Unit
		}
	}

	if unit == "" || unitMap[strings.ToLower(unit)] == 0 {
		switch {
		case mag >= unitMap["pb"]:
			unit = "PB"
		case mag >= unitMap["tb"]:
			unit = "TB"
		case mag >= unitMap["gb"]:
			unit = "GB"
		case mag >= unitMap["mb"]:
			unit = "MB"
		case mag >= unitMap["kb"]:
			unit = "KB"
		default:
			unit = "B"
		}
	}

	val := value / unitMap[strings.ToLower(unit)]
	str := fmt.Sprintf("%.*f", decimalPlaces, val)

	if !fixedDecimals {
		str = trimDecimals(str)
	}

	if thousandsSeparator != "" {
		parts := strings.Split(str, ".")
		parts[0] = insertThousands(parts[0], thousandsSeparator)
		str = strings.Join(parts, ".")
	}

	res := str + unitSeparator + unit
	return res
}

func trimDecimals(s string) string {
	if !strings.Contains(s, ".") {
		return s
	}
	s = strings.TrimRight(s, "0")
	return strings.TrimSuffix(s, ".")
}

func insertThousands(s, sep string) string {
	negative := false
	if strings.HasPrefix(s, "-") {
		negative = true
		s = s[1:]
	}

	n := len(s)
	if n <= 3 {
		if negative {
			return "-" + s
		}
		return s
	}

	var out []byte
	start := n % 3
	if start == 0 {
		start = 3
	}

	out = append(out, s[:start]...)
	for i := start; i < n; i += 3 {
		out = append(out, sep...)
		out = append(out, s[i:i+3]...)
	}

	if negative {
		return "-" + string(out)
	}
	return string(out)
}

func parseJSInt10(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return math.NaN()
	}
	sign := 1.0
	if strings.HasPrefix(s, "-") {
		sign = -1.0
		s = s[1:]
	} else if strings.HasPrefix(s, "+") {
		s = s[1:]
	}
	digits := ""
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			digits += string(ch)
		} else {
			break
		}
	}
	if digits == "" {
		return math.NaN()
	}
	val, err := strconv.ParseFloat(digits, 64)
	if err != nil {
		return math.NaN()
	}
	return sign * val
}
