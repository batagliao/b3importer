package data

import (
	"b3importer/configuration"
	"time"

	"github.com/rs/zerolog/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	BATCH_SIZE = 100
)

type dbContext struct {
	db *gorm.DB
}

var instance *dbContext

func loadInstance() *dbContext {
	log.Debug().Msg("connecting to database...")
	dsn := configuration.Get().MysqlAddress
	dialector := mysql.Open(dsn)
	db, err := gorm.Open(dialector, &gorm.Config{})
	db.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
	}
	if err != nil {
		log.Fatal().Err(err)
	}

	instance = &dbContext{
		db: db,
	}

	// log.Println("realizing migrations...")
	// instance.db.AutoMigrate(
	// 	&balance.BalanceActiveIndividual{},
	// 	&balance.BalanceActiveConsolidated{},
	// )

	return instance
}

func Instance() *dbContext {
	if instance == nil {
		instance = loadInstance()
	}
	return instance
}

func (d *dbContext) Create(values interface{}) error {
	log.Debug().Msg("inserting data")
	if err := d.db.CreateInBatches(values, BATCH_SIZE).Error; err != nil {
		return err
	}
	return nil
}

func (d *dbContext) ClearCompanyOldData(entity interface{}, cnpj string, dtrefer time.Time) error {
	log.Debug().Msgf("removing data based on cnpj %s and dt_refer %v", cnpj, dtrefer)
	if err := d.db.Unscoped().Where("CNPJ_Cia LIKE ? AND Data_Referencia = ?", cnpj, dtrefer).Delete(entity).Error; err != nil {
		return err
	}
	return nil
}

func (d *dbContext) Migrate(objs ...interface{}) error {
	return d.db.AutoMigrate(objs...)
}
