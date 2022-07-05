package commons

const (
	ConnectionTimeout   = 90 // seconds
	DatabaseName        = "exercise"
	MigrationFolderPath = "file://server/repositories/db/migrations"
	NatsUrl             = "localhost:4222"
	//NatsUrl				= "nats:4222"
	DatabaseHost = "localhost"
	//DatabaseHost				= "postgres"
)

type RepositoryType string

const (
	RepositoryTypeDB       RepositoryType = "DB"
	RepositoryTypeFile     RepositoryType = "File"
	RepositoryTypeInMemory RepositoryType = "InMemory"
)
