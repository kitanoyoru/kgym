package prometheus

var GlobalRegistry *Registry

func init() {
	GlobalRegistry = NewRegistry()
}
