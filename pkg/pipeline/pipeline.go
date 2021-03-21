package pipeline

type Pipeline func() (next bool, err error)

func Run(pipelines ...Pipeline) error {
	for _, p := range pipelines {
		next, err := p()
		if next {
			continue
		}

		if err != nil {
			return err
		}

		if !next {
			return nil
		}
	}

	return nil
}
