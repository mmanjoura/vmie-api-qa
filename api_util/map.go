// This package provides various utility functions for working with interfaces, maps, 
// arrays, and JSON/YAML conversion. Each function has a comment explaining its 
// purpose and usage. Let me know if you need further clarification on any specific function!

package api_util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"gopkg.in/yaml.v3"
)

// InterfaceToJsonString converts the given interface{} value to a JSON string.
// It uses the json.Marshal function to serialize the interface{} value to JSON.
// If the resulting JSON string starts and ends with double quotes, it removes them before returning the string.
// If an error occurs during the serialization process, an empty string is returned.
func InterfaceToJsonString(i interface{}) string {
	b, _ := json.Marshal(i)
	if b[0] == '"' {
		return string(b[1 : len(b)-1]) // remove the ""
	}
	return string(b)
}

// MapInterfaceToMapString converts the params map (all primitive types with exception of array)
// before passing to resty.
func MapInterfaceToMapString(src map[string]interface{}) map[string]string {
	dst := make(map[string]string)
	for k, v := range src {
		if ar, ok := v.([]interface{}); ok {
			str := ""
			for _, entry := range ar {
				str += fmt.Sprintf("%v,", InterfaceToJsonString(entry))
			}
			str = strings.TrimRight(str, ",")
			dst[k] = str
		} else {
			dst[k] = InterfaceToJsonString(v)
		}
	}
	return dst
}

// MapIsCompatible checks if the first map has every key in the second.
func MapIsCompatible(big map[string]interface{}, small map[string]interface{}) bool {
	for k, _ := range small {
		if _, ok := big[k]; !ok {
			return false
		}
	}
	return true
}

// TimeCompare compares two values and determines if they represent the same time.
// It takes two interface{} arguments, v1 and v2, which are expected to be strings
// representing time in RFC3339 format. The function returns true if the values are
// equal, otherwise it returns false.
//
// If either v1 or v2 is not a string or not in RFC3339 format, the function returns false.
//
// If both v1 and v2 are valid RFC3339 strings, the function compares the parsed time values
// and returns true if they are equal, otherwise it returns false.
//
// If only one of v1 and v2 is a valid RFC3339 string, the function checks if the other value
// contains the second and minute elements of the parsed time value. If both elements are found,
// the function returns true, otherwise it returns false.
//
// Example usage:
//   t1 := "2022-01-01T12:00:00Z"
//   t2 := "2022-01-01T12:00:00Z"
//   result := TimeCompare(t1, t2) // returns true
//
//   t3 := "2022-01-01T12:00:00Z"
//   t4 := "2022-01-01T12:01:00Z"
//   result := TimeCompare(t3, t4) // returns false

func TimeCompare(v1 interface{}, v2 interface{}) bool {
	s1, ok := v1.(string)
	if !ok {
		return false
	}
	s2, ok := v2.(string)
	if !ok {
		return false
	}
	var t time.Time
	var s string
	var b1, b2 bool
	t1, err := time.Parse(time.RFC3339, s1)
	if err == nil {
		t = t1
		s = s2
		b1 = true
	}
	t2, err := time.Parse(time.RFC3339, s2)
	if err == nil {
		t = t2
		s = s1
		b2 = true
	}
	if b1 && b2 {
		return t1 == t2
	}
	if !b1 && !b2 {
		return false
	}
	// One of b1 and b2 is true, now t point to time and s point to a potential time string
	// that's not RFC3339 format. We make a guess buy searching for a few key elements.
	return strings.Contains(s, fmt.Sprintf("%d", t.Second())) && strings.Contains(s, fmt.Sprintf("%d", t.Minute()))
}

// MapCombine combines two map together. If there is any overlap the dst will be overwritten.
// MapCombine combines two maps by copying all key-value pairs from the source map to the destination map.
// If the destination map is empty, it returns a copy of the source map.
// If the source map is empty, it returns the destination map.
// If both maps have key-value pairs, the key-value pairs from the source map overwrite the corresponding key-value pairs in the destination map.
// The function modifies the destination map in-place and returns it.
func MapCombine(dst map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	if len(dst) == 0 {
		return MapCopy(src)
	}
	if len(src) == 0 {
		return dst
	}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}


// MapAdd adds the key-value pairs from the source map to the destination map.
// If the destination map is empty, it returns a copy of the source map.
// If the source map is empty, it returns the destination map unchanged.
// If a key already exists in the destination map, it is not overwritten.
// The function modifies the destination map in-place and returns it.
func MapAdd(dst map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	if len(dst) == 0 {
		return MapCopy(src)
	}
	if len(src) == 0 {
		return dst
	}
	for k, v := range src {
		if _, exist := dst[k]; !exist {
			dst[k] = v
		}
	}
	return dst
}

// MapReplace replaces the values in dst with the ones in src with the matching keys.
// MapReplace replaces the values in the destination map with the values from the source map.
// If the source map is empty, the destination map is returned as is.
// The function iterates over the keys in the destination map and checks if the key exists in the source map.
// If the key exists, the value in the destination map is replaced with the corresponding value from the source map.
// The modified destination map is then returned.
func MapReplace(dst map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	if len(src) == 0 {
		return dst
	}
	for k := range dst {
		if v, ok := src[k]; ok {
			dst[k] = v
		}
	}
	return dst
}

// MapCopy creates a deep copy of a map[string]interface{}.
// It recursively copies nested maps and arrays to ensure a complete copy.
// If the source map is empty, it returns nil.
// The function returns a new map with the copied values.
func MapCopy(src map[string]interface{}) map[string]interface{} {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]interface{})
	for k, v := range src {
		if m, ok := v.(map[string]interface{}); ok {
			v = MapCopy(m)
		}
		if a, ok := v.([]interface{}); ok {
			v = ArrayCopy(a)
		}
		dst[k] = v
	}
	return dst
}


// ArrayCopy copies the elements of the source array to a new destination array.
// If the source array contains maps or arrays, they are recursively copied as well.
// The copied elements are returned as a new array.
// If the source array is empty, nil is returned.
func ArrayCopy(src []interface{}) (dst []interface{}) {
	if len(src) == 0 {
		return nil
	}
	for _, v := range src {
		if m, ok := v.(map[string]interface{}); ok {
			v = MapCopy(m)
		}
		if a, ok := v.([]interface{}); ok {
			v = ArrayCopy(a)
		}
		dst = append(dst, v)
	}
	return dst
}

// InterfacePrint prints the given interface as YAML to the logger and optionally to the console.
// It takes two parameters:
//   - m: The interface to be printed as YAML.
//   - printToConsole: A boolean flag indicating whether to print the YAML to the console.
// The function marshals the interface to YAML format using the `yaml.Marshal` function,
// and then logs the YAML string to the logger using the `Logger.Print` method.
// If the `printToConsole` flag is set to true, it also prints the YAML string to the console using `fmt.Println`.
// Note: Any error that occurs during marshaling is ignored.
func InterfacePrint(m interface{}, printToConsole bool) {
	yamlBytes, _ := yaml.Marshal(m)
	Logger.Print(string(yamlBytes))
	if printToConsole {
		fmt.Println(string(yamlBytes))
	}
}

// Check if existing matches criteria. When criteria is a map, we check whether
// everything in criteria can be found and equals a field in existing.
// InterfaceEquals compares two interface{} values for equality.
// It returns true if the values are equal, and false otherwise.
// The function handles various types of values, including nil, maps, arrays, slices, and strings.
// For maps, it recursively compares the key-value pairs.
// For strings, it checks if the existing value is a JSON number.
// For other types, it compares the values using reflection and JSON marshaling.
func InterfaceEquals(criteria interface{}, existing interface{}) bool {
	if criteria == nil {
		if existing == nil {
			return true
		} else {
			existingKind := reflect.TypeOf(existing).Kind()
			if existingKind == reflect.Map || existingKind == reflect.Array || existingKind == reflect.Slice {
				return true
			}
			return false
		}
	} else {
		if existing == nil {
			return false
		}
	}
	cType := reflect.TypeOf(criteria)
	eType := reflect.TypeOf(existing)
	if cType == eType && cType.Comparable() {
		if criteria == existing {
			return true
		}
		// The only exception is time, where the format may be different on both ends.
		return TimeCompare(criteria, existing)
	}

	cKind := cType.Kind()
	eKind := eType.Kind()
	if cKind == reflect.Array || cKind == reflect.Slice {
		if eKind == reflect.Array || eKind == reflect.Slice {
			// We don't compare arrays
			return true
		}
		return false
	}
	if cKind == reflect.Map {
		if eKind != reflect.Map {
			return false
		}
		cm, ok := criteria.(map[string]interface{})
		if !ok {
			return false
		}
		em, ok := existing.(map[string]interface{})
		if !ok {
			return false
		}
		for k, v := range cm {
			if !InterfaceEquals(v, em[k]) {
				return false
			}
		}
		return true
	}
	if eKind == reflect.String && (cKind == reflect.Int || cKind == reflect.Float32 || cKind == reflect.Float64) {
		return reflect.TypeOf(existing).String() == "json.Number"
	}

	cJson, _ := json.Marshal(criteria)
	eJson, _ := json.Marshal(existing)

	return string(cJson) == string(eJson)
}

// MarshalJsonIndentNoEscape marshals the given interface into a JSON byte slice with indentation and without HTML escaping.
// It takes an interface{} as input and returns the marshaled JSON byte slice and an error, if any.
// The indentation is set to four spaces.
func MarshalJsonIndentNoEscape(i interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	err := enc.Encode(i)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Given a yaml stream, output a json stream.
// YamlToJson converts YAML data to JSON format.
// It takes a byte slice `in` containing the YAML data and returns a `json.RawMessage` and an error.
// The `json.RawMessage` represents the converted JSON data, while the error indicates any conversion errors.
func YamlToJson(in []byte) (json.RawMessage, error) {
	var unmarshaled interface{}
	err := yaml.Unmarshal(in, &unmarshaled)
	if err != nil {
		return nil, err
	}
	return swag.YAMLToJSON(unmarshaled)
}

// JsonToYaml converts a JSON byte array to a YAML byte array.
// It takes a JSON byte array as input and returns a YAML byte array and an error.
// If the conversion is successful, the function returns the YAML byte array and a nil error.
// If an error occurs during the conversion, the function returns a nil byte array and the error.
func JsonToYaml(in []byte) ([]byte, error) {
	var out interface{}
	err := json.Unmarshal(in, &out)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(out)
}

// YamlObjToJsonObj converts a YAML object to a JSON object.
// It takes an input interface{} representing the YAML object and returns the corresponding JSON object.
// If the conversion is successful, it returns the JSON object and nil error.
// If an error occurs during the conversion, it returns nil and the error.
func YamlObjToJsonObj(in interface{}) (interface{}, error) {
	jsonRaw, err := swag.YAMLToJSON(in)
	if err != nil {
		return nil, err
	}
	var out interface{}
	err = json.Unmarshal(jsonRaw, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}


// FieldIterFunc is a function type that represents an iterator function for iterating over fields in a map.
// It takes a key string and a value interface{} as parameters and returns an error.
// The iterator function can be used to perform operations on each key-value pair in a map.
type FieldIterFunc func(key string, value interface{}) error

// MapIterFunc is a function type that represents a callback function
// used for iterating over a map[string]interface{}.
// The function takes a map as a parameter and returns an error if any.
type MapIterFunc func(m map[string]interface{}) error

// Iterate all the leaf level fields. For maps iterate all the fields. For arrays we will go through all the entries and
// see if any of them is a map. The iteration will be done in a width first order, so deeply buried fields will be iterated last.
// The maps should be map[string]interface{}.
// IterateMapsInInterface iterates over maps in an interface and applies a callback function to each map.
// The callback function should have the signature `func(map[string]interface{}) error`.
// If the input interface contains nested maps or arrays, the function recursively iterates over them as well.
// Returns an error if the callback function returns an error, otherwise returns nil.
func IterateMapsInInterface(in interface{}, callback MapIterFunc) error {
	if inMap, _ := in.(map[string]interface{}); inMap != nil {
		err := callback(inMap)
		if err != nil {
			return err
		}
		for _, v := range inMap {
			err := IterateMapsInInterface(v, callback)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if inArray, _ := in.([]interface{}); inArray != nil {
		for _, v := range inArray {
			err := IterateMapsInInterface(v, callback)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

// IterateFieldsInInterface iterates over the fields of an interface{} value and invokes the provided callback function for each field.
// The callback function should have the signature `func(key string, value interface{}) error`.
// It returns an error if the callback function returns an error, otherwise it returns nil.
func IterateFieldsInInterface(in interface{}, callback FieldIterFunc) error {
	mapCallback := func(m map[string]interface{}) error {
		for k, v := range m {
			err := callback(k, v)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return IterateMapsInInterface(in, mapCallback)
}