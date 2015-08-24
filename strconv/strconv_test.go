package strconv

import "fmt"
import "math"
import "strings"
import "strconv"
import "testing"
import "github.com/stretchr/testify/assert"

var (
	nan  = math.NaN()
	pinf = math.Inf(+1)
	ninf = math.Inf(-1)
)

const (
	hiThres = 1.000000000000001
	loThres = 0.999999999999999
)

func atof(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func signStr(f float64) string {
	if f >= 0 {
		return ""
	}
	return "-"
}

func stripDot(s string) string {
	if s == "." {
		return ""
	}
	return s
}

func addIEC(s string) string {
	if s == "" {
		return ""
	}
	return s + "i"
}

func normToAlt(s string) string {
	for _, ch := range s {
		if alt, ok := mapNormToAlt[ch]; ok {
			return string(alt)
		}
	}
	return s
}

func altToNorm(s string) string {
	for _, ch := range s {
		if norm, ok := mapAltToNorm[ch]; ok {
			return string(norm)
		}
	}
	return s
}

func split(s string) (string, string) {
	i := strings.IndexAny(s, parsePrefixes)
	if i < 0 {
		return s, ""
	}
	return s[:i], s[i:]
}

func TestPrefixExactSI(t *testing.T) {
	for _, sign := range []float64{-1, +1} {
		for i, f := range scaleSI {
			str := FormatPrefix(sign*f, SI, -1)
			flt, err := ParsePrefix(str, SI)

			pre := normToAlt(stripDot(string(prefixes[i])))
			assert.Equal(t, signStr(sign)+"1"+pre, str)
			assert.Equal(t, sign*f, flt)
			assert.Equal(t, nil, err)
		}
	}
}

func testPrefixExact(t *testing.T, minExp, mode int, wrap func(string) string) {
	var pres = prefixes[minExp+len(divPrefixes):]
	var str, pre string
	var flt float64
	var err error

	for _, sign := range []float64{-1, +1} {
		var f0 = sign * scaleIEC[minExp+len(divPrefixes)]

		str = FormatPrefix(f0/2, mode, -1)
		flt, err = ParsePrefix(str, mode)
		pre = wrap(string(pres[0]))
		assert.Equal(t, signStr(sign)+"0.5"+pre, str)
		assert.Equal(t, f0/2, flt)
		assert.Equal(t, nil, err)

		for i := 0; i < 10*len(pres); i++ {
			str = FormatPrefix(f0, mode, -1)
			flt, err = ParsePrefix(str, mode)

			pre := wrap(string(pres[(i / 10)]))
			ord := fmt.Sprintf("%d", 1<<uint(i%10))
			assert.Equal(t, signStr(sign)+ord+pre, str)
			assert.Equal(t, f0, flt)
			assert.Equal(t, nil, err)
			f0 *= 2
		}

		str = FormatPrefix(f0, mode, -1)
		flt, err = ParsePrefix(str, mode)
		pre = wrap(string(pres[len(pres)-1]))
		assert.Equal(t, signStr(sign)+"1024"+pre, str)
		assert.Equal(t, f0, flt)
		assert.Equal(t, nil, err)
	}
}

func TestPrefixExactBase1024(t *testing.T) {
	wrap := func(s string) string { return stripDot(s) }
	testPrefixExact(t, minExp, Base1024, wrap)
}

func TestPrefixExactIEC(t *testing.T) {
	wrap := func(s string) string { return addIEC(stripDot(s)) }
	testPrefixExact(t, 0, IEC, wrap)
}

func testPrefixBoundary(t *testing.T, scales []float64, prefixes string, mode int, wrap func(string) string) {
	var str, str1, str2, pre string
	var flt, fabs, fnum float64
	var err error

	base := 1024.0
	if mode == SI {
		base = 1000.0
	}

	for _, sign := range []float64{-1, +1} {
		for i, f := range scales {
			// Round towards zero.
			str = FormatPrefix(math.Nextafter(sign*f, sign*ninf), mode, -1)
			flt, err = ParsePrefix(str, mode)
			fabs = math.Abs(flt)

			pre = string(prefixes[0])
			if i > 0 {
				pre = string(prefixes[i-1])
			}
			pre = wrap(pre)
			str1, str2 = split(str)
			fnum = math.Abs(atof(str1))
			if i == 0 {
				assert.True(t, 1.0*loThres <= fnum && fnum <= 1.0)
			} else {
				assert.True(t, base*loThres <= fnum && fnum <= base)
			}
			assert.Equal(t, pre, str2)
			assert.True(t, f*loThres <= fabs && fabs <= f)
			assert.Equal(t, math.Signbit(flt), math.Signbit(sign))
			assert.Equal(t, nil, err)

			// Round away from zero.
			str = FormatPrefix(math.Nextafter(sign*f, sign*pinf), mode, -1)
			flt, err = ParsePrefix(str, mode)
			fabs = math.Abs(flt)

			pre = wrap(string(prefixes[i]))
			str1, str2 = split(str)
			fnum = math.Abs(atof(str1))
			assert.True(t, 1.0 <= fnum && fnum <= 1.0*hiThres)
			assert.Equal(t, pre, str2)
			assert.True(t, f <= fabs && fabs <= f*hiThres)
			assert.Equal(t, math.Signbit(flt), math.Signbit(sign))
			assert.Equal(t, nil, err)
		}
	}
}

func TestPrefixBoundarySI(t *testing.T) {
	wrap := func(s string) string { return normToAlt(stripDot(s)) }
	testPrefixBoundary(t, scaleSI, prefixes, SI, wrap)
}

func TestPrefixBoundaryBase1024(t *testing.T) {
	wrap := func(s string) string { return altToNorm(stripDot(s)) }
	testPrefixBoundary(t, scaleIEC, prefixes, Base1024, wrap)
}

func TestPrefixBoundaryIEC(t *testing.T) {
	idx := len(prefixes) / 2
	wrap := func(s string) string { return addIEC(altToNorm(stripDot(s))) }
	testPrefixBoundary(t, scaleIEC[idx:], prefixes[idx:], IEC, wrap)
}

func TestPrefixFailParse(t *testing.T) {
	for _, x := range []struct {
		str  string
		mode int
		ok   bool
		flt  float64
	}{
		{"", SI, false, 0},
		{"NaN1M", SI, false, 0},
		{"1", IEC, true, Unit},
		{"1 ", IEC, false, 0},
		{"1M", IEC, false, 0},
		{"1Mi", SI, false, 0},
		{"+1M", Base1024, true, +Mebi},
		{"-1Mi", Base1024, true, -Mebi},
		{"+1Mi", Base1024, true, +Mebi},
		{"1E-3", SI, false, 0},
		{"1e-3", SI, false, 0},
		{"1ki", SI, false, 0},
		{"1ki", IEC, false, 0},
		{"1ki", Base1024, true, Kibi},
		{"+1ki", Base1024, true, Kibi},
		{"1μi", SI, false, 0},
		{"1μi", IEC, false, 0},
		{"1μi", Base1024, false, 0},
		{"1k", SI, true, Kilo},
		{"1k", IEC, false, 0},
		{"1k", Base1024, true, Kibi},
		{"1μ", SI, true, Micro},
		{"1μ ", SI, false, 0},
		{" 1μ", SI, false, 0},
		{"+1μ", SI, true, Micro},
		{"1μ", IEC, false, 0},
		{"1μ", Base1024, true, 1.0 / Mebi},
		{"+1μ", Base1024, true, 1.0 / Mebi},
		{"1mi", IEC, false, 0},
		{"0.000001", SI, true, Micro},
		{"1000000u", SI, true, Unit},
		{"1048576", Base1024, true, Mebi},
		{"1048576Ki", IEC, true, Gibi},
		{"nAn", SI, true, nan},
		{"+nan", Base1024, false, 0},
		{"-NAN", IEC, false, 0},
		{"INF", SI, true, pinf},
		{"+iNf", Base1024, true, pinf},
		{"-inF", IEC, true, ninf},
		{"", AutoParse, false, 0},
		{"123", AutoParse, true, 123},
		{"123Ki", AutoParse, true, 123 * Kibi},
		{"123k", AutoParse, true, 123 * Kilo},
		{"123K", AutoParse, true, 123 * Kilo},
		{"3Mi", AutoParse, true, 3 * Mebi},
		{"3M", AutoParse, true, 3 * Mega},
		{"3E-3", AutoParse, true, 3E-3},
		{"2E2", AutoParse, true, 2E2},
	} {
		flt, err := ParsePrefix(x.str, x.mode)
		if x.ok {
			assert.Nil(t, err)
			if !math.IsNaN(x.flt) || !math.IsNaN(flt) {
				assert.Equal(t, x.flt, flt)
			}
		} else {
			assert.NotNil(t, err)
		}
	}
}

func TestPrefix(t *testing.T) {
	// Test for zero, NaN, -Inf, and +Inf.
	for _, mode := range []int{SI, Base1024, IEC} {
		for _, prec := range []int{-1, 0, +1} {
			for _, f := range []float64{-0.0, +0.0, nan, ninf, pinf} {
				str := FormatPrefix(f, mode, prec)
				flt, err := ParsePrefix(str, mode)

				assert.Equal(t, str, strconv.FormatFloat(f, 'f', prec, 64))
				if !math.IsNaN(f) || !math.IsNaN(flt) {
					assert.Equal(t, f, flt)
				}
				assert.Equal(t, nil, err)
			}
		}
	}

	// Test for a huge range of values.
	for _, mode := range []int{SI, Base1024, IEC} {
		for _, prec := range []int{-1, 0, +1, +2} {
			for i := -100; i <= +100; i++ {
				f := 1.234567890123456 * math.Pow(10, float64(i))
				str := FormatPrefix(f, mode, prec)
				flt, err := ParsePrefix(str, mode)
				str1, _ := split(str)
				fnum := math.Abs(atof(str1))

				// Ensure that we maintain decent precision.
				if prec < 0 {
					assert.True(t, f*loThres <= flt && flt <= f*hiThres)
				} else if flt != 0 {
					assert.True(t, f*0.5 <= flt && flt <= f*2)
				}

				// Ensure that we choose the best scale if possible.
				var base, minScale, maxScale float64
				switch mode {
				case SI:
					base = 1000.0
					minScale, maxScale = scaleSI[0], scaleSI[len(scaleSI)-1]
				case Base1024:
					base = 1024.0
					minScale, maxScale = scaleIEC[0], scaleIEC[len(scaleIEC)-1]
				case IEC:
					base = 1024.0
					minScale, maxScale = 1, scaleIEC[len(scaleIEC)-1]
				}
				if minScale <= f && f <= maxScale {
					assert.True(t, 1.0 <= fnum && fnum <= base)
				}

				assert.Equal(t, nil, err)
			}
		}
	}
}
