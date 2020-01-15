package models

import (
	_ "github.com/mattn/go-sqlite3" // SQLite driver support
	"log"
	"xorm.io/core"
	"xorm.io/xorm"
)

var (
	engine *xorm.Engine
	tables []interface{}
)

func init() {
	tables = append(tables,
		new(Device),
		new(Room),
		new(RoomStat),
		new(Statistic),
		new(User),
	)
}

// SetupEngine sets up an XORM engine and syncs the schema.
// It will return an xorm engine.
func SetupEngine() *xorm.Engine {
	engine, err := xorm.NewEngine("sqlite3", "data.db")
	if err != nil {
		log.Fatalln("Failed to setup engine!", err)
	}

	engine.SetMapper(core.GonicMapper{})
	err = engine.Sync(tables...)

	if err != nil {
		log.Fatalln("Failed to sync schema!", err)
	}

	return engine
}
