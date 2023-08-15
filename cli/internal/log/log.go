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
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"io"
	"sync"
)

var (
	once   sync.Once
	Logger = NewLogger(io.Discard)
)

// init sets up the time format and an error marshaller that lets us record an error's stack trace
func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

// NewLogger creates a new zerolog.Logger with default values
func NewLogger(w io.Writer) *zerolog.Logger {
	l := zerolog.New(w).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	return &l
}

func Init(logFile string) {
	once.Do(func() {
		writer := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    128,
			MaxAge:     7,
			MaxBackups: 4,
		}

		Logger = NewLogger(writer)
		Logger.Info().Msg("logger initialized")
	})
}
