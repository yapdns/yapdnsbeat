package http

type config struct {
	ClientId        string `config:"client_id"`
	ClientSecretKey string `config:"client_secret_key"`
	ApiEndpoint     string `config:"api_endpoint"`
	BulkApiEndpoint string `config:"bulk_api_endpoint"`
}

var (
	defaultConfig = config{}
)
