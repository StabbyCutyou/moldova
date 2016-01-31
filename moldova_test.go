package moldova

import (
	"strconv"
	"testing"
)

// TODO Test each random function individually, under a number of inputs to make supported
// all the options behave as expected.

func TestBuildSQL(t *testing.T) {
	template := "INSERT INTO floof VALUES ('{guid}','{guid:0}','{country}',{int:-2000:0},{int:100:1000},{float:-1000.0:-540.0},{int:1:40},'{now}','{now:0}','{char:2:up}',NULL,-3)"
	_, err := ParseTemplate(template)
	if err != nil {
		t.Error(err)
	}
}

func TestCountries(t *testing.T) {
	template := "INSERT INTO `floop` VALUES ('{country}','{country:up:0}','{country}','{country:down:1}')"
	_, err := ParseTemplate(template)
	if err != nil {
		t.Error(err)
	}
}

func TestInteger(t *testing.T) {
	template := "{int:5:6}"
	tp, err := ParseTemplate(template)
	if err != nil {
		t.Error(err)
	}
	c, err := strconv.Atoi(tp)
	if err != nil {
		t.Error(err)
	}
	if c < 5 || c > 6 {
		t.Error("Integer out of range")
	}
}

func TestNowOrdinal(t *testing.T) {
	template := "{now:1}"
	_, err := ParseTemplate(template)
	if err == nil {
		t.Error("Did not return an error on an invalid {now} ordinal")
	}
}

func TestGuidOrdinal(t *testing.T) {
	template := "{guid:1}"
	_, err := ParseTemplate(template)
	if err == nil {
		t.Error("Did not return an error on an invalid {gui} ordinal")
	}
}

func BenchmarkBuildSQL(b *testing.B) {
	template := "INSERT INTO `floop` VALUES ('{guid}','{guid:0}',{int:-2000:0},{int:100:1000},{int:1:40},'{now}','{now:0}','{char:2:up}',NULL)"

	for n := 0; n < b.N; n++ {
		ParseTemplate(template)
	}
}
