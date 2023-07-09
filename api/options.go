package api

// Option -.
type Option func(*ByBit)

// BaseUrl -
func BaseUrl(url string) Option {
	return func(bb *ByBit) {
		bb.baseURL = url
	}
}

// DebugMode -.
func DebugMode(mode bool) Option {
	return func(bb *ByBit) {
		bb.debugMode = mode
	}
}
