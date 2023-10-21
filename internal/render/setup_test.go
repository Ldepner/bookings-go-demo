package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/ldepner/bookings/internal/config"
	"github.com/ldepner/bookings/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// what am I going to store in my session?
	gob.Register(models.Reservation{})

	testApp.InProduction = false

	session = scs.New()
	testApp.Session = session
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	app = &testApp

	os.Exit(m.Run())
}
