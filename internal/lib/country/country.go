package country

import (
	"time"
)

const (
	TjID     = 1
	TjName   = "Tajikistan"
	TjPrefix = "tj"
)

const (
	_utc5         = "+5"
	_utc5Duration = 5 * time.Hour
)

type Country struct {
	Shard     string
	Name      string
	UTC       string
	ToUTC     time.Duration
	FromUTC   time.Duration
	CountryID int8
}

var Countries = [3]Country{{CountryID: TjID, Shard: TjPrefix, Name: TjName, UTC: _utc5, ToUTC: -_utc5Duration, FromUTC: _utc5Duration}}

func (c Country) GetID() int8 {
	return c.CountryID
}

func (c Country) Prefix() string {
	return c.Shard + "_"
}

func (c Country) Postfix() string {
	return "_" + c.Shard
}

func ByID(id int8) Country {
	for _, country := range Countries {
		if id == country.CountryID {
			return country
		}
	}
	return Countries[0]
}

func ConvertTimeToUTC(t time.Time, countryID int8) time.Time {
	return t.Add(ByID(countryID).ToUTC)
}

func ConvertTimeFromUTC(t time.Time, countryID int8) time.Time {
	return t.Add(ByID(countryID).FromUTC)
}
