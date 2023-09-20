//go:build !integration && !generate

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

package typescript

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

const testData = `
{
  "name": "semver",
  "version": "7.5.4",
  "description": "The semantic version parser used by npm.",
  "main": "index.js",
  "scripts": {
    "test": "tap",
    "snap": "tap",
    "lint": "eslint \"**/*.js\"",
    "postlint": "template-oss-check",
    "lintfix": "npm run lint -- --fix",
    "posttest": "npm run lint",
    "template-oss-apply": "template-oss-apply --force"
  },
  "devDependencies": {
    "@npmcli/eslint-config": "^4.0.0",
    "@npmcli/template-oss": "4.18.0",
    "tap": "^16.0.0"
  },
  "license": "ISC",
  "repository": {
    "type": "git",
    "url": "https://github.com/npm/node-semver.git"
  },
  "bin": {
    "semver": "bin/semver.js"
  },
  "files": [
    "bin/",
    "lib/",
    "classes/",
    "functions/",
    "internal/",
    "ranges/",
    "index.js",
    "preload.js",
    "range.bnf"
  ],
  "tap": {
    "timeout": 30,
    "coverage-map": "map.js",
    "nyc-arg": [
      "--exclude",
      "tap-snapshots/**"
    ]
  },
  "engines": {
    "node": ">=10"
  },
  "dependencies": {
    "lru-cache": "^6.0.0"
  },
  "author": "GitHub Inc.",
  "templateOSS": {
    "//@npmcli/template-oss": "This file is partially managed by @npmcli/template-oss. Edits may be overwritten.",
    "version": "4.18.0",
    "engines": ">=10",
    "ciVersions": [
      "10.0.0",
      "10.x",
      "12.x",
      "14.x",
      "16.x",
      "18.x"
    ],
    "npmSpec": "8",
    "distPaths": [
      "classes/",
      "functions/",
      "internal/",
      "ranges/",
      "index.js",
      "preload.js",
      "range.bnf"
    ],
    "allowPaths": [
      "/classes/",
      "/functions/",
      "/internal/",
      "/ranges/",
      "/index.js",
      "/preload.js",
      "/range.bnf"
    ],
    "publish": "true"
  }
}`

func TestPackage(t *testing.T) {
	m, err := ParseManifest([]byte(testData))
	require.NoError(t, err)

	data, err := m.Write()
	require.NoError(t, err)

	testMap := make(map[string]interface{})
	dataMap := make(map[string]interface{})

	err = json.Unmarshal([]byte(testData), &testMap)
	require.NoError(t, err)

	err = json.Unmarshal(data, &dataMap)
	require.NoError(t, err)

	require.Equal(t, testMap, dataMap)
}
