package cryco

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrParse Some kind of error when parsing a value from the config/setings files
	ErrParse = errors.New("Parse error")
	// ErrNotStructPtr The receiving veriable is not a struct passed as a pointer
	ErrNotStructPtr = errors.New("Argument is not a pointer to a struct")
	// ErrNotExported ...
	ErrNotExported = errors.New("Can't use non-exported fields in struct")
	// ErrBadFileFormat ...
	ErrBadFileFormat = errors.New("Bad file format")
	// ErrUnhandledType ...
	ErrUnhandledType = errors.New("Not handling field type")
)

func setValue(p interface{}, key string, value string) error {
	// Elem returns the value that the pointer u points to.
	v := reflect.ValueOf(p).Elem()
	f := v.FieldByName(key)
	// make sure that this field is defined, and can be changed.
	if !f.IsValid() || !f.CanSet() {
		return fmt.Errorf("%w - field %s", ErrNotExported, key)
	}
	if f.Kind() == reflect.String {
		f.SetString(value)
		return nil
	}
	if f.Kind() == reflect.Int64 {
		i64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("%w %v", ErrParse, err)
		}
		f.SetInt(i64)
		return nil
	}
	if f.Kind() == reflect.Float64 {
		f64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("Can't parse %v as float64", value)
		}
		f.SetFloat(f64)
		return nil
	}

	return fmt.Errorf("%w %s", ErrUnhandledType, f.Kind())

}

// SetDefaults ...
func SetDefaults(struc interface{}) error {
	var err error

	if reflect.TypeOf(struc).Kind() != reflect.Ptr || reflect.ValueOf(struc).Elem().Kind() != reflect.Struct {
		return ErrNotStructPtr
	}

	for i := 0; i < reflect.ValueOf(struc).Elem().NumField(); i++ {
		defv := reflect.ValueOf(struc).Elem().Type().Field(i).Tag.Get("default")
		_, ok := reflect.ValueOf(struc).Elem().Type().Field(i).Tag.Lookup("")
		if defv == "" && !ok {
			continue
		}
		err = setValue(struc, reflect.ValueOf(struc).Elem().Type().Field(i).Name, defv)
		if err != nil {
			return err
		}
	}
	return nil
}

// ParseReaders parses data from one or more io.Readers
func ParseReaders(struc interface{}, readers []io.Reader) error {
	if err := SetDefaults(struc); err != nil {
		return err
	}

	processed := false
	cnt := 0
	for _, r := range readers {
		cnt++
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			s := strings.TrimSpace(scanner.Text())
			if s == "" || string(s[0]) == "#" {
				continue
			}
			ss := strings.SplitN(s, "=", 2)
			if len(ss) < 2 {
				return fmt.Errorf("%w, missing = at '%s'", ErrBadFileFormat, s)
			}
			key := strings.TrimSpace(ss[0])
			value := strings.TrimSpace(ss[1])
			if err := setValue(struc, key, value); err != nil {
				return err
			}
			processed = true
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		if processed {
			break
		}
	}
	return nil
}

// ParseFiles tries to parse each file in the list and stops after the first parseable file.
func ParseFiles(struc interface{}, filenames ...string) error {
	var rdrs []io.Reader
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			continue
		}
		defer f.Close()
		rdrs = append(rdrs, f)
	}
	return ParseReaders(struc, rdrs)
}
