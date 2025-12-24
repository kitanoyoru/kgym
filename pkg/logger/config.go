package logger

type option func(*config)

type config struct {
	dev bool
}

func WithDev() option {
	return func(c *config) {
		c.dev = true
	}
}
