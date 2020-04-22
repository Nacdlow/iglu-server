package models

import (
	_ "github.com/mattn/go-sqlite3" // SQLite driver support

	"log"
	"time"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/Nacdlow/iglu-server/modules/settings"
	"xorm.io/core"
	"xorm.io/xorm"
)

var (
	engine     *xorm.Engine
	tables     []interface{}
	sqlitePath = "data.db?cache=shared&mode=ro&_busy_timeout=9999999"
)

func init() {
	tables = append(tables,
		new(Alert),
		new(Device),
		new(Room),
		new(RoomStat),
		new(Schedule),
		new(Statistic),
		new(User),
	)
	gofakeit.Seed(0) // Using 0 will use current Unix time.
}

// SetupEngine sets up an XORM engine and syncs the schema.
// It will return an xorm engine.
func SetupEngine() *xorm.Engine {
	var err error
	engine, err = xorm.NewEngine("sqlite3", sqlitePath)
	if err != nil {
		log.Fatalln("Failed to setup engine! ", err)
	}

	engine.SetMapper(core.GonicMapper{})
	err = engine.Sync(tables...)

	if err != nil {
		log.Fatalln("Failed to sync schema! ", err)
	}

	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 2000)
	engine.SetDefaultCacher(cacher)
	location, err := time.LoadLocation(settings.Config.GetString("Timezone"))
	if err != nil {
		log.Fatalln("Failed to load timezone! ", err)
	}
	engine.TZLocation = location

	return engine
}

// SetupTestEngine sets up an XORM engine for unit testing and syncs the
// schema to it.
func SetupTestEngine() *xorm.Engine {
	var err error
	engine, err = xorm.NewEngine("sqlite3", ":memory:")
	if err != nil {
		log.Fatalln("Failed to setup test engine! ", err)
	}

	engine.SetMapper(core.GonicMapper{})
	err = engine.Sync(tables...)

	if err != nil {
		log.Fatalln("Failed to sync schema with testing database! ", err)
	}

	return engine
}
