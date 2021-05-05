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
	cipher5      = "-jOb83fMxevZJ5VwDRKrNu8NZdfV9wVYrSvzS3M="         // 5 encrypted
	cipher5d5    = "CIgq4gXo_gew-86-Lpcla-6UcfvDpLCjy_FU6shTyw=="     // 5.5 encrypted
	cipherFive   = "brSSFJ8WHPoMbz0-5bNrNl0ixtK23wyyHrWEEy6QT6U="     // Five encrypted
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
			st = testStruct{}
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
		st = testStruct{}
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
		I  int64   `default:"VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4=" env:"EnvI"`     // 1
		i2 int64   `default:"VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4="`                // 1
		F  float64 `default:"ZWuWGl8sOQ_gMFsz_l0IllFBmYemsNAennDesZ81ew==" env:"EnvF"` // 1.1
		S  string  `default:"ZfgUJkrHKNc3_1kOGq0441Guz7GIOs9FzxuQOHfaTg==" env:"EnvS"` // One
	}
	type testGoodStruct struct {
		I  int64   `default:"VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4=" env:"EnvI"`     // 1
		i2 int64   `other:"Foobar"`                                                    //
		F  float64 `default:"ZWuWGl8sOQ_gMFsz_l0IllFBmYemsNAennDesZ81ew==" env:"EnvF"` // 1.1
		S  string  `default:"ZfgUJkrHKNc3_1kOGq0441Guz7GIOs9FzxuQOHfaTg==" env:"EnvS"` // One
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
			stGood = testGoodStruct{}
			stBad = testBadStruct{}
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
		stGood = testGoodStruct{}
		stBad = testBadStruct{}
		_ = SetDefaults(&stGood, bKeyGood)
		want := testGoodStruct{1, 0, 1.1, "One"}
		if stGood != want {
			t.Errorf("setValue() got %v, want %v", stGood, want)
			return
		}
	})
}

const cfgOk3 = `
#
# Hello world

I=3
S=Three
F=3.3
`

const cfgOk4 = `
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

//
func setEnvs(envs string) {
	os.Unsetenv("EnvI")
	os.Unsetenv("EnvF")
	os.Unsetenv("EnvS")
	for _, c := range envs {
		switch string(c) {
		case "I":
			os.Setenv("EnvI", cipher5)
		case "F":
			os.Setenv("EnvF", cipher5d5)
		case "S":
			os.Setenv("EnvS", cipherFive)
		case "i":
			os.Setenv("EnvI", "(6)")
		case "f":
			os.Setenv("EnvF", "(6.6)")
		case "s":
			os.Setenv("EnvS", "(Six)")
		}
	}
}

func resetReaders(rdrs0 *[]io.Reader, rdrs3 *[]io.Reader, rdrs4 *[]io.Reader) {
	const cfgOk3 = `
	#
	# Hello world
	
	I=3
	S=Three
	F=3.3
	`

	const cfgOk4 = `
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
	*rdrs0 = nil
	*rdrs3 = nil
	*rdrs3 = append(*rdrs3, strings.NewReader(cfgOk3))
	*rdrs3 = append(*rdrs3, strings.NewReader(cfgOk4))
	*rdrs4 = nil
	*rdrs4 = append(*rdrs4, strings.NewReader(cfgEmpty))
	*rdrs4 = append(*rdrs4, strings.NewReader(cfgOk4))

}

func TestParseReaders(t *testing.T) {
	var rdrs0 []io.Reader
	var rdrs3 []io.Reader
	var rdrs4 []io.Reader
	type testGoodStruct struct {
		I  int64   `default:"VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4=" env:"EnvI"`     // 1
		i2 int64   `other:"Foobar"`                                                    //
		F  float64 `default:"ZWuWGl8sOQ_gMFsz_l0IllFBmYemsNAennDesZ81ew==" env:"EnvF"` // 1.1
		S  string  `default:"ZfgUJkrHKNc3_1kOGq0441Guz7GIOs9FzxuQOHfaTg==" env:"EnvS"` // One
	}
	var stGood testGoodStruct

	type args struct {
		struc       interface{}
		readersName string
	}
	tests := []struct {
		name        string
		args        args
		envs        string
		key         string
		want        testGoodStruct
		wantErr     bool
		wantErrType error
	}{
		{"Wrong key", args{&stGood, "rdrs0"}, "", keyWrongB64, testGoodStruct{0, 0, 0.0, ""}, true, ErrInvalidKey},
		{"Bad key", args{&stGood, "rdrs0"}, "", keyBadB64, testGoodStruct{0, 0, 0.0, ""}, true, ErrBase64},
		{"No readers", args{&stGood, "rdrs0"}, "", keyGoodB64, testGoodStruct{1, 0, 1.1, "One"}, false, nil},
		{"Use first of 2", args{&stGood, "rdrs3"}, "", keyGoodB64, testGoodStruct{3, 0, 3.3, "Three"}, false, nil},
		{"Use second of two", args{&stGood, "rdrs4"}, "", keyGoodB64, testGoodStruct{4, 0, 4.4, "Four"}, false, nil},

		{"Wrong key", args{&stGood, "rdrs0"}, "I", keyWrongB64, testGoodStruct{0, 0, 0.0, ""}, true, ErrInvalidKey},
		{"Bad key", args{&stGood, "rdrs0"}, "I", keyBadB64, testGoodStruct{0, 0, 0.0, ""}, true, ErrBase64},
		{"No readers", args{&stGood, "rdrs0"}, "I", keyGoodB64, testGoodStruct{5, 0, 1.1, "One"}, false, nil},
		{"Use first of 2", args{&stGood, "rdrs3"}, "I", keyGoodB64, testGoodStruct{5, 0, 3.3, "Three"}, false, nil},
		{"Use second of two", args{&stGood, "rdrs4"}, "I", keyGoodB64, testGoodStruct{5, 0, 4.4, "Four"}, false, nil},

		{"Wrong key", args{&stGood, "rdrs0"}, "Ifs", keyWrongB64, testGoodStruct{0, 0, 0.0, ""}, true, ErrInvalidKey},
		{"Bad key", args{&stGood, "rdrs0"}, "Ifs", keyBadB64, testGoodStruct{0, 0, 0.0, ""}, true, ErrBase64},
		{"No readers", args{&stGood, "rdrs0"}, "Ifs", keyGoodB64, testGoodStruct{5, 0, 6.6, "Six"}, false, nil},
		{"Use first of 2", args{&stGood, "rdrs3"}, "Ifs", keyGoodB64, testGoodStruct{5, 0, 6.6, "Six"}, false, nil},
		{"Use second of two", args{&stGood, "rdrs4"}, "Ifs", keyGoodB64, testGoodStruct{5, 0, 6.6, "Six"}, false, nil},
	}
	_ = rdrs0
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.args.struc = &testGoodStruct{}
			resetReaders(&rdrs0, &rdrs3, &rdrs4)
			setEnvs(tt.envs)
			os.Setenv(envKeyName, tt.key)
			switch tt.args.readersName {
			case "rdrs0":
				err = ParseReaders(tt.args.struc, rdrs0)
			case "rdrs3":
				err = ParseReaders(tt.args.struc, rdrs3)
			case "rdrs4":
				err = ParseReaders(tt.args.struc, rdrs4)
			}
			os.Unsetenv(envKeyName)
			setEnvs("")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.wantErrType) {
				t.Errorf("ParseReaders() error = '%v', wantErr '%v'", err, tt.wantErrType)
				return
			}
			if !reflect.DeepEqual(tt.args.struc, &tt.want) {
				t.Errorf("ParseReaders() got = '%v', want '%v'", tt.args.struc, tt.want)
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
	f1.WriteString(cfgOk3)
	f2.WriteString(cfgOk4)
	f3.WriteString(cfgEmpty)
	f1.Sync()
	f2.Sync()
	f3.Sync()

	type testGoodStruct struct {
		I  int64   `default:"VsA2dNX5VkXVwqC-JMHQWCtUWNZ78OPz61OKbB4=" env:"EnvI"`     // 1
		i2 int64   `other:"Foobar"`                                                    //
		F  float64 `default:"ZWuWGl8sOQ_gMFsz_l0IllFBmYemsNAennDesZ81ew==" env:"EnvF"` // 1.1
		S  string  `default:"ZfgUJkrHKNc3_1kOGq0441Guz7GIOs9FzxuQOHfaTg==" env:"EnvS"` // One
	}
	var stGood testGoodStruct

	type args struct {
		struc     interface{}
		filenames []string
	}
	tests := []struct {
		name        string
		args        args
		envs        string
		want        testGoodStruct
		wantErr     bool
		wantErrType error
	}{
		{"UseFirst", args{&stGood, []string{f1.Name(), f2.Name()}}, "", testGoodStruct{3, 0, 3.3, "Three"}, false, nil},
		{"UseSecond", args{&stGood, []string{f3.Name(), f2.Name()}}, "", testGoodStruct{4, 0, 4.4, "Four"}, false, nil},
		{"Empty&Nonexistent", args{&stGood, []string{f3.Name(), "NoFile.TXT"}}, "", testGoodStruct{1, 0, 1.1, "One"}, false, nil},
		{"Nofiles", args{&stGood, []string{}}, "", testGoodStruct{1, 0, 1.1, "One"}, false, nil},

		{"UseFirst", args{&stGood, []string{f1.Name(), f2.Name()}}, "I", testGoodStruct{5, 0, 3.3, "Three"}, false, nil},
		{"UseSecond", args{&stGood, []string{f3.Name(), f2.Name()}}, "I", testGoodStruct{5, 0, 4.4, "Four"}, false, nil},
		{"Empty&Nonexistent", args{&stGood, []string{f3.Name(), "NoFile.TXT"}}, "I", testGoodStruct{5, 0, 1.1, "One"}, false, nil},
		{"Nofiles", args{&stGood, []string{}}, "I", testGoodStruct{5, 0, 1.1, "One"}, false, nil},

		{"UseFirst", args{&stGood, []string{f1.Name(), f2.Name()}}, "IFS", testGoodStruct{5, 0, 5.5, "Five"}, false, nil},
		{"UseSecond", args{&stGood, []string{f3.Name(), f2.Name()}}, "IFS", testGoodStruct{5, 0, 5.5, "Five"}, false, nil},
		{"Empty&Nonexistent", args{&stGood, []string{f3.Name(), "NoFile.TXT"}}, "IFS", testGoodStruct{5, 0, 5.5, "Five"}, false, nil},
		{"Nofiles", args{&stGood, []string{}}, "IFS", testGoodStruct{5, 0, 5.5, "Five"}, false, nil},

		{"UseFirst", args{&stGood, []string{f1.Name(), f2.Name()}}, "ifs", testGoodStruct{6, 0, 6.6, "Six"}, false, nil},
		{"UseSecond", args{&stGood, []string{f3.Name(), f2.Name()}}, "ifs", testGoodStruct{6, 0, 6.6, "Six"}, false, nil},
		{"Empty&Nonexistent", args{&stGood, []string{f3.Name(), "NoFile.TXT"}}, "ifs", testGoodStruct{6, 0, 6.6, "Six"}, false, nil},
		{"Nofiles", args{&stGood, []string{}}, "ifs", testGoodStruct{6, 0, 6.6, "Six"}, false, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.struc = &testGoodStruct{}
			setEnvs(tt.envs)
			os.Setenv(envKeyName, keyGoodB64)
			err := ParseFiles(tt.args.struc, tt.args.filenames...)
			os.Unsetenv(envKeyName)
			setEnvs("")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, tt.wantErrType) {
				t.Errorf("ParseFiles() error = '%v', wantErr '%v'", err, tt.wantErrType)
				return
			}
			if !reflect.DeepEqual(tt.args.struc, &tt.want) {
				t.Errorf("ParseFiles() got = '%v', want '%v'", stGood, tt.want)
				return
			}

		})
	}
}
