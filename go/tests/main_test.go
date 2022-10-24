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
	"github.com/loopholelabs/scale/go/scalefunc"
	"os"
	"os/exec"
	"testing"
)

type TestCase struct {
	Name   string
	Module string
	Run    func(scalefunc.ScaleFunc, *testing.T)
}

func TestMain(m *testing.M) {
	err := exec.Command("sh", "compile.sh").Run()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	err = exec.Command("sh", "cleanup.sh").Run()
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}
