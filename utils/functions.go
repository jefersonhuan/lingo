package utils

func MapWithErrors(functions ...func() error) error {
	for _, f := range functions {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}
