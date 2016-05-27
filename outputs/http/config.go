package http

type config struct {
	Pretty bool `config:"pretty"`
	Foo bool `config:"foo"`
}

var (
	defaultConfig = config{
		Pretty: false,
		Foo: true,
	}
)
