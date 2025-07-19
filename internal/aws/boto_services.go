package aws

type BotoServiceDataSources map[string]BotoDataSource

// BotoDataSource for loading each service data schema
type BotoDataSource struct {
	ApiVersion string
	Filename   string
	Sha        string
}
