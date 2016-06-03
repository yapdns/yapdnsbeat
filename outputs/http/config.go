package http

type config struct {
	ApiEndpoint     string `config:"api_endpoint"`
	BulkApiEndpoint string `config:"bulk_api_endpoint"`
}

var (
	defaultConfig = config{}
)
