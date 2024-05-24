package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"payment/pkg/config"
	"strings"
)

type DB struct {
	*gorm.DB
}

// LoadDSN constructs a Data Source Name (DSN) string for connecting to the database
// using the provided configuration settings.
func LoadDSN(config *config.Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		config.PostgresConfig.Host,
		config.PostgresConfig.User,
		config.PostgresConfig.Password,
		config.PostgresConfig.Database,
		config.PostgresConfig.Port,
		config.PostgresConfig.Timezone)

}

// New initializes a new database connection using the provided DSN string.
// It returns a pointer to a DB instance or an error if the connection fails.
func New(dsn string) (*DB, error) {
	dialector := postgres.Open(dsn)
	db, err := gorm.Open(dialector, &gorm.Config{FullSaveAssociations: false})
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// AutoMigrate runs the GORM AutoMigrate function for the provided models,
// ensuring that the database schema matches the model definitions.
func AutoMigrate(db *DB, model ...interface{}) error {
	err := db.AutoMigrate(model...)
	if err != nil {
		return err
	}
	return nil
}

func CreateEnumType(db *DB, typeName string, values []string) error {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = ?);"
	err := db.Raw(query, typeName).Scan(&exists).Error
	if err != nil {
		return err
	}

	if !exists {
		valuesList := "'" + strings.Join(values, "', '") + "'"
		createTypeQuery := fmt.Sprintf("CREATE TYPE %s AS ENUM (%s);", typeName, valuesList)
		if err := db.Exec(createTypeQuery).Error; err != nil {
			return err
		}
		log.Printf("type %s created with values: %v", typeName, values)
		return nil
	}

	log.Printf("type %s already exists", typeName)
	return nil
}
