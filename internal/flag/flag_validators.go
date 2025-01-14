package flag

type flagValidators []func(Flag) error

const (
	ErrFlagApplyDebugConflict = FlagErr("Apply and debug flags cannot be set to true at the same time")
	ErrFlagPeriodLessThanOne  = FlagErr("Period flag (-p|--period) cannot be negative")
)

func (f Flag) validateFlags() error {

	flagValidatorList := flagValidators{
		applyFlagValidator,
		periodFlagValidator,
	}

	for _, validator := range flagValidatorList {
		err := validator(f)

		if err != nil {
			return err
		}
	}

	return nil
}

func applyFlagValidator(f Flag) error {

	if f.Apply && f.Debug {

		return ErrFlagApplyDebugConflict
	}

	return nil
}

func periodFlagValidator(f Flag) error {

	if f.Period < 1 {
		return ErrFlagPeriodLessThanOne
	}

	return nil
}
