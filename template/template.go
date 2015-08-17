package template

type Variables struct {
	Variables map[string]string
	Raw       []byte
	Path      string
}

type Config struct {
	Path string
	Raw  []byte

	Imports []struct{ Path string }
}
