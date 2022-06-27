package commons

const (
	ConnectionTimeout   = 90 // seconds
	DatabaseName        = "GoExercise"
	MigrationFolderPath = "file://server/repositories/db/migrations"
)

type RepositoryType string

const (
	RepositoryTypeDB       RepositoryType = "DB"
	RepositoryTypeFile     RepositoryType = "File"
	RepositoryTypeInMemory RepositoryType = "InMemory"
)
