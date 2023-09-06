/*
	Copyright 2022 Loophole Labs

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

package scale

import "strings"

type Parsed struct {
	Organization string
	Name         string
	Tag          string
}

// Parse parses a function or signature name of the form <org>/<name>:<tag> into its organization, name, and tag
func Parse(name string) *Parsed {
	orgSplit := strings.Split(name, "/")
	if len(orgSplit) == 1 {
		orgSplit = []string{"", name}
	}
	tagSplit := strings.Split(orgSplit[1], ":")
	if len(tagSplit) == 1 {
		tagSplit = []string{tagSplit[0], ""}
	}
	return &Parsed{
		Organization: orgSplit[0],
		Name:         tagSplit[0],
		Tag:          tagSplit[1],
	}
}

func unpackUint32(packed uint64) (uint32, uint32) {
	return uint32(packed >> 32), uint32(packed)
}
