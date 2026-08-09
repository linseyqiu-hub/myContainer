package errorHandlers

func RunSteps(steps ...func() error) error {
	for _, step := range steps {

		err := step()
		if err != nil {
			return err
		}

	}

	return nil
}
