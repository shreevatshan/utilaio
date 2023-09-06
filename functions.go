package utilaio

import (
	"archive/zip"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func RemoveString(str *string, remove string) {
	*str = strings.ReplaceAll(*str, remove, EMPTY_STRING)
}

func ReplaceString(str *string, from string, to string) {
	*str = strings.ReplaceAll(*str, from, to)
}

func RemoveWhiteSpace(str *string) {
	*str = strings.ReplaceAll(*str, SPACE, EMPTY_STRING)
}

func ReturnKeyAndValueFromString(keyvalue_string string) (key, value string) {

	key = EMPTY_STRING
	value = EMPTY_STRING

	last_index := strings.LastIndex(keyvalue_string, EQUAL_TO)
	if last_index != -1 {
		key = keyvalue_string[:last_index]
		value = keyvalue_string[last_index+1:]
	}
	return key, value
}

func ConvertNewlineSeparatedStringToKeyValuePairBasedOnEqual(newline_separated_string string) map[string]string {

	result_map := make(map[string]string)

	if newline_separated_string != EMPTY_STRING {
		individual_lines := strings.Split(newline_separated_string, NEW_LINE)
		for i := range individual_lines {
			individual_line := individual_lines[i]
			key, value := ReturnKeyAndValueFromString(individual_line)
			RemoveWhiteSpace(&key)
			RemoveWhiteSpace(&value)
			result_map[key] = value
		}
	}
	return result_map
}

func ConvertNewlineSeparatedStringToKeyValuePairBasedOnEqualAndComma(newline_separated_string string) map[string]string {

	result_map := make(map[string]string)

	if newline_separated_string != EMPTY_STRING {
		individual_lines := strings.Split(newline_separated_string, NEW_LINE)
		for i := range individual_lines {
			individual_line := individual_lines[i]
			keys, value := ReturnKeyAndValueFromString(individual_line)
			RemoveWhiteSpace(&keys)
			RemoveWhiteSpace(&value)
			individual_keys := strings.Split(keys, COMMA)
			for j := range individual_keys {
				key := individual_keys[j]
				result_map[key] = value
			}
		}
	}
	return result_map
}

func DecodeBase64Encoded(base64encoded []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(base64encoded))
}

func GetCurrentTimeInUnixMilli() int64 {
	current_time := time.Now()
	return current_time.UnixMilli()
}

func HasPrefixCaseInsensitive(string_to_check string, string_to_compare string) bool {
	return strings.HasPrefix(strings.ToLower(string_to_check), strings.ToLower(string_to_compare))
}

func FormKeyValuePairFromMapOfValueInterface(from_map map[string]interface{}) []KeyValuePair {
	var to_array []KeyValuePair
	for key, value := range from_map {
		var pair KeyValuePair
		pair.Key = key
		pair.Value = value
		to_array = append(to_array, pair)
	}
	return to_array
}

func FormKeyValuePairFromMapOfValueInterfaceSlice(from_map map[string][]interface{}) []KeyValuePair {
	var to_array []KeyValuePair
	for key, value := range from_map {
		var pair KeyValuePair
		pair.Key = key
		pair.Value = value
		to_array = append(to_array, pair)
	}
	return to_array
}

func GetValueFromInterfaceMapAsString(key string, interface_map map[string]interface{}) string {
	var value string
	if _, exists := interface_map[key]; exists {
		value = interface_map[key].(string)
	}
	return value
}

func ConvertIntegerToBoolean(integer int) bool {
	return integer > 0
}

func Unzip(source string, destination string) error {

	zipReader, _ := zip.OpenReader(source)
	for _, file := range zipReader.Reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zippedFile.Close()

		extractedFilePath := filepath.Join(
			destination,
			file.Name,
		)

		if file.FileInfo().IsDir() {
			os.MkdirAll(extractedFilePath, file.Mode())
		} else {
			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				return err
			}
			defer outputFile.Close()

			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func MapToHashCode(value interface{}) ([]byte, error) {

	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(data)

	return hash[:], nil
}

func WriteCurrentPID(file_path string) error {
	pid := os.Getpid()

	pidBytes := []byte(fmt.Sprintf("%d", pid))

	err := ioutil.WriteFile(file_path, pidBytes, 0777)

	return err
}

func CopyFile(from_location string, to_location string) error {
	var err error

	original, err := os.Open(from_location)
	if err != nil {
		return err
	}

	new, err := os.Create(to_location)
	if err != nil {
		original.Close()
		return err
	}

	_, err = io.Copy(new, original)

	original.Close()
	new.Close()
	return err
}

func CurrentUsername() (string, error) {

	var username string

	current_user, err := user.Current()
	if err == nil {
		username = current_user.Username
	}

	return username, err
}

func RoundTo(n float64, decimals uint32) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}

func Float64ComparisonWithTolerance(a float64, b float64, tolerance float64) bool {
	if diff := math.Abs(a - b); diff < tolerance {
		return true
	} else {
		return false
	}
}

// Enter len as a multiple of 2
func CreateUUID(len int) (string, error) {
	bytes := make([]byte, len/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	uid := hex.EncodeToString(bytes)
	return uid, nil
}
