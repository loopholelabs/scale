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

package log

import (
	"fmt"
	"io"
)

var _ io.Writer = (*NamedLogger)(nil)

type NamedLogger struct {
	nameTemplate      string
	nameTemplateBytes []byte
	writer            io.Writer
}

func NewNamedLogger(name string, writer io.Writer) *NamedLogger {
	nameTemplate := fmt.Sprintf("%s: ", name)
	return &NamedLogger{
		nameTemplate:      nameTemplate,
		nameTemplateBytes: []byte(nameTemplate),
		writer:            writer,
	}
}

func (l *NamedLogger) Write(p []byte) (int, error) {
	_, err := l.writer.Write(append(append(l.nameTemplateBytes, p...), '\n'))
	return len(p), err
}
