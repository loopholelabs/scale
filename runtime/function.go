package runtime

import (
	"context"
	"fmt"
	"github.com/loopholelabs/scale-go/scalefunc"
	"github.com/tetratelabs/wazero"
)

// Function is the runtime representation of a scale function.
type Function struct {
	ScaleFunc scalefunc.ScaleFunc
	Compiled  wazero.CompiledModule
}

func (r *Runtime) registerFunction(ctx context.Context, scaleFunc scalefunc.ScaleFunc) error {
	compiled, err := r.runtime.CompileModule(ctx, scaleFunc.Function, r.compileConfig)
	if err != nil {
		return fmt.Errorf("failed to compile function '%s': %w", scaleFunc.ScaleFile.Name, err)
	}

	f := &Function{
		ScaleFunc: scaleFunc,
		Compiled:  compiled,
	}

	r.functions = append(r.functions, f)
	return nil
}
