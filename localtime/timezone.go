package localtime

import (
	"time"
	_ "time/tzdata"

	"github.com/rs/zerolog/log"
)

// Sets time.Location, which is used in the localtime package
func SetLocalTimeZone(location string) {
	time.Local = MustLoadLocation(location)
}

func MustLoadLocation(location string) *time.Location {
	loc, err := time.LoadLocation(location)
	if err != nil {
		log.Panic().Err(err).Msgf("Kunne ikke laste tidssone %s", location)
	}
	return loc
}
