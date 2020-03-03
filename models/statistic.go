package models

import (
	"errors"
	"time"

	"github.com/brianvoe/gofakeit/v4"
)

// Statistic represents a statistic log at a period of time, which spans an
// hour.
type Statistic struct {
	StatID               int64   `xorm:"pk autoincr" json:"id" xml:"id,attr"`
	StatTime             int64   `xorm:"unique index" json:"statTime" xml:"timestamps>stat_time"`
	PowerGenAvg          float64 `json:"powerGenAvg" xml:"avg>power_gen,attr"`               // Power generated, kWh
	PowerConAvg          float64 `json:"powerConAvg" xml:"avg>power_con,attr"`               // Power conserved, kWh
	MainDoorsOpenedCount int64   `json:"mainDoorsOpenedCount" xml:"main_doors_opened_count"` // How many times the doors opened
	CreatedUnix          int64   `xorm:"created" json:"createdUnix" xml:"timestamps>created_unix"`
	UpdatedUnix          int64   `xorm:"updated" json:"updatedUnix" xml:"timestamps>updated_unix"`
}

// GetFakeStat returns a new randomly generated statistic. This is used for
// testing purposes.
func GetFakeStat() (s *Statistic) {
	s = new(Statistic)
	s.CreatedUnix = time.Now().UnixNano() - int64(gofakeit.Number(0, 99999))
	s.PowerGenAvg = gofakeit.Float64Range(0, 45)
	s.PowerConAvg = gofakeit.Float64Range(0, s.PowerGenAvg)
	return
}

// GetLatestStats returns the latest stats (from the past 24 hours).
func GetLatestStats() (s []Statistic, err error) {
	minTime := time.Now().Add(-(24 * time.Hour))
	var stats []Statistic
	stats, err = GetStats()
	if err != nil {
		return
	}
	for _, stat := range stats {
		if stat.StatTime > minTime.Unix() {
			s = append(s, stat)
		}
	}
	return
}

// StatExists checks whether a statistic exists based on statistic time.
func StatExists(time int64) bool {
	total, _ := engine.Where("stat_time = ?", time).Count(new(Statistic))
	return total > 0
}

// GetStat gets a Statistic based on its ID from the database.
func GetStat(id int64) (*Statistic, error) {
	s := new(Statistic)
	has, err := engine.ID(id).Get(s)
	if err != nil {
		return s, err
	} else if !has {
		return s, errors.New("Statistic does not exist")
	}
	return s, nil
}

// GetStats returns an array of all Statistics from the database.
func GetStats() (stat []Statistic, err error) {
	err = engine.Find(&stat)
	return
}

// AddStat adds a Statistic in the database.
func AddStat(s *Statistic) (err error) {
	_, err = engine.Insert(s)
	return
}

// HasStat returns whether an Statistic is in the database or not.
func HasStat(id int64) (has bool) {
	has, _ = engine.Get(&Statistic{StatID: id})
	return
}

// UpdateStat updates an Statistic in the database.
func UpdateStat(d *Statistic) (err error) {
	_, err = engine.Id(d.StatID).Update(d)
	return
}

// UpdateStatCols will update the columns of an item even if they are empty.
func UpdateStatCols(d *Statistic, cols ...string) (err error) {
	_, err = engine.ID(d.StatID).Cols(cols...).Update(d)
	return
}

// DeleteStat deletes a Statistic from the database.
func DeleteStat(id int64) (err error) {
	_, err = engine.ID(id).Delete(&Statistic{})
	return
}
