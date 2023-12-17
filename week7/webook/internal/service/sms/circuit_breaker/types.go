package cb

type CircuitBreaker interface {
	Do(args ...any) error
}
