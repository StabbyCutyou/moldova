package moldova

import (
	"bytes"
	"strconv"
	"testing"
)

// TODO Test each random function individually, under a number of inputs to make supported
// all the options behave as expected.

func TestBuildSQL(t *testing.T) {
	template := "INSERT INTO floof VALUES ('{guid}','{guid:ordinal:0}','{country}',{int:min:-2000|max:0},{int:min:100|max:1000},{float:min:-1000.0|max:-540.0},{int:min:1|max:40},'{now}','{now:ordinal:0}','{unicode:length:2|case:up}',NULL,-3)"
	cs, err := BuildCallstack(template)
	if err != nil {
		t.Error(err)
	}
	result := &bytes.Buffer{}
	err = cs.Write(result)
	if err != nil {
		t.Error(err)
	}
}

func TestCountries(t *testing.T) {
	template := "INSERT INTO `floop` VALUES ('{country}','{country:case:up|ordinal:0}','{country}','{country:case:down|ordinal:1}')"
	cs, err := BuildCallstack(template)
	if err != nil {
		t.Error(err)
	}
	result := &bytes.Buffer{}
	err = cs.Write(result)
	if err != nil {
		t.Error(err)
	}
}

func TestInteger(t *testing.T) {
	template := "{int:min:5|max:6}"
	cs, err := BuildCallstack(template)
	if err != nil {
		t.Error(err)
	}
	result := &bytes.Buffer{}
	err = cs.Write(result)
	if err != nil {
		t.Error(err)
	}

	c, err := strconv.Atoi(result.String())
	if err != nil {
		t.Error(err)
	}
	if c < 5 || c > 6 {
		t.Error("Integer out of range")
	}
}

func TestNowOrdinal(t *testing.T) {
	template := "{now:ordinal:1}"
	cs, err := BuildCallstack(template)
	if err != nil {
		t.Error(err)
	}
	result := &bytes.Buffer{}
	err = cs.Write(result)
	if err == nil {
		t.Error("Did not return an error on an invalid {now} ordinal")
	}
}

func TestGuidOrdinal(t *testing.T) {
	template := "{guid:ordinal:1}"
	cs, err := BuildCallstack(template)
	if err != nil {
		t.Error(err)
	}
	result := &bytes.Buffer{}
	err = cs.Write(result)
	if err == nil {
		t.Error("Did not return an error on an invalid {guid} ordinal")
	}
}

func BenchmarkBuildCallstackRuns(b *testing.B) {
	template := "INSERT INTO floof VALUES ('{guid}','{guid:ordinal:0}','{country}',{int:min:-2000|max:0},{int:min:100|max:1000},{float:min:-1000.0|max:-540.0},{int:min:1|max:40},'{now}','{now:ordinal:0}','{unicode:length:2|case:up}',NULL,-3)"
	var cs *Callstack
	var err error
	for n := 0; n < b.N; n++ {
		if n == 0 {
			if cs, err = BuildCallstack(template); err != nil {
				b.Error(err)
			}
		}
		result := &bytes.Buffer{}
		cs.Write(result)
	}
}
