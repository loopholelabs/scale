//go:build !integration

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

package format

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

const data = `
use crate::{ lazy::{Lazy, SyncLazy, SyncOnceCell}, panic,
        sync::{ atomic::{AtomicUsize, Ordering::SeqCst},
            mpsc::channel, Mutex, },
      thread,
    };
    impl<T, U> Into<U> for T where U: From<T> {
        fn into(self) -> U { U::from(self) }
    }
`

func TestFormatter(t *testing.T) {
	f, err := New()
	require.NoError(t, err)

	formatted, err := f.Format(context.Background(), data)
	require.NoError(t, err)

	t.Logf("Formatted:\n\n%s\n", formatted)
}
