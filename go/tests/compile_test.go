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

package tests

import (
	"bytes"
	"github.com/loopholelabs/scale/go/scalefile"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"text/template"
)

func TestCompileDefaultTemplate(t *testing.T) {
	tmplFile, err := os.ReadFile("../compile/go.mod.tmpl")
	assert.NoError(t, err)
	tmpl, err := template.New("dependencies").Parse(string(tmplFile))
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)

	err = tmpl.Execute(buf, scalefile.DefaultDependencies)
	assert.NoError(t, err)

	assert.Equal(t, "module github.com/loopholelabs/scale/go/compile\n\ngo 1.18\n\nrequire github.com/loopholelabs/scale/go v0.0.8\n\n", buf.String())
}

func TestCompileCustomTemplate(t *testing.T) {
	tmplFile, err := os.ReadFile("../compile/go.mod.tmpl")
	assert.NoError(t, err)
	tmpl, err := template.New("dependencies").Parse(string(tmplFile))
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)

	err = tmpl.Execute(buf, []scalefile.Dependency{
		{
			Name:    "github.com/loopholelabs/scale/go",
			Version: "v0.0.1",
		},
		{
			Name:    "github.com/loopholelabs/scale/go/context",
			Version: "v0.0.1",
		},
		{
			Name:    "github.com/loopholelabs/scale/go/scalefunc",
			Version: "v0.0.2",
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, "module github.com/loopholelabs/scale/go/compile\n\ngo 1.18\n\nrequire github.com/loopholelabs/scale/go v0.0.1\n\nrequire github.com/loopholelabs/scale/go/context v0.0.1\n\nrequire github.com/loopholelabs/scale/go/scalefunc v0.0.2\n\n", buf.String())
}
