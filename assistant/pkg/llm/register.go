package llm

import "fmt"

type Factory func(cfg Config) (Client, error)

var providers = map[string]Factory{}

func Register(name string, f Factory) {
	providers[name] = f
}

func NewClient(provider string, cfg Config) (Client, error) {
	f, ok := providers[provider]
	if !ok {
		return nil, fmt.Errorf("llm: unknown provider %q", provider)
	}
	return f(cfg)
}

func GetProvider(name string) (Factory, bool) {
	f, ok := providers[name]
	return f, ok
}

func RegisteredProviders() []string {
	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	return names
}
