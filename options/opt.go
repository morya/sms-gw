package options

type Options struct {
	Account    string
	Password   string
	RemoteAddr string

	ConcurrentCount int

	AutoSuicide bool
}
