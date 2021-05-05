package cryco

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	keylen = 16 // AES128
)

var (
	key = "" // Set by go build -ldflags "-X github.com/mengstr/cryco.key=......."
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

// GetKey Returns the active key decoded from its original Base64
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
func Decrypt(bKey []byte, cipherB64 string) (string, error) {
	// If string is bracketed with paranthesis () then if should be treated as cleartext so
	// remove the paranthesises and return as is
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
func setValue(p interface{}, field string, value string) error {
	log.Printf("setValue '%s' to '%s' ", field, value)
	// Elem returns the value that the pointer p points to.
	v := reflect.ValueOf(p).Elem()
	f := v.FieldByName(field)
	// make sure that this field is defined, and can be changed.
	if !f.IsValid() || !f.CanSet() {
		return fmt.Errorf("%w - field %s", ErrNotExported, field)
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
			return fmt.Errorf("%w %v", ErrParse, err)
		}
		f.SetFloat(f64)
		return nil
	}

	return fmt.Errorf("%w %s", ErrUnhandledType, f.Kind())

}

// SetFromEnv ...
func SetFromEnv(struc interface{}, bKey []byte) error {
	log.Println("SETENVS")
	var err error

	if reflect.TypeOf(struc).Kind() != reflect.Ptr || reflect.ValueOf(struc).Elem().Kind() != reflect.Struct {
		return ErrNotStructPtr
	}

	for i := 0; i < reflect.ValueOf(struc).Elem().NumField(); i++ {
		envname := reflect.ValueOf(struc).Elem().Type().Field(i).Tag.Get("env")
		_, ok := reflect.ValueOf(struc).Elem().Type().Field(i).Tag.Lookup("")
		if envname == "" && !ok {
			continue
		}
		defv, ok := os.LookupEnv(envname)
		if !ok {
			continue
		}
		defv, err = Decrypt(bKey, defv)
		if err != nil {
			return err
		}
		err = setValue(struc, reflect.ValueOf(struc).Elem().Type().Field(i).Name, defv)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetDefaults ...
func SetDefaults(struc interface{}, bKey []byte) error {
	log.Println("SETDEFAULT")
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
		defv, err = Decrypt(bKey, defv)
		if err != nil {
			return err
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
	log.Println("PARSEREADERS")
	bKey, err := GetKey()
	if err != nil {
		return err
	}
	if err := SetDefaults(struc, bKey); err != nil {
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
			field := strings.TrimSpace(ss[0])
			value := strings.TrimSpace(ss[1])
			if err := setValue(struc, field, value); err != nil {
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

	return SetFromEnv(struc, bKey)
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
