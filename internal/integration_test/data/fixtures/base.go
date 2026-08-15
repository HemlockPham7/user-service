package fixtures

import (
	"testing"

	"github.com/HemlockPham7/common-libs/pkg/sqldb"
	"gorm.io/gorm"
)

type Fixture interface {
	SetupDB(db *gorm.DB)
	Migrate() error
	GenerateData() error
	DB() *gorm.DB
}

type base struct {
	db *gorm.DB
}

func (b *base) SetupDB(db *gorm.DB) {
	b.db = db
}

func (b *base) DB() *gorm.DB {
	return b.db
}

func NewFixture(t *testing.T, fix Fixture) *gorm.DB {
	// create test db
	fix.SetupDB(sqldb.InitMockDB(t))

	// migrate schema
	err := fix.Migrate()
	if err != nil {
		t.Fatalf("failed to migrate fixture: %s", err)
	}

	// generate data
	err = fix.GenerateData()
	if err != nil {
		t.Fatalf("failed to generate fixture data: %s", err)
	}

	// return db
	return fix.DB()
}
