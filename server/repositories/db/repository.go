package db

import (
	generatedModel "Exercise/graph/model"
	"Exercise/server/commons"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var saveQuery = fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1,$2,$3,$4,$5) RETURNING id;", repositoryTableName, insertFields)

//var updateQuery = fmt.Sprintf("UPDATE %s SET status = $1, modified = CURRENT_TIMESTAMP, modified_by = $2 WHERE service_tag = $3 RETURNING %s;", repositoryTableName, returnFields)
var updateQuery = fmt.Sprintf("UPDATE %s SET status = $1, modified = CURRENT_TIMESTAMP, modified_by = $2 WHERE service_tag = $3;", repositoryTableName)
var getQuery = fmt.Sprintf("SELECT %s FROM %s;", returnFields, repositoryTableName)
var getByQuery = fmt.Sprintf("SELECT %s FROM %s WHERE", returnFields, repositoryTableName)

const (
	repositoryTableName = "dataentries"
	insertFields        = " id, title, \"content\", \"views\", \"timestamp\""
	returnFields        = insertFields
)

type databaseRepository struct {
	con *sql.DB
}

func NewRepository() (*databaseRepository, error) {
	ctx := context.Background()
	config := &Config{}
	*config = parseEnv()
	con, err := NewDatabase(*config)
	if err != nil {
		log.Panic(ctx, err, "unable to setup database")
		return nil, errors.New("unable to setup database")
	}

	err = Migrate(con)
	if err != nil {
		log.Panic(ctx, err, "unable to setup database")
		return nil, errors.New("Failed while migration")
	}
	return &databaseRepository{con}, nil
}

func (repo *databaseRepository) Save(ctx context.Context, model generatedModel.InputDataEntry) (data interface{}, err error) {
	var id string
	repo.con.QueryRow(saveQuery, model.ID, model.Title, model.Content, model.Views, model.Timestamp).Scan(&id)
	//TODO: update below print statement with database persist functionality
	//log.Panic(ctx, "type cast model to eoSkeleton %v", model)
	return nil, err
}

func (repo *databaseRepository) Get(ctx context.Context) (data []*generatedModel.DataEntry, err error) {

	var dataEntries []*generatedModel.DataEntry

	rows, e := repo.con.Query(getQuery)

	if e != nil {
		// handle this error better than this
		panic(e)
	}
	defer rows.Close()
	for rows.Next() {
		var retId string
		var title string
		var content string
		var views int
		var timestamp string
		err = rows.Scan(&retId, &title, &content, &views, &timestamp)
		if err != nil {
			// handle this error
			panic(err)
		}
		dataEntry := generatedModel.DataEntry{
			ID:        retId,
			Title:     title,
			Content:   content,
			Views:     views,
			Timestamp: timestamp,
		}
		dataEntries = append(dataEntries, &dataEntry)

	}

	return dataEntries, e
}

func (repo *databaseRepository) GetById(ctx context.Context, query string) (data *generatedModel.DataEntry, err error) {
	var retId string
	var title string
	var content string
	var views int
	var timestamp string

	e := repo.con.QueryRow(fmt.Sprintf("%s %s;", getByQuery, query)).Scan(&retId, &title, &content, &views, &timestamp)

	if e != nil && e != sql.ErrNoRows {
		// handle this error better than this
		panic(e)
	}

	dataEntry := generatedModel.DataEntry{
		ID:        retId,
		Title:     title,
		Content:   content,
		Views:     views,
		Timestamp: timestamp,
	}

	return &dataEntry, e
}

func (repo *databaseRepository) Update(ctx context.Context, model interface{}) (data interface{}, err error) {
	//TODO: update below print statement with database persist functionality
	return nil, err
}

func parseEnv() Config {
	return Config{
		Host:     "10.43.93.252",
		Port:     "5432",
		User:     "postgres",
		Password: "123456",
		Database: commons.DatabaseName,
	}
}
