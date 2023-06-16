/*
	Copyright 2023 Loophole Labs
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at
		   http://www.apache.org/licenses/LICENSE-2.0
	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package utils

import (
	"errors"
	"unicode"
)

// Params is a function that creates a map of parameters for use in go templates to pass multiple values
func Params(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("parameters must be a list of key/value pairs")
	}
	params := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("keys must be strings")
		}
		params[key] = values[i+1]
	}
	return params, nil
}

// CamelCase converts a PascalCase string to camelCase
func CamelCase(s string) string {
	runes := []rune(s)
	// Lowercase consecutive uppercase characters at the beginning of the string
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i]) && unicode.IsUpper(runes[i+1]) {
			runes[i] = unicode.ToLower(runes[i])
		} else {
			break
		}
	}
	// Convert to camelCase
	for i := 0; i < len(runes); i++ {
		if unicode.IsUpper(runes[i]) {
			if i > 0 {
				return string(runes[:i]) + string(unicode.ToLower(runes[i])) + string(runes[i+1:])
			}
			return string(unicode.ToLower(runes[i])) + string(runes[i+1:])
		}
	}
	return string(runes)
}
