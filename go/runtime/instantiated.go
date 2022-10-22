package runtime

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tetratelabs/wazero/api"
)

type Instantiated struct {
	m      *Module
	module api.Module
	run    api.Function
	resize api.Function
}

func NewInstantiated(ctx context.Context, i *Instance, m *Module) (*Instantiated, error) {
	module, err := i.runtime.runtime.InstantiateModule(ctx, m.function.Compiled, i.runtime.moduleConfig.WithName(fmt.Sprintf("%s.%s.%s", i.id, m.function.ScaleFunc.ScaleFile.Name, uuid.New().String())))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate function '%s' for instance %s: %w", m.function.ScaleFunc.ScaleFile.Name, i.id, err)
	}

	run := module.ExportedFunction("run")
	resize := module.ExportedFunction("resize")
	if run == nil || resize == nil {
		return nil, fmt.Errorf("failed to find run or resize functions for instance %s", i.id)
	}

	return &Instantiated{
		m:      m,
		module: module,
		run:    run,
		resize: resize,
	}, nil
}
