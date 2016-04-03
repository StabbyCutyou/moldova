package moldova

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type TestComparator func(string) error

type TestCase struct {
	Template     string
	Comparator   TestComparator
	ParseFailure bool
	WriteFailure bool
}

var GUIDCases = []TestCase{
	{
		Template: "{guid}",
		Comparator: func(s string) error {
			p := strings.Split(s, "-")
			if len(p) == 5 &&
				len(p[0]) == 8 &&
				len(p[1]) == len(p[2]) && len(p[2]) == len(p[3]) && len(p[3]) == 4 &&
				len(p[4]) == 12 {
				return nil
			}
			return errors.New("Guid not in correct format: " + s)
		},
	},
	{
		Template: "{guid}@{guid:ordinal:0}",
		Comparator: func(s string) error {
			p := strings.Split(s, "@")
			if p[0] == p[1] {
				return nil
			}
			return errors.New("Guid at position 1 not equal to guid at position 0: " + p[0] + " " + p[1])
		},
	},
	{
		Template:     "{guid}@{guid:ordinal:1}",
		WriteFailure: true,
	},
}

var NowCases = []TestCase{
	{
		// There is no proper deterministic way to test what the value of now is, without
		// something like rubys timecop (but the go-equivalent is not viable) or relying
		// on luck, which will run out if tests are run at just the wrong moment.
		// Therefore, for the basic test, i'm just asserting nothing went wrong for now.
		Template: "{now}",
		Comparator: func(s string) error {
			if len(s) > 0 {
				return nil
			}
			return errors.New("Now not in correct format: " + s)
		},
	},
	{
		Template: "{now}@{now:ordinal:0}",
		Comparator: func(s string) error {
			p := strings.Split(s, "@")
			if p[0] == p[1] {
				return nil
			}
			return errors.New("Now at position 1 not equal to now at position 0: " + p[0] + " " + p[1])
		},
	},
	{
		Template:     "{now}@{now:ordinal:1}",
		WriteFailure: true,
	},
}

var TimeCases = []TestCase{
	{
		Template: "{time:min:1|max:1|format:simple|zone:EST}",
		Comparator: func(s string) error {
			if s == "1969-12-31 19:00:01" {
				return nil
			}
			return errors.New("Time value was not the expected value")
		},
	},
	{
		Template: "{time:min:1|max:1|format:simpletz|zone:EST}",
		Comparator: func(s string) error {
			if s == "1969-12-31 19:00:01 -0500" {
				return nil
			}
			return errors.New("Time value was not the expected value")
		},
	},
	{
		Template: "{time:min:1|max:1|format:2006//01//02@@15_04_05|zone:EST}",
		Comparator: func(s string) error {
			if s == "1969//12//31@@19_00_01" {
				return nil
			}
			return errors.New("Time value was not the expected value")
		},
	},
	{
		Template: "{time}@{time:ordinal:0}",
		Comparator: func(s string) error {
			p := strings.Split(s, "@")
			if p[0] == p[1] {
				return nil
			}
			return errors.New("Time at position 1 not equal to time at position 0: " + p[0] + " " + p[1])
		},
	},
	{
		Template:     "{time}@{time:ordinal:1}",
		WriteFailure: true,
	},
}

var CountryCases = []TestCase{
	{
		Template: "{country}",
		Comparator: func(s string) error {
			// TODO better check here in case we ever support different types of country codes
			if len(s) == 2 {
				return nil
			}
			return errors.New("Invalid country code generated somehow")
		},
	},
	{
		Template: "{country:case:up}",
		Comparator: func(s string) error {
			// Since I can't know which country comes out, i'll invert the result
			// If the ToLowered result is not the same as the original result, we know
			// that the original was successfully output in upper case
			if strings.ToLower(s) != s {
				return nil
			}
			return errors.New("Country was returned in lowercase, but was requested in uppercase")
		},
	},
	{
		Template: "{country:case:down}",
		Comparator: func(s string) error {
			// Since I can't know which country comes out, i'll invert the result
			// If the ToLowered result is not the same as the original result, we know
			// that the original was successfully output in upper case
			if strings.ToUpper(s) != s {
				return nil
			}
			return errors.New("Country was returned in uppercase, but was requested in lowercase")
		},
	},
	{
		Template: "{country}@{country:ordinal:0}",
		Comparator: func(s string) error {
			p := strings.Split(s, "@")
			if p[0] == p[1] {
				return nil
			}
			return errors.New("Country at position 1 not equal to country at position 0: " + p[0] + " " + p[1])
		},
	},
	{
		Template:     "{country}@{country:ordinal:1}",
		WriteFailure: true,
	},
}

// Placeholders
var FloatCases = []TestCase{
	{
		Template: "{float}",
		Comparator: func(s string) error {
			i, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			if i >= 0.0 && i <= 100.0 {
				return nil
			}
			return errors.New("Float out of range for default min/max values")
		},
	},
	{
		Template: "{float:max:5000.0|min:4999.0}",
		Comparator: func(s string) error {
			i, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			if i >= 4999.0 && i <= 5000.0 {
				return nil
			}
			return errors.New("Float out of range for custom min/max values")
		},
	},
	{
		Template: "{float:max:-5000.0|min:-5001.0}",
		Comparator: func(s string) error {
			i, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			if i >= -5001.0 && i <= -5000.0 {
				return nil
			}
			return errors.New("Float out of range for custom min/max values")
		},
	},
	{
		Template: "{float}@{float:ordinal:0}",
		Comparator: func(s string) error {
			p := strings.Split(s, "@")
			if p[0] == p[1] {
				return nil
			}
			return errors.New("Float at position 1 not equal to Float at position 0: " + p[0] + " " + p[1])
		},
	},
	{
		Template:     "{float}@{float:ordinal:1}",
		WriteFailure: true,
	},
}

var IntegerCases = []TestCase{
	{
		Template: "{int}",
		Comparator: func(s string) error {
			i, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			if i >= 0 && i <= 100 {
				return nil
			}
			return errors.New("Int out of range for default min/max values")
		},
	},
	{
		Template: "{int:max:5000|min:4999}",
		Comparator: func(s string) error {
			i, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			if i >= 4999 && i <= 5000 {
				return nil
			}
			return errors.New("Int out of range for custom min/max values")
		},
	},
	{
		Template: "{int:max:-5000|min:-5001}",
		Comparator: func(s string) error {
			i, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			if i >= -5001 && i <= -5000 {
				return nil
			}
			return errors.New("Int out of range for custom min/max values")
		},
	},
	{
		Template: "{int}@{int:ordinal:0}",
		Comparator: func(s string) error {
			p := strings.Split(s, "@")
			if p[0] == p[1] {
				return nil
			}
			return errors.New("Int at position 1 not equal to int at position 0: " + p[0] + " " + p[1])
		},
	},
	{
		Template:     "{int}@{int:ordinal:1}",
		WriteFailure: true,
	},
}

var UnicodeCases = []TestCase{
	{
		Template: "{unicode}",
		Comparator: func(s string) error {
			if len([]rune(s)) == 2 {
				return nil
			}
			return errors.New("Unicode string not the correct length")
		},
	},
	{
		Template: "{unicode:length:10}",
		Comparator: func(s string) error {
			if len([]rune(s)) == 10 {
				return nil
			}
			return errors.New("Unicode string not the correct length")
		},
	},
	{
		Template: "{unicode:length:10|case:up}",
		Comparator: func(s string) error {
			if len([]rune(s)) == 10 || strings.ToLower(s) != strings.ToUpper(s) {
				return nil
			}
			return errors.New("Unicode string not the correct length")
		},
	},
	{
		Template: "{unicode}@{unicode:ordinal:0}",
		Comparator: func(s string) error {
			p := strings.Split(s, "@")
			if p[0] == p[1] {
				return nil
			}
			return errors.New("Unicode at position 1 not equal to Unicode at position 0: " + p[0] + " " + p[1])
		},
	},
	{
		Template:     "{unicode}@{unicode:ordinal:1}",
		WriteFailure: true,
	},
}

var FirstNameCases = []TestCase{
	{
		Template: "{firstname}",
		Comparator: func(s string) error {
			if len(s) > 0 {
				return nil
			}
			return errors.New("First Name string not the correct length")
		},
	},
	{
		Template: "{firstname}@{firstname:ordinal:0}",
		Comparator: func(s string) error {
			p := strings.Split(s, "@")
			if p[0] == p[1] {
				return nil
			}
			return errors.New("First Name at position 1 not equal to Unicode at position 0: " + p[0] + " " + p[1])
		},
	},
	{
		Template:     "{firstname}@{firstname:ordinal:1}",
		WriteFailure: true,
	},
}

var LastNameCases = []TestCase{
	{
		Template: "{lastname}",
		Comparator: func(s string) error {
			if len(s) > 0 {
				return nil
			}
			return errors.New("First Name string not the correct length")
		},
	},
	{
		Template: "{lastname}@{lastname:ordinal:0}",
		Comparator: func(s string) error {
			p := strings.Split(s, "@")
			if p[0] == p[1] {
				return nil
			}
			return errors.New("First Name at position 1 not equal to Unicode at position 0: " + p[0] + " " + p[1])
		},
	},
	{
		Template:     "{lastname}@{lastname:ordinal:1}",
		WriteFailure: true,
	},
}

var FullNameCases = []TestCase{
	{
		Template: "{firstname} {lastname}",
		Comparator: func(s string) error {
			if len(s) > 0 {
				return nil
			}
			return errors.New("First / Last Name string not the correct length")
		},
	},
}

var AllCases = [][]TestCase{
	GUIDCases,
	NowCases,
	TimeCases,
	CountryCases,
	FloatCases,
	IntegerCases,
	UnicodeCases,
	FirstNameCases,
	LastNameCases,
	FullNameCases,
}

// TODO Test each random function individually, under a number of inputs to make supported
// all the options behave as expected.

func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())
	os.Exit(m.Run())
}

func TestAllCases(t *testing.T) {
	// TODO The library should be threadsafe, I should go wide here to run all specs
	// in parallel, like the natural tests would be. Channel + waitgroup to collect
	// and report on errors once they all finish
	for _, cs := range AllCases {
		for _, c := range cs {
			cs, err := BuildCallstack(c.Template)
			// If we get an error and weren't expecting it
			// Or, if we didn't get one but were expecting it
			if err != nil && !c.ParseFailure {
				t.Error(err)
			} else if err == nil && c.ParseFailure {
				t.Error("Expected to encounter Parse Failure, but did not for Test Case ", c.Template)
			}

			result := &bytes.Buffer{}
			err = cs.Write(result)

			// If we get an error and weren't expecting it
			// Or, if we didn't get one but were expecting it
			if err != nil && !c.WriteFailure {
				t.Error(err)
			} else if err == nil && c.ParseFailure {
				t.Error("Expected to encounter Write Failure, but did not for Test Case ", c.Template)
			}

			if c.Comparator != nil {
				if err := c.Comparator(result.String()); err != nil {
					t.Error(err)
				}
			}
		}
	}
}

func TestGeneratedStringLength(t *testing.T) {
	template := "Hey I'm {int:min:1|max:9} years old"
	sampleresult := "Hey I'm 1 years old"
	cs, err := BuildCallstack(template)
	if err != nil {
		t.Error(err)
	}
	result := &bytes.Buffer{}
	err = cs.Write(result)
	if err != nil {
		t.Error(err)
	}
	if len([]rune(sampleresult)) != len([]rune(result.String())) {
		fmt.Println(result.String())
		t.Error("Missing parts of the rendered templtate")
	}
}

func BenchmarkGUID(b *testing.B) {
	c := GUIDCases[0]
	var cs *Callstack
	var err error
	if cs, err = BuildCallstack(c.Template); err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result := &bytes.Buffer{}
		err = cs.Write(result)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkNow(b *testing.B) {
	c := NowCases[0]
	var cs *Callstack
	var err error
	if cs, err = BuildCallstack(c.Template); err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result := &bytes.Buffer{}
		err = cs.Write(result)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkTime(b *testing.B) {
	c := TimeCases[0]
	var cs *Callstack
	var err error
	if cs, err = BuildCallstack(c.Template); err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result := &bytes.Buffer{}
		err = cs.Write(result)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkInteger(b *testing.B) {
	c := IntegerCases[0]
	var cs *Callstack
	var err error
	if cs, err = BuildCallstack(c.Template); err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result := &bytes.Buffer{}
		err = cs.Write(result)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkFloat(b *testing.B) {
	c := FloatCases[0]
	var cs *Callstack
	var err error
	if cs, err = BuildCallstack(c.Template); err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result := &bytes.Buffer{}
		err = cs.Write(result)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkCountry(b *testing.B) {
	c := CountryCases[0]
	var cs *Callstack
	var err error
	if cs, err = BuildCallstack(c.Template); err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result := &bytes.Buffer{}
		err = cs.Write(result)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkUnicode(b *testing.B) {
	c := UnicodeCases[0]
	var cs *Callstack
	var err error
	if cs, err = BuildCallstack(c.Template); err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result := &bytes.Buffer{}
		err = cs.Write(result)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkFirstName(b *testing.B) {
	c := FirstNameCases[0]
	var cs *Callstack
	var err error
	if cs, err = BuildCallstack(c.Template); err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result := &bytes.Buffer{}
		err = cs.Write(result)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkLastName(b *testing.B) {
	c := FirstNameCases[0]
	var cs *Callstack
	var err error
	if cs, err = BuildCallstack(c.Template); err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result := &bytes.Buffer{}
		err = cs.Write(result)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkFullName(b *testing.B) {
	c := FullNameCases[0]
	var cs *Callstack
	var err error
	if cs, err = BuildCallstack(c.Template); err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result := &bytes.Buffer{}
		err = cs.Write(result)
		if err != nil {
			b.Error(err)
		}
	}
}
