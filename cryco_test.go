package cryco

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	envKeyName   = "KEYcrycotest"
	keyGoodB64   = "QWFhYWFhYWFhYWFhYWFhQQ=="                         // AaaaaaaaaaaaaaaA
	keyWrongB64  = "WGFhYWFhYWFhYWFhYWFhWA=="                         // XaaaaaaaaaaaaaaX
	keyBadB64    = "WFhYWFhYWFhYWFhYWFhQQ=="                          // Too short key
	goodBase64   = "R29vZA=="                                         // Good as BASE64
	badBase64    = "R29!ZA=="                                         // invalid character in BASE64
	shortCipher  = "QUJDMTIz"                                         // ABC123 as BASE64
	cipherABC123 = "iVKgKeNMAPVGXU2XJP__yFHDMP0tj5kyRALAsgI0jXWfsg==" // ABC123 encrypted
	cipher1      = "VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4="         // 1 encrypted
	cipher1d1    = "ZWuWGl8sOQ_gMFsz_l0IllFBmYemsNAennDesZ81ew=="     // 1.1 encrypted
	cipherOne    = "ZfgUJkrHKNc3_1kOGq0441Guz7GIOs9FzxuQOHfaTg=="     // One encrypted
	cipher2      = "G1XjPYt9TUEEwoSXXNFnzpcDJTd_FqpNlBNFUPg="         // 2 encrypted
	cipher2d2    = "O4F8tmdIqB9Rd1c_eaHVdmeO74XmP1Vvn37QcqHyTA=="     // 2.2 encrypted
	cipherTwo    = "xhGpyDq3QzaEE0mLfN7fv9Zylsl5Zk0Co7srhOe9GA=="     // Two encrypted
)

var (
	bKeyZero  = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	bKeyShort = []byte{65, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97}
	bKeyGood  = []byte{65, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 65}
	bKeyWrong = []byte{111, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 111}
)

func Test_exeName(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"", "crycotest", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exeName()
			if (err != nil) != tt.wantErr {
				t.Errorf("exeName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("exeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetKey(t *testing.T) {
	tests := []struct {
		name        string
		e           string
		want        []byte
		wantErr     bool
		wantErrType error
	}{
		{"nothing", "", bKeyZero, false, nil},
		{"bad env key", keyBadB64, bKeyZero, true, ErrBase64},
		{"good env key", keyGoodB64, bKeyGood, false, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(envKeyName)
			if tt.e != "" {
				os.Setenv(envKeyName, tt.e)
			}
			got, err := GetKey()
			os.Unsetenv(envKeyName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.wantErrType) {
				t.Errorf("GetKey() error = '%v', wantErr '%v'", err, tt.wantErrType)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		bKey      []byte
		cipherB64 string
	}
	tests := []struct {
		name        string
		args        args
		want        string
		wantErr     bool
		wantErrType error
	}{
		{"bad base64", args{bKey: bKeyGood, cipherB64: badBase64}, "", true, ErrBase64},
		{"short key", args{bKey: bKeyShort, cipherB64: goodBase64}, "", true, ErrInternal},
		{"short cipherdata", args{bKey: bKeyGood, cipherB64: shortCipher}, "", true, ErrInternal},
		{"wrong key", args{bKey: bKeyWrong, cipherB64: cipherABC123}, "", true, ErrInvalidKey},
		{"good cipherdata", args{bKey: bKeyGood, cipherB64: cipherABC123}, "ABC123", false, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.bKey, tt.args.cipherB64)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.wantErrType) {
				t.Errorf("Decrypt() error = '%v', wantErr '%v'", err, tt.wantErrType)
				return
			}
			if got != tt.want {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setValue(t *testing.T) {
	type testStruct struct {
		I  int64   `default:"VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4="`     // 1
		i2 int64   `default:"VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4="`     // 1
		F  float64 `default:"ZWuWGl8sOQ_gMFsz_l0IllFBmYemsNAennDesZ81ew=="` // 1.1
		S  string  `default:"ZfgUJkrHKNc3_1kOGq0441Guz7GIOs9FzxuQOHfaTg=="` // One
	}
	var st testStruct

	type args struct {
		p     interface{}
		field string
		value string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantErrType error
	}{
		{"Int64", args{&st, "I", "2"}, false, nil},
		{"Not exported int64", args{&st, "i2", "2"}, true, ErrNotExported},
		{"Float64", args{&st, "F", "2.2"}, false, nil},
		{"String", args{&st, "S", "Two"}, false, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setValue(tt.args.p, tt.args.field, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("setValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.wantErrType) {
				t.Errorf("setValue() error = '%v', wantErr '%v'", err, tt.wantErrType)
				return
			}

		})
	}
	t.Run("setValue() results", func(t *testing.T) {
		for _, tt := range tests {
			_ = setValue(tt.args.p, tt.args.field, tt.args.value)
		}
		want := testStruct{2, 0, 2.2, "Two"}
		if st != want {
			t.Errorf("setValue() got %v, want %v", st, want)
			return
		}
	})
}

func TestSetDefaults(t *testing.T) {
	type testBadStruct struct {
		I  int64   `default:"VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4="`     // 1
		i2 int64   `default:"VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4="`     // 1
		F  float64 `default:"ZWuWGl8sOQ_gMFsz_l0IllFBmYemsNAennDesZ81ew=="` // 1.1
		S  string  `default:"ZfgUJkrHKNc3_1kOGq0441Guz7GIOs9FzxuQOHfaTg=="` // One
	}
	type testGoodStruct struct {
		I  int64   `default:"VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4="`     // 1
		i2 int64   `other:"Foobar"`                                         //
		F  float64 `default:"ZWuWGl8sOQ_gMFsz_l0IllFBmYemsNAennDesZ81ew=="` // 1.1
		S  string  `default:"ZfgUJkrHKNc3_1kOGq0441Guz7GIOs9FzxuQOHfaTg=="` // One
	}
	var stBad testBadStruct
	var stGood testGoodStruct
	var testInt int64

	type args struct {
		struc interface{}
		bKey  []byte
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantErrType error
	}{
		{"Int", args{testInt, bKeyGood}, true, ErrNotStructPtr},
		{"Int pointer", args{&testInt, bKeyGood}, true, ErrNotStructPtr},
		{"Not exported", args{&stBad, bKeyGood}, true, ErrNotExported},
		{"Struct", args{stGood, bKeyGood}, true, ErrNotStructPtr},
		{"Struct pointer", args{&stGood, bKeyGood}, false, nil},
		{"Wrong key", args{&stGood, bKeyWrong}, true, ErrInvalidKey},
		{"Short key", args{&stGood, bKeyShort}, true, ErrInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetDefaults(tt.args.struc, tt.args.bKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDefaults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.wantErrType) {
				t.Errorf("SetDefaults() error = '%v', wantErr '%v'", err, tt.wantErrType)
				return
			}
		})
	}
	t.Run("SetDDefault values", func(t *testing.T) {
		_ = SetDefaults(&stGood, bKeyGood)
		want := testGoodStruct{1, 0, 1.1, "One"}
		if stGood != want {
			t.Errorf("setValue() got %v, want %v", stGood, want)
			return
		}
	})
}

const cfgOk1 = `
#
# Hello world

I=3
S=Three
F=3.3
`

const cfgOk2 = `
#
# Hello world

I=4
S=Four
F=4.4
`

const cfgEmpty = `
#
# Hello world
`

func TestParseReaders(t *testing.T) {
	var rdrs0 []io.Reader
	rdrs0 = nil
	var rdrs1 []io.Reader
	rdrs1 = nil
	rdrs1 = append(rdrs1, strings.NewReader(cfgOk1))
	rdrs1 = append(rdrs1, strings.NewReader(cfgOk2))
	var rdrs2 []io.Reader
	rdrs2 = nil
	rdrs2 = append(rdrs2, strings.NewReader(cfgEmpty))
	rdrs2 = append(rdrs2, strings.NewReader(cfgOk2))

	type testGoodStruct struct {
		I  int64   `default:"VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4="`     // 1
		i2 int64   `other:"Foobar"`                                         //
		F  float64 `default:"ZWuWGl8sOQ_gMFsz_l0IllFBmYemsNAennDesZ81ew=="` // 1.1
		S  string  `default:"ZfgUJkrHKNc3_1kOGq0441Guz7GIOs9FzxuQOHfaTg=="` // One
	}
	var stGood testGoodStruct

	type args struct {
		struc   interface{}
		readers []io.Reader
	}
	tests := []struct {
		name        string
		args        args
		key         string
		want        testGoodStruct
		wantErr     bool
		wantErrType error
	}{
		{"Wrong key", args{&stGood, rdrs0}, keyWrongB64, testGoodStruct{0, 0, 0.0, ""}, true, ErrInvalidKey},
		{"Bad key", args{&stGood, rdrs0}, keyBadB64, testGoodStruct{0, 0, 0.0, ""}, true, ErrBase64},
		{"No readers", args{&stGood, rdrs0}, keyGoodB64, testGoodStruct{1, 0, 1.1, "One"}, false, nil},
		{"Use first of 2", args{&stGood, rdrs1}, keyGoodB64, testGoodStruct{3, 0, 3.3, "Three"}, false, nil},
		{"Use second of two", args{&stGood, rdrs2}, keyGoodB64, testGoodStruct{4, 0, 4.4, "Four"}, false, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(envKeyName, tt.key)
			err := ParseReaders(tt.args.struc, tt.args.readers)
			os.Unsetenv(envKeyName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.wantErrType) {
				t.Errorf("ParseReaders() error = '%v', wantErr '%v'", err, tt.wantErrType)
				return
			}
			if stGood != tt.want {
				t.Errorf("ParseReaders() got = '%v', want '%v'", stGood, tt.want)
				return
			}
		})
	}
}

func TestParseFiles(t *testing.T) {
	f1, err1 := ioutil.TempFile("", "cryco")
	f2, err2 := ioutil.TempFile("", "cryco")
	f3, err3 := ioutil.TempFile("", "cryco")
	if err1 != nil || err2 != nil || err3 != nil {
		t.Error("Can't create tempfiles")
		t.Fail()
	}
	defer os.Remove(f1.Name())
	defer os.Remove(f2.Name())
	defer os.Remove(f3.Name())
	f1.WriteString(cfgOk1)
	f2.WriteString(cfgOk2)
	f3.WriteString(cfgEmpty)
	f1.Sync()
	f2.Sync()
	f3.Sync()

	type testGoodStruct struct {
		I  int64   `default:"VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4="`     // 1
		i2 int64   `other:"Foobar"`                                         //
		F  float64 `default:"ZWuWGl8sOQ_gMFsz_l0IllFBmYemsNAennDesZ81ew=="` // 1.1
		S  string  `default:"ZfgUJkrHKNc3_1kOGq0441Guz7GIOs9FzxuQOHfaTg=="` // One
	}
	var stGood testGoodStruct

	type args struct {
		struc     interface{}
		filenames []string
	}
	tests := []struct {
		name        string
		args        args
		want        testGoodStruct
		wantErr     bool
		wantErrType error
	}{
		{"UseFirst", args{&stGood, []string{f1.Name(), f2.Name()}}, testGoodStruct{3, 0, 3.3, "Three"}, false, nil},
		{"UseSecond", args{&stGood, []string{f3.Name(), f2.Name()}}, testGoodStruct{4, 0, 4.4, "Four"}, false, nil},
		{"Empty&Nonexistent", args{&stGood, []string{f3.Name(), "NoFile.TXT"}}, testGoodStruct{1, 0, 1.1, "One"}, false, nil},
		{"Nofiles", args{&stGood, []string{}}, testGoodStruct{1, 0, 1.1, "One"}, false, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(envKeyName, keyGoodB64)
			err := ParseFiles(tt.args.struc, tt.args.filenames...)
			os.Unsetenv(envKeyName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, tt.wantErrType) {
				t.Errorf("ParseFiles() error = '%v', wantErr '%v'", err, tt.wantErrType)
				return
			}
			if stGood != tt.want {
				t.Errorf("ParseFiles() got = '%v', want '%v'", stGood, tt.want)
				return
			}

		})
	}
}
