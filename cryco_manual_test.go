package cryco_test

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/mengstr/cryco"
)

type testStruct struct {
	I  int64 `default:"afAlg7xxVE_9lyNZI5WqfRV0mqr0YTAElWo_Oiw="`
	I1 int64 `json:"ID"`
	i2 int64
	I4 int64
	F  float64 `default:"s6UB2c7qbmJ6uVVxyTHd8Tg5wiRHA8tIcnwtrIF1dw=="`
	S  string  `default:"Q_mpS5vn_fl0cY3L44foaz-LJ3p5Qn9fCmmFWBH8xbHAnw=="`
}

type testStructNonExported struct {
	I int64   `default:"afAlg7xxVE_9lyNZI5WqfRV0mqr0YTAElWo_Oiw="`
	j int64   `default:"afAlg7xxVE_9lyNZI5WqfRV0mqr0YTAElWo_Oiw="`
	F float64 `default:"s6UB2c7qbmJ6uVVxyTHd8Tg5wiRHA8tIcnwtrIF1dw=="`
	S string  `default:"Q_mpS5vn_fl0cY3L44foaz-LJ3p5Qn9fCmmFWBH8xbHAnw=="`
}

type testStructNotHandled struct {
	I int32 `default:"afAlg7xxVE_9lyNZI5WqfRV0mqr0YTAElWo_Oiw="`
}

const cfgOk = `
#
# Hello world

I=2
S=Bletch
F=2.2
`

const cfgBadInt = `
#
# Hello world

I=3.1
S=Bletch
F=3.3
`

const cfgBadFloat = `
#
# Hello world

I=4
S=Bletch
F=4x4
`
const cfgNewVar = `
#
# Hello world

I=5
S=Bletch
F=5.5
X=5
`
const cfgMisingEqual = `
#
# Hello world

I=6
S:Bletch
F=5.5
`

func Test(t *testing.T) {
	var err error
	key := []byte{65, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 65}
	os.Setenv("KEYcrycotest", "QWFhYWFhYWFhYWFhYWFhQQ==")

	// var a testStruct
	testStructOk := testStruct{}
	testStructNE := testStructNonExported{}
	testStructNH := testStructNotHandled{}

	I := 123
	err = cryco.SetDefaults(I, key)
	if !errors.Is(err, cryco.ErrNotStructPtr) {
		t.Errorf("SetDefault: Passing int - expected '%v' got '%v'", cryco.ErrNotStructPtr, err)
	}

	err = cryco.SetDefaults(&I, key)
	if !errors.Is(err, cryco.ErrNotStructPtr) {
		t.Errorf("SetDefault: Passing &int - expected '%v' got '%v'", cryco.ErrNotStructPtr, err)
	}
	err = cryco.SetDefaults(testStructOk, key)
	if !errors.Is(err, cryco.ErrNotStructPtr) {
		t.Errorf("SetDefault: Passing struct - expected '%v' got '%v'", cryco.ErrNotStructPtr, err)
	}

	// cryco.ParseFiles(&I, "a.txt", "b.txt", "c.txt")

	err = cryco.SetDefaults(testStructOk, key)
	if !errors.Is(err, cryco.ErrNotStructPtr) {
		t.Errorf("Passing struct - expected '%v' got '%v'", cryco.ErrNotStructPtr, err)
	}

	err = cryco.SetDefaults(&testStructOk, key)
	if err != nil {
		t.Errorf("Got error '%v' when passing &struct", err)
	}

	if testStructOk.I != 1 {
		t.Errorf("I is '%v', expected 1", testStructOk.I)
	}
	if testStructOk.F != 1.1 {
		t.Errorf("F is '%v', expected 1.1", testStructOk.F)
	}
	if testStructOk.S != "Foobar" {
		t.Errorf("S is '%v', expected 'Foobar'", testStructOk.S)
	}

	var rdrs []io.Reader

	rdrs = nil
	rdrs = append(rdrs, strings.NewReader(cfgOk))

	err = cryco.ParseReaders(testStructOk, rdrs)
	if !errors.Is(err, cryco.ErrNotStructPtr) {
		t.Errorf("Passing struct - expected '%v' got '%v'", cryco.ErrNotStructPtr, err)
	}

	err = cryco.ParseReaders(&testStructOk, rdrs)
	if err != nil {
		t.Errorf("ParseReaders returned %v", err)
	}
	if testStructOk.I != 2 {
		t.Errorf("I is '%v', expected 2", testStructOk.I)
	}
	if testStructOk.F != 2.2 {
		t.Errorf("F is '%v', expected 2.2", testStructOk.F)
	}
	if testStructOk.S != "Bletch" {
		t.Errorf("S is '%v', expected 'Bletch'", testStructOk.S)
	}

	rdrs = nil
	rdrs = append(rdrs, strings.NewReader(cfgBadInt))
	err = cryco.ParseReaders(&testStructOk, rdrs)
	log.Println(errors.Is(err, cryco.ErrParse))

	rdrs = nil
	rdrs = append(rdrs, strings.NewReader(cfgBadFloat))
	err = cryco.ParseReaders(&testStructOk, rdrs)
	log.Println(errors.Is(err, cryco.ErrParse))

	// Test missing = in file
	rdrs = nil
	rdrs = append(rdrs, strings.NewReader(cfgMisingEqual))
	err = cryco.ParseReaders(&testStructOk, rdrs)
	if !errors.Is(err, cryco.ErrBadFileFormat) {
		t.Errorf("Expected '%v' got '%v'", cryco.ErrBadFileFormat, err)
	}

	// Test non-exported field
	err = cryco.SetDefaults(&testStructNE, key)
	if !errors.Is(err, cryco.ErrNotExported) {
		t.Errorf("Non-exported field - expected '%v' got '%v'", cryco.ErrNotExported, err)
	}

	// Test not supported field type
	err = cryco.SetDefaults(&testStructNH, key)
	if !errors.Is(err, cryco.ErrUnhandledType) {
		t.Errorf("Unhandled int32 - expected '%v' got '%v'", cryco.ErrUnhandledType, err)
	}

}
