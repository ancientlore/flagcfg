package flagcfg

import (
	"flag"
	//	"github.com/ancientlore/flagcfg"
	"strings"
	"testing"
	"time"
)

const data = `
string = "hello"
int = 32
uint = 33
int64 = 64
uint64 = 65
duration = "10s"
bool = true
float64 = 64.32
anotherInt = 44
some_var = "hello again"
another_var = "oh the pain"
`

func TestTypes(t *testing.T) {
	s := flag.NewFlagSet("TestTypes", flag.PanicOnError)
	String := s.String("string", "", "")
	Int := s.Int("int", 0, "")
	Uint := s.Uint("uint", 0, "")
	Int64 := s.Int64("int64", 0, "")
	Uint64 := s.Uint64("uint64", 0, "")
	Duration := s.Duration("duration", time.Second, "")
	Bool := s.Bool("bool", false, "")
	Float64 := s.Float64("float64", 0.0, "")
	AnotherInt := s.Int("anotherInt", 0, "")
	SomeVar := s.String("some-var", "", "")
	AnotherVar := s.String("another.var", "", "")

	s.Parse(strings.Split("-anotherInt 55", " "))

	err := ParseSet([]byte(data), s)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if *String != "hello" {
		t.Errorf("String should be \"hello\": \"%s\"", *String)
	}
	if *Int != 32 {
		t.Errorf("Int should be 32: %d", *Int)
	}
	if *Uint != 33 {
		t.Errorf("Uint should be 33: %d", *Uint)
	}
	if *Int64 != 64 {
		t.Errorf("Int64 should be 64: %d", *Int64)
	}
	if *Uint64 != 65 {
		t.Errorf("Uint64 should be 65: %d", *Uint64)
	}
	if (*Duration).String() != "10s" {
		t.Errorf("Duration should be 10s: %s", (*Duration).String())
	}
	if *Bool != true {
		t.Errorf("Bool should be true: %v", *Bool)
	}
	if *Float64 != 64.32 {
		t.Errorf("Float64 should be 64.32: %f", *Float64)
	}
	if *AnotherInt != 55 {
		t.Errorf("AnotherInt should be 55: %d", *AnotherInt)
	}
	if *SomeVar != "hello again" {
		t.Errorf("SomeVar should be \"hello again\": \"%s\"", *SomeVar)
	}
	if *AnotherVar != "oh the pain" {
		t.Errorf("AnotherVar should be \"oh the pain\": \"%s\"", *AnotherVar)
	}
}
