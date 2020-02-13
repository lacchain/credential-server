package util

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	mrand "math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/cloudflare/cfssl/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ocsp"
)

var (
	rnd = mrand.NewSource(time.Now().UnixNano())
	// ErrNotImplemented used to return errors for functions not implemented
	ErrNotImplemented = errors.New("NOT YET IMPLEMENTED")
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RevocationReasonCodes is a map between string reason codes to integers as defined in RFC 5280
var RevocationReasonCodes = map[string]int{
	"unspecified":          ocsp.Unspecified,
	"keycompromise":        ocsp.KeyCompromise,
	"cacompromise":         ocsp.CACompromise,
	"affiliationchanged":   ocsp.AffiliationChanged,
	"superseded":           ocsp.Superseded,
	"cessationofoperation": ocsp.CessationOfOperation,
	"certificatehold":      ocsp.CertificateHold,
	"removefromcrl":        ocsp.RemoveFromCRL,
	"privilegewithdrawn":   ocsp.PrivilegeWithdrawn,
	"aacompromise":         ocsp.AACompromise,
}

// SecretTag to tag a field as secret as in password, token
const SecretTag = "mask"

// URLRegex is the regular expression to check if a value is an URL
var URLRegex = regexp.MustCompile("(ldap|http)s*://(\\S+):(\\S+)@")

//ECDSASignature forms the structure for R and S value for ECDSA
type ECDSASignature struct {
	R, S *big.Int
}

// RandomString returns a random string
func RandomString(n int) string {
	b := make([]byte, n)

	for i, cache, remain := n-1, rnd.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rnd.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RemoveQuotes removes outer quotes from a string if necessary
func RemoveQuotes(str string) string {
	if str == "" {
		return str
	}
	if (strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'")) ||
		(strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"")) {
		str = str[1 : len(str)-1]
	}
	return str
}

// ReadFile reads a file
func ReadFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

// WriteFile writes a file
func WriteFile(file string, buf []byte, perm os.FileMode) error {
	dir := path.Dir(file)
	// Create the directory if it doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.Wrapf(err, "Failed to create directory '%s' for file '%s'", dir, file)
		}
	}
	return ioutil.WriteFile(file, buf, perm)
}

// FileExists checks to see if a file exists
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// ReplaceTextFile replace a text in a file
func ReplaceTextFile(path , old, new string) error {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	
	newContents := strings.Replace(string(read), old, new, -1)

	err = ioutil.WriteFile(path, []byte(newContents), 0)
	if err != nil {
		panic(err)
	}

	return nil
}

// Marshal to bytes
func Marshal(from interface{}, what string) ([]byte, error) {
	buf, err := json.Marshal(from)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to marshal %s", what)
	}
	return buf, nil
}

// Unmarshal from bytes
func Unmarshal(from []byte, to interface{}, what string) error {
	err := json.Unmarshal(from, to)
	if err != nil {
		return errors.Wrapf(err, "Failed to unmarshal %s", what)
	}
	return nil
}

//GetECPrivateKey get *ecdsa.PrivateKey from key pem
func GetECPrivateKey(raw []byte) (*ecdsa.PrivateKey, error) {
	decoded, _ := pem.Decode(raw)
	if decoded == nil {
		return nil, errors.New("Failed to decode the PEM-encoded ECDSA key")
	}
	ECprivKey, err := x509.ParseECPrivateKey(decoded.Bytes)
	if err == nil {
		return ECprivKey, nil
	}
	key, err2 := x509.ParsePKCS8PrivateKey(decoded.Bytes)
	if err2 == nil {
		switch key.(type) {
		case *ecdsa.PrivateKey:
			return key.(*ecdsa.PrivateKey), nil
		case *rsa.PrivateKey:
			return nil, errors.New("Expecting EC private key but found RSA private key")
		default:
			return nil, errors.New("Invalid private key type in PKCS#8 wrapping")
		}
	}
	return nil, errors.Wrap(err2, "Failed parsing EC private key")
}

//GetRSAPrivateKey get *rsa.PrivateKey from key pem
func GetRSAPrivateKey(raw []byte) (*rsa.PrivateKey, error) {
	decoded, _ := pem.Decode(raw)
	if decoded == nil {
		return nil, errors.New("Failed to decode the PEM-encoded RSA key")
	}
	RSAprivKey, err := x509.ParsePKCS1PrivateKey(decoded.Bytes)
	if err == nil {
		return RSAprivKey, nil
	}
	key, err2 := x509.ParsePKCS8PrivateKey(decoded.Bytes)
	if err2 == nil {
		switch key.(type) {
		case *ecdsa.PrivateKey:
			return nil, errors.New("Expecting RSA private key but found EC private key")
		case *rsa.PrivateKey:
			return key.(*rsa.PrivateKey), nil
		default:
			return nil, errors.New("Invalid private key type in PKCS#8 wrapping")
		}
	}
	return nil, errors.Wrap(err, "Failed parsing RSA private key")
}

// B64Encode base64 encodes bytes
func B64Encode(buf []byte) string {
	return base64.StdEncoding.EncodeToString(buf)
}

// B64Decode base64 decodes a string
func B64Decode(str string) (buf []byte, err error) {
	return base64.StdEncoding.DecodeString(str)
}

// StrContained returns true if 'str' is in 'strs'; otherwise return false
func StrContained(str string, strs []string) bool {
	for _, s := range strs {
		if strings.ToLower(s) == strings.ToLower(str) {
			return true
		}
	}
	return false
}

// IsSubsetOf returns an error if there is something in 'small' that
// is not in 'big'.  Both small and big are assumed to be comma-separated
// strings.  All string comparisons are case-insensitive.
// Examples:
// 1) IsSubsetOf('a,B', 'A,B,C') returns nil
// 2) IsSubsetOf('A,B,C', 'B,C') returns an error because A is not in the 2nd set.
func IsSubsetOf(small, big string) error {
	bigSet := strings.Split(big, ",")
	smallSet := strings.Split(small, ",")
	for _, s := range smallSet {
		if s != "" && !StrContained(s, bigSet) {
			return errors.Errorf("'%s' is not a member of '%s'", s, big)
		}
	}
	return nil
}

// HTTPRequestToString returns a string for an HTTP request for debuggging
func HTTPRequestToString(req *http.Request) string {
	body, _ := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	return fmt.Sprintf("%s %s\n%s",
		req.Method, req.URL, string(body))
}

// HTTPResponseToString returns a string for an HTTP response for debuggging
func HTTPResponseToString(resp *http.Response) string {
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	return fmt.Sprintf("statusCode=%d (%s)\n%s",
		resp.StatusCode, resp.Status, string(body))
}

// GetDefaultConfigFile gets the default path for the config file to display in usage message
func GetDefaultConfigFile(cmdName string) string {
	if cmdName == "credential-provider-server" {
		var fname = fmt.Sprintf("%s-config.yaml", cmdName)
		// First check home env variables
		home := "."
		envs := []string{"CREDENTIAL_PROVIDER_SERVER_HOME", "CREDENTIAL_PROVIDER_HOME", "CREDENTIAL_PROVIDER_CFG_PATH"}
		for _, env := range envs {
			envVal := os.Getenv(env)
			if envVal != "" {
				home = envVal
				break
			}
		}
		return path.Join(home, fname)
	}

	var fname = fmt.Sprintf("%s-config.yaml", cmdName)

	return path.Join(os.Getenv("HOME"), ".credential-provider-server", fname)
}

// MakeFileAbs makes 'file' absolute relative to 'dir' if not already absolute
func MakeFileAbs(file, dir string) (string, error) {
	if file == "" {
		return "", nil
	}
	if filepath.IsAbs(file) {
		return file, nil
	}
	path, err := filepath.Abs(filepath.Join(dir, file))
	if err != nil {
		return "", errors.Wrapf(err, "Failed making '%s' absolute based on '%s'", file, dir)
	}
	return path, nil
}

// MakeFileNamesAbsolute makes all file names in the list absolute, relative to home
func MakeFileNamesAbsolute(files []*string, home string) error {
	for _, filePtr := range files {
		abs, err := MakeFileAbs(*filePtr, home)
		if err != nil {
			return err
		}
		*filePtr = abs
	}
	return nil
}

// Fatal logs a fatal message and exits
func Fatal(format string, v ...interface{}) {
	log.Fatalf(format, v...)
	os.Exit(1)
}

// GetSerialAsHex returns the serial number from certificate as hex format
func GetSerialAsHex(serial *big.Int) string {
	hex := fmt.Sprintf("%x", serial)
	return hex
}

// StructToString converts a struct to a string. If a field
// has a 'secret' tag, it is masked in the returned string
func StructToString(si interface{}) string {
	rval := reflect.ValueOf(si).Elem()
	tipe := rval.Type()
	var buffer bytes.Buffer
	buffer.WriteString("{ ")
	for i := 0; i < rval.NumField(); i++ {
		tf := tipe.Field(i)
		if !rval.FieldByName(tf.Name).CanSet() {
			continue // skip unexported fields
		}
		var fStr string
		tagv := tf.Tag.Get(SecretTag)
		if tagv == "password" || tagv == "username" {
			fStr = fmt.Sprintf("%s:**** ", tf.Name)
		} else if tagv == "url" {
			val, ok := rval.Field(i).Interface().(string)
			if ok {
				val = GetMaskedURL(val)
				fStr = fmt.Sprintf("%s:%v ", tf.Name, val)
			} else {
				fStr = fmt.Sprintf("%s:%v ", tf.Name, rval.Field(i).Interface())
			}
		} else {
			fStr = fmt.Sprintf("%s:%v ", tf.Name, rval.Field(i).Interface())
		}
		buffer.WriteString(fStr)
	}
	buffer.WriteString(" }")
	return buffer.String()
}

// GetMaskedURL returns masked URL. It masks username and password from the URL
// if present
func GetMaskedURL(url string) string {
	matches := URLRegex.FindStringSubmatch(url)

	// If there is a match, there should be four entries: 1 for
	// the match and 3 for submatches
	if len(matches) == 4 {
		matchIdxs := URLRegex.FindStringSubmatchIndex(url)
		matchStr := url[matchIdxs[0]:matchIdxs[1]]
		for idx := 2; idx < len(matches); idx++ {
			if matches[idx] != "" {
				matchStr = strings.Replace(matchStr, matches[idx], "****", 1)
			}
		}
		url = url[:matchIdxs[0]] + matchStr + url[matchIdxs[1]:len(url)]
	}
	return url
}

// NormalizeStringSlice checks for seperators
func NormalizeStringSlice(slice []string) []string {
	var normalizedSlice []string

	if len(slice) > 0 {
		for _, item := range slice {
			// Remove surrounding brackets "[]" if specified
			if strings.HasPrefix(item, "[") && strings.HasSuffix(item, "]") {
				item = item[1 : len(item)-1]
			}
			// Split elements based on comma and add to normalized slice
			if strings.Contains(item, ",") {
				normalizedSlice = append(normalizedSlice, strings.Split(item, ",")...)
			} else {
				normalizedSlice = append(normalizedSlice, item)
			}
		}
	}

	return normalizedSlice
}

// NormalizeFileList provides absolute pathing for the list of files
func NormalizeFileList(files []string, homeDir string) ([]string, error) {
	var err error

	files = NormalizeStringSlice(files)

	for i, file := range files {
		files[i], err = MakeFileAbs(file, homeDir)
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

// Read reads from Reader into a byte array
func Read(r io.Reader, data []byte) ([]byte, error) {
	j := 0
	for {
		n, err := r.Read(data[j:])
		j = j + n
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrapf(err, "Read failure")
		}

		if (n == 0 && j == len(data)) || j > len(data) {
			return nil, errors.New("Size of requested data is too large")
		}
	}

	return data[:j], nil
}

// Hostname name returns the hostname of the machine
func Hostname() string {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "localhost"
	}
	return hostname
}

// ValidateAndReturnAbsConf checks to see that there are no conflicts between the
// configuration file path and home directory. If no conflicts, returns back the absolute
// path for the configuration file and home directory.
func ValidateAndReturnAbsConf(configFilePath, homeDir, cmdName string) (string, string, error) {
	var err error
	var homeDirSet bool
	var configFileSet bool

	defaultConfig := GetDefaultConfigFile(cmdName) // Get the default configuration

	if configFilePath == "" {
		configFilePath = defaultConfig // If no config file path specified, use the default configuration file
	} else {
		configFileSet = true
	}

	if homeDir == "" {
		homeDir = filepath.Dir(defaultConfig) // If no home directory specified, use the default directory
	} else {
		homeDirSet = true
	}

	// Make the home directory absolute
	homeDir, err = filepath.Abs(homeDir)
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to get full path of config file")
	}
	homeDir = strings.TrimRight(homeDir, "/")

	if configFileSet && homeDirSet {
		log.Warning("Using both --config and --home CLI flags; --config will take precedence")
	}

	if configFileSet {
		configFilePath, err = filepath.Abs(configFilePath)
		if err != nil {
			return "", "", errors.Wrap(err, "Failed to get full path of configuration file")
		}
		return configFilePath, filepath.Dir(configFilePath), nil
	}

	configFile := filepath.Join(homeDir, filepath.Base(defaultConfig)) // Join specified home directory with default config file name
	return configFile, homeDir, nil
}

// FatalError will check to see if an error occured if so it will cause the test cases exit
func FatalError(t *testing.T, err error, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args)
	}
	if !assert.NoError(t, err, msg) {
		t.Fatal(msg)
	}
}

// ErrorContains will check to see if an error occurred, if so it will check that it contains
// the appropriate error message
func ErrorContains(t *testing.T, err error, contains, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args)
	}
	if assert.Error(t, err, msg) {
		assert.Contains(t, err.Error(), contains)
	}
}

// GetSliceFromList will return a slice from a list
func GetSliceFromList(split string, delim string) []string {
	return strings.Split(strings.Replace(split, " ", "", -1), delim)
}

// ListContains looks through a comma separated list to see if a string exists
func ListContains(list, find string) bool {
	items := strings.Split(list, ",")
	for _, item := range items {
		item = strings.TrimPrefix(item, " ")
		if item == find {
			return true
		}
	}
	return false
}
