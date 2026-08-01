package bytesutil

import (
	"encoding/json"
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

// Bytes processes JSON string inputs from C bridge and returns JSON string result
func Bytes(valJson string, optsJson string) string {
	val := parseJSONVal(valJson)
	if val == nil {
		return returnJSON(nil)
	}

	switch v := val.(type) {
	case string:
		return returnJSON(parseImpl(v))
	case float64:
		opts := parseJSONOpts(optsJson)
		return returnJSON(formatImpl(v, opts))
	default:
		return returnJSON(nil)
	}
}

// Parse processes JSON string input from C bridge and returns JSON string result
func Parse(valJson string) string {
	val := parseJSONVal(valJson)
	return returnJSON(parseImpl(val))
}

// Format processes JSON string inputs from C bridge and returns JSON string result
func Format(valJson string, optsJson string) string {
	val := parseJSONVal(valJson)
	num, ok := val.(float64)
	if !ok {
		return returnJSON(nil)
	}
	opts := parseJSONOpts(optsJson)
	return returnJSON(formatImpl(num, opts))
}

func parseImpl(val interface{}) interface{} {
	switch v := val.(type) {
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil
		}
		return v
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

		return math.Floor(unitMap[unit] * floatValue)
	default:
		return nil
	}
}

func formatImpl(value float64, options *Options) interface{} {
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

func parseJSONVal(valJson string) interface{} {
	if valJson == "" || valJson == "undefined" {
		return nil
	}
	var val interface{}
	if err := json.Unmarshal([]byte(valJson), &val); err != nil {
		return nil
	}
	return val
}

func parseJSONOpts(optsJson string) *Options {
	if optsJson == "" || optsJson == "undefined" || optsJson == "null" {
		return nil
	}
	var opts Options
	if err := json.Unmarshal([]byte(optsJson), &opts); err != nil {
		return nil
	}
	return &opts
}

func returnJSON(val interface{}) string {
	if val == nil {
		return "null"
	}
	b, err := json.Marshal(val)
	if err != nil {
		return "null"
	}
	return string(b)
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
