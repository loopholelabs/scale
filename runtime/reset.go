package runtime

func (r *Runtime) Reset() error {
	for _, f := range r.functions {
		err := f.Module.Close(r.ctx)
		if err != nil {
			return err
		}
		f.Module, err = r.runtime.InstantiateModule(r.ctx, f.Compiled, r.moduleConfig.WithName(f.ScaleFunc.ScaleFile.Name))
		if err != nil {
			return err
		}
	}
	return nil
}
