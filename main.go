package cryco

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	keylen     = 16 // AES128
	tagDefVal  = "def"
	tagFileVal = "fil"
	tagEnvVal  = "env"
)

var (
	// Set by go build -ldflags "-X github.com/mengstr/cryco.key=......."
	key = ""
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
	// ErrBase64 ...
	ErrBase64 = errors.New("Bad Base64 format")
	// ErrInternal ...
	ErrInternal = errors.New("Internal/OS error")
	// ErrInvalidKey ...
	ErrInvalidKey = errors.New("Invalid key")
)

// Returns the sanatized name of the running program
func exeName() (string, error) {
	s, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("%w from os.Executable %v", ErrInternal, err)
	}
	s = filepath.Base(s)

	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	if err != nil {
		return "", fmt.Errorf("%w - %v", ErrInternal, err)
	}
	return reg.ReplaceAllString(s, ""), nil
}

// GetKey Returns the active key decoded from its original Base64 encoding
// The key is retreived from either an environment variable named KEY<executable name>
// or locally from the executable using a variable that got its value patched into
// it during build.
func GetKey() ([]byte, error) {
	name, err := exeName()
	if err != nil {
		return nil, err
	}
	s := os.Getenv("KEY" + name)
	if s == "" {
		s = key
	}
	if s == "" {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil
	}
	bKey, err := base64.StdEncoding.DecodeString(s)
	if err != nil || len(bKey) != keylen {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, fmt.Errorf("%w (%s)", ErrBase64, s)
	}
	return bKey, nil
}

// Decrypt takes a base64 encoded ciphertext string and decrypts it into a cleartext string
// If value string is bracketed with paranthesis () then it should be treated as cleartext so
// remove the paranthesises and return as is
func Decrypt(bKey []byte, cipherB64 string) (string, error) {
	// Cleartext?
	if len(cipherB64) > 1 && cipherB64[0:1] == "(" && cipherB64[len(cipherB64)-1:] == ")" {
		return cipherB64[1 : len(cipherB64)-1], nil
	}
	encryptData, err := base64.URLEncoding.DecodeString(cipherB64)
	if err != nil {
		return "", fmt.Errorf("%w %v", ErrBase64, err)
	}
	cipherBlock, err := aes.NewCipher(bKey)
	if err != nil {
		return "", fmt.Errorf("%w (a)", ErrInternal)
	}
	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", fmt.Errorf("%w (b)", ErrInternal)
	}
	nonceSize := aead.NonceSize()
	if len(encryptData) < nonceSize {
		return "", fmt.Errorf("%w (c)", ErrInternal)
	}
	nonce, cipherText := encryptData[:nonceSize], encryptData[nonceSize:]
	plainData, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("%w %v", ErrInvalidKey, err)
	}
	return string(plainData), nil
}

//
func setFieldValue(p interface{}, field string, value string) error {
	var err error

	if err = CheckParam(p); err != nil {
		return err
	}
	fld := reflect.ValueOf(p).Elem().FieldByName(field)
	if !fld.IsValid() || !fld.CanSet() {
		return fmt.Errorf("%w - field %s", ErrNotExported, field)
	}
	switch fld.Kind() {
	case reflect.String:
		fld.SetString(value)
	case reflect.Int64:
		i64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("%w %v", ErrParse, err)
		}
		fld.SetInt(i64)
	case reflect.Float64:
		f64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("%w %v", ErrParse, err)
		}
		fld.SetFloat(f64)
	default:
		return fmt.Errorf("%w %s", ErrUnhandledType, fld.Kind())
	}
	return nil
}

//
func setValueFromTag(p interface{}, tagType string, tagName string, value string) error {
	var err error
	if err = CheckParam(p); err != nil {
		return err
	}
	// Iterate over the fields until the right one is found
	for i := 0; i < reflect.ValueOf(p).Elem().NumField(); i++ {
		tv, ok := reflect.ValueOf(p).Elem().Type().Field(i).Tag.Lookup(tagType)
		if ok && tv == tagName {
			fieldName := reflect.ValueOf(p).Elem().Type().Field(i).Name
			setFieldValue(p, fieldName, value)
			return nil
		}
	}
	return nil
}

// SetFromEnv ...
func SetFromEnv(p interface{}, bKey []byte) error {
	var err error

	if err = CheckParam(p); err != nil {
		return err
	}
	// Iterate over the fields until the right one is found
	for i := 0; i < reflect.ValueOf(p).Elem().NumField(); i++ {
		tv, ok := reflect.ValueOf(p).Elem().Type().Field(i).Tag.Lookup(tagEnvVal)
		if ok {
			fieldName := reflect.ValueOf(p).Elem().Type().Field(i).Name
			value, ok := os.LookupEnv(tv)
			if !ok {
				continue
			}
			value, err = Decrypt(bKey, value)
			if err != nil {
				return err
			}
			setFieldValue(p, fieldName, value)
		}
	}
	return nil
}

// CheckParam verifies that the param is pointer to a struct
func CheckParam(p interface{}) error {
	if reflect.TypeOf(p).Kind() != reflect.Ptr {
		return ErrNotStructPtr
	}
	if reflect.ValueOf(p).Elem().Kind() != reflect.Struct {
		return ErrNotStructPtr
	}
	return nil
}

// SetDefaults ...
func SetDefaults(struc interface{}, bKey []byte) error {
	var err error
	if err = CheckParam(struc); err != nil {
		return err
	}
	// Scan entire array for default tags and apply them
	e := reflect.ValueOf(struc).Elem()
	for i := 0; i < e.NumField(); i++ {
		fld := e.Type().Field(i)
		value, ok := fld.Tag.Lookup(tagDefVal)
		if ok {
			if value, err = Decrypt(bKey, value); err != nil {
				return err
			}
			if err = setFieldValue(struc, fld.Name, value); err != nil {
				return err
			}
		}
	}
	return nil
}

// ParseReaders parses data from one or more io.Readers.
// First set the dafault values,
// then apply values from the files,
// finally set values from environment variables
func ParseReaders(struc interface{}, readers []io.Reader) error {
	var err error
	if err = CheckParam(struc); err != nil {
		return err
	}
	bKey, err := GetKey()
	if err != nil {
		return err
	}
	if err := SetDefaults(struc, bKey); err != nil {
		return err
	}
	// Process all lines in each reader. As soon as one reader have had
	// any values in it stop processing the rest of the readers.
	processed := false
	cnt := 0
	for _, r := range readers {
		cnt++
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			s := strings.TrimSpace(scanner.Text())
			// Skip empty lines and comments
			if s == "" || string(s[0]) == "#" {
				continue
			}
			// Split line into key (the tag name) and value
			ss := strings.SplitN(s, "=", 2)
			if len(ss) < 2 {
				return fmt.Errorf("%w, missing = at '%s'", ErrBadFileFormat, s)
			}
			// Decrypt the value
			value, err := Decrypt(bKey, strings.TrimSpace(ss[1]))
			if err != nil {
				return err
			}
			// Set the value in the struct, using the tag name
			if err := setValueFromTag(struc, tagFileVal, strings.TrimSpace(ss[0]), value); err != nil {
				return err
			}
			processed = true
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		// Stop scanning files as soon as the first usable file has been fully processed
		if processed {
			break
		}
	}
	// Finish with setting values from envronment variables
	return SetFromEnv(struc, bKey)
}

// ParseFiles tries to parse each file in the list and stops after the first parseable file.
func ParseFiles(struc interface{}, filenames ...string) error {
	var err error

	if err = CheckParam(struc); err != nil {
		return err
	}

	// Opens all specified files...
	var rdrs []io.Reader
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			continue
		}
		defer f.Close()
		rdrs = append(rdrs, f)
	}
	// ...and pass the readers into the ParseReaders() for processing
	return ParseReaders(struc, rdrs)
}
