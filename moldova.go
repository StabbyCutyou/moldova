// Package moldova is a lightweight generator of random data, based on a provided
// template. It supports a number of tokens which will be replaced with random values,
// based on the type and arguments of each token.
package moldova

import (
	"bytes"
	crand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	// I want to keep files that only exist to help provide sources of data or are
	// helpers to Moldova in their own subdirectory, for organization reasons. Go
	// requires that this be it's own package, which means I'd need to reference them
	// with their package name if I wanted to use them, but I'd rather just have them
	// all be "considered" part of the same package/namespace. So, I purposefully use
	// a dot import here to do so, despite that in most cases dot importing is not great
	. "github.com/StabbyCutyou/moldova/data"
)

type cmdOptions map[string]string
type objectCache map[string]interface{}

// TokenWriter is a closure that wraps a call to generate random data, and places
// the result into the provided buffer
type tokenWriter func(*bytes.Buffer, objectCache) error

// Callstack is a list of closures to invoke in order to generate the result of a
// parsed template. Callstack is a FIFO implementation, making it more akin to a queue
// than a stack.
type Callstack struct {
	stack []tokenWriter
	cache objectCache
}

func newCallstack() *Callstack {
	return &Callstack{
		stack: make([]tokenWriter, 0),
	}
}

// Push will place the given tokenWriter function onto the stack. The first function
// placed onto the stack will be the first one called when Write is called
func (c *Callstack) Push(t tokenWriter) {
	c.stack = append(c.stack, t)
}

// Write will take a bytes.Buffer pointer and fill it with the results of calling
// each known function on the Callstack.
func (c *Callstack) Write(result *bytes.Buffer) error {
	c.cache = newObjectCache()
	for _, f := range c.stack {
		if err := f(result, c.cache); err != nil {
			return err
		}
	}
	return nil
}

var defaultOptions = map[string]cmdOptions{
	"guid":    cmdOptions{"ordinal": "-1"},
	"now":     cmdOptions{"ordinal": "-1", "format": "simple"},
	"time":    cmdOptions{"ordinal": "-1", "format": "simple", "min": "0", "max": "1455512165"},
	"int":     cmdOptions{"min": "0", "maxr": "100", "ordinal": "-1"},
	"float":   cmdOptions{"min": "0.0", "maxr": "100.0", "ordinal": "-1"},
	"ascii":   cmdOptions{"length": "2", "case": "down", "ordinal": "-1"},
	"unicode": cmdOptions{"length": "2", "case": "down", "ordinal": "-1"},
	"country": cmdOptions{"ordinal": "-1", "case": "up"},
}

func newObjectCache() objectCache {
	return objectCache{
		"guid":    make([]string, 0),
		"now":     make([]string, 0),
		"time":    make([]string, 0),
		"country": make([]string, 0),
		"unicode": make([]string, 0),
		"ascii":   make([]string, 0),
		"int":     make([]int, 0),
		"float":   make([]float64, 0),
	}
}

// BuildCallstack will parse the template, and return a callstack of closures to
// invoke in order, which will produce static/random values that can be turned into
// a string
func BuildCallstack(inputTemplate string) (*Callstack, error) {
	stack := newCallstack()
	wordBuffer := &bytes.Buffer{}
	foundWord := false
	for _, c := range inputTemplate {
		if !foundWord && c == '{' {
			// We're starting a word to parse
			foundWord = true
			// Dump the current buffer into a closure
			// Assigning to 'cb', ClosureBuster, will get around this issue
			// THANKS .NET PRIOR TO 4.0 FOR TEACHING ME ABOUT ACCESS TO A MODIFIED CLOSURE!
			cb := wordBuffer.String()
			wordBuffer.Reset()
			f := func(result *bytes.Buffer, cache objectCache) error {
				result.WriteString(cb)
				return nil
			}
			stack.Push(f)
		} else if foundWord && c == '}' {
			// We're closing a word, so eval it and get the data to put in the string
			foundWord = false
			parts := strings.SplitN(wordBuffer.String(), ":", 2)
			rawOpts := ""
			if len(parts) > 1 {
				rawOpts = parts[1]
			}
			opts, err := optionsToMap(parts[0], rawOpts)
			if err != nil {
				return nil, err
			}
			// Build the closure that will invoke resolveWord
			f := func(result *bytes.Buffer, cache objectCache) error {
				val := ""
				if val, err = resolveWord(cache, parts[0], opts); err != nil {
					return err
				}
				result.WriteString(val)
				return nil
			}
			stack.Push(f)
			wordBuffer.Reset()
		} else {
			// Straight pass through
			wordBuffer.WriteRune(c)
		}
	}

	return stack, nil
}

// This function was borrowed with permission from the following location
// https://github.com/dgryski/trifles/blob/master/uuid/uuid.go
// All credit / lawsuits can be forwarded to Damian Gryski and Russ Cox
func uuidv4() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(crand.Reader, b)
	if err != nil {
		// probably "shouldn't happen"
		log.Fatal(err)
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func optionsToMap(name string, options string) (map[string]string, error) {
	parts := strings.Split(options, "|")
	m := make(map[string]string)
	defaults := defaultOptions[name]
	for k, v := range defaults {
		m[k] = v
	}
	// If there were no options specified, just use defaults
	if len(options) == 0 {
		return m, nil
	}
	for _, p := range parts {
		// Some options, like format, can have : in them. Only split the first :, which
		// should have the arg name, ad a value with an arbitrary number of : inside of it
		opt := strings.SplitN(p, ":", 2)
		m[opt[0]] = opt[1]
	}
	return m, nil
}

func resolveWord(oc objectCache, word string, opts cmdOptions) (string, error) {
	// If there were options provided, convert them to a lookup map prior to invoking
	// a randomizer.
	switch word {
	case "guid":
		return guid(oc, opts)
	case "int":
		return integer(oc, opts)
	case "now":
		return now(oc, opts)
	case "time":
		return datetime(oc, opts)
	case "float":
		return float(oc, opts)
	case "unicode":
		return unicode(oc, opts)
	case "country":
		return country(oc, opts)
	}
	// TODO make this an error
	return "", nil
}

// TODO All the below functions need way better commenting and parameter annotations
// It's described in the readme, but I should probably make these public and then
// give them proper comments, so that GoDoc can also document them

func integer(oc objectCache, opts cmdOptions) (string, error) {
	lb := opts["min"]
	ub := opts["max"]
	min, err := strconv.Atoi(lb)
	if err != nil {
		return "", err
	}
	max, err := strconv.Atoi(ub)
	if err != nil {
		return "", err
	}
	o := opts["ordinal"]
	ord, err := strconv.Atoi(o)
	if err != nil {
		return "", err
	}

	if ord >= 0 {
		c := oc["int"]
		cache := c.([]int)
		if len(cache)-1 < ord {
			return "", fmt.Errorf("Ordinal %d has not yet been encountered for integers. Please check your input string", ord)
		}
		i := cache[ord]
		return strconv.Itoa(i), nil
	}

	if min > max {
		return "", errors.New("You cannot generate a random number whose lower bound is greater than it's upper bound. Please check your input string")
	}

	// Incase we need to tell the function to invert the case
	negateResult := false
	// get the difference between them
	diff := max - min
	// Since this supports negatives, need to handle some special corner cases?
	if max < 0 && min <= 0 {
		// if the range is entirely negative
		negateResult = true
		// Swap them, so they are still the same relative distance from eachother, but positive - invert the result
		oldLower := min
		min = -max
		max = -oldLower
	}
	// neg to pos ranges currently not supported
	// else both are positive
	// get a number from 0 to diff
	n := rand.Intn(diff)
	// add lowerbound to it - now it's between lower and upper
	n += min
	if negateResult {
		n = -n
	}

	// store it in the cache
	ca := oc["int"]
	cache := ca.([]int)
	oc["int"] = append(cache, n)

	return strconv.Itoa(n), nil
}

func float(oc objectCache, opts cmdOptions) (string, error) {
	lb := opts["min"]
	ub := opts["max"]
	min, err := strconv.ParseFloat(lb, 64)
	if err != nil {
		return "", err
	}
	max, err := strconv.ParseFloat(ub, 64)
	if err != nil {
		return "", err
	}
	o := opts["ordinal"]
	ord, err := strconv.Atoi(o)
	if err != nil {
		return "", err
	}

	if ord >= 0 {
		c := oc["float"]
		cache := c.([]float64)
		if len(cache)-1 < ord {
			return "", fmt.Errorf("Ordinal %d has not yet been encountered for integers. Please check your input string", ord)
		}
		n := cache[ord]
		return fmt.Sprintf("%f", n), nil
	}

	if min > max {
		return "", errors.New("You cannot generate a random number whose lower bound is greater than it's upper bound. Please check your input string")
	}

	// Incase we need to tell the function to invert the case
	negateResult := false
	// get the difference between them
	diff := max - min
	// Since this supports negatives, need to handle some special corner cases?
	if min < 0.0 && max <= 0.0 {
		// if the range is entirely negative
		negateResult = true
		// Swap them, so they are still the same relative distance from eachother, but positive - invert the result
		oldLower := min
		min = -max
		max = -oldLower
	}
	// neg to pos ranges currently not supported
	// else both are positive
	// get a number from 0 to diff
	n := (rand.Float64() * diff) + min

	if negateResult {
		n = -n
	}

	// store it in the cache
	ca := oc["float"]
	cache := ca.([]float64)
	oc["float"] = append(cache, n)

	return fmt.Sprintf("%f", n), nil
}

func country(oc objectCache, opts cmdOptions) (string, error) {
	cCase := opts["case"]
	o := opts["ordinal"]
	ord, err := strconv.Atoi(o)
	if err != nil {
		return "", err
	}

	if ord >= 0 {
		c, _ := oc["country"]
		cache := c.([]string)
		if len(cache)-1 < ord {
			return "", fmt.Errorf("Ordinal %d has not yet been encountered for countries. Please check your input string", ord)
		}
		country := cache[ord]
		// Countries go into the cache upper case, only check for lowering it
		if c == "down" {
			return strings.ToLower(country), nil
		}
		return country, nil
	}
	// Generate a new one
	n := rand.Intn(len(CountryCodes))
	country := CountryCodes[n]
	// store it in the cache
	ca := oc["country"]
	cache := ca.([]string)
	oc["country"] = append(cache, country)

	if cCase == "down" {
		return strings.ToLower(country), nil
	}

	return country, nil

}

func unicode(oc objectCache, opts cmdOptions) (string, error) {
	cCase := opts["case"]
	n := opts["length"]
	num, err := strconv.Atoi(n)
	if err != nil {
		return "", err
	} else if num <= 0 {
		return "", errors.New("You have specified a number of characters to generate which is not a number greater than zero. Please check your input string")
	}
	o := opts["ordinal"]
	ord, err := strconv.Atoi(o)
	if err != nil {
		return "", err
	}

	if ord >= 0 {
		c, _ := oc["unicode"]
		cache := c.([]string)
		if len(cache)-1 < ord {
			return "", fmt.Errorf("Ordinal %d has not yet been encountered for unicode strings. Please check your input string", ord)
		}
		str := cache[ord]
		// Countries go into the cache upper case, only check for lowering it
		if c == "up" {
			return strings.ToUpper(str), nil
		}
		return str, nil
	}

	result := generateRandomString(num)
	// store it in the cache
	ca := oc["unicode"]
	cache := ca.([]string)
	oc["unicode"] = append(cache, result)
	if cCase == "up" {
		return strings.ToUpper(string(result)), nil
	}
	return string(result), nil
}

func generateRandomString(length int) string {
	rarr := make([]rune, length)
	for i := 0; i < length; i++ {
		// First, pick which range this character comes from
		o := rand.Intn(len(PrintableRanges))
		r := PrintableRanges[o]

		minCharCode := r[0]
		maxCharCode := r[1]

		// Get the delata between max and min
		diff := maxCharCode - minCharCode
		// Get a random value within the range specified
		num := rand.Intn(diff) + minCharCode
		// Turn it into a rune, set it on the result object
		rarr[i] = rune(num)
	}
	return string(rarr)
}

func now(oc objectCache, opts cmdOptions) (string, error) {
	f := opts["format"]
	o := opts["ordinal"]
	ord, err := strconv.Atoi(o)
	if err != nil {
		return "", err
	}
	if ord >= 0 {
		c, _ := oc["now"]
		cache := c.([]string)
		if len(cache)-1 < ord {
			return "", fmt.Errorf("Ordinal %d has not yet been encountered for time-now. Please check your input string", ord)
		}
		return cache[ord], nil
	}
	now := time.Now()
	ts := formatTime(&now, f)

	// store it in the cache
	c, _ := oc["now"]
	cache := c.([]string)
	oc["now"] = append(cache, ts)

	return ts, nil
}

func datetime(oc objectCache, opts cmdOptions) (string, error) {
	lb := opts["min"]
	ub := opts["max"]
	min, err := strconv.Atoi(lb)
	if err != nil {
		return "", err
	}
	max, err := strconv.Atoi(ub)
	if err != nil {
		return "", err
	}
	if min > max {
		return "", errors.New("You cannot generate a random time whose lower bound is greater than it's upper bound. Please check your input string")
	}

	f := opts["format"]

	o := opts["ordinal"]
	ord, err := strconv.Atoi(o)
	if err != nil {
		return "", err
	}
	if ord >= 0 {
		c, _ := oc["time"]
		cache := c.([]string)
		if len(cache)-1 < ord {
			return "", fmt.Errorf("Ordinal %d has not yet been encountered for time-now. Please check your input string", ord)
		}
		return cache[ord], nil
	}
	ut := rand.Int63n(int64(max)) - int64(min)
	t := time.Unix(ut, 0)
	ts := formatTime(&t, f)
	// store it in the cache
	c, _ := oc["time"]
	cache := c.([]string)
	oc["time"] = append(cache, ts)

	return ts, nil
}

func formatTime(t *time.Time, format string) string {
	if f, ok := TimeFormats[format]; ok {
		return t.Format(f)
	}
	return t.Format(format)
}

func guid(oc objectCache, opts cmdOptions) (string, error) {
	o := opts["ordinal"]
	ord, err := strconv.Atoi(o)
	if err != nil {
		return "", err
	}
	if ord >= 0 {
		c, _ := oc["guid"]
		cache := c.([]string)
		if len(cache)-1 < ord {
			return "", fmt.Errorf("Ordinal %d has not yet been encountered for guids. Please check your input string", ord)
		}
		return cache[ord], nil
	}

	guid := uuidv4()
	// store it in the cache
	c, _ := oc["guid"]
	cache := c.([]string)
	oc["guid"] = append(cache, guid)

	return guid, nil

}
