package prettyslog_test

import (
	"testing"

	prettyslog "github.com/rautaruukkipalich/prettyslog"
)

func TestPrettyLogger(t *testing.T) {
	log := prettyslog.NewPrettyLogger()
	log.Debug("Debug message")
}