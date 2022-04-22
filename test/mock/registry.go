// build mock

// The global registry for mocked object automatically assigned by mock constructors with '+build mock' tags.
package mock

var (
	registry = make(map[string]interface{})
)

// RegisterMockObject registers mock object into global registry, if object is singleton name could be a name of
// structure, in case of many object of one type name must be constructed as primary index value.
func RegisterMockObject(name string, mock interface{}) {
	registry[name] = mock
}

// GetMockObject receives registered mock object from registry.
func GetMockObject(name string) interface{} {
	return registry[name]
}

const (
	SolanaProvider    = "SolanaProvider"
	PostMarkProvider  = "PostMarkProvider"
	CoingeckoProvider = "CoingeckoProvider"
)
