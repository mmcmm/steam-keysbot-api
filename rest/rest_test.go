package rest

import (
	"database/sql"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/db"
)

var ts *httptest.Server
var body, expected string

func TestMain(m *testing.M) {
	dbconn := db.Open()
	defer dbconn.Close()
	r := StartRouter(dbconn)
	ts = httptest.NewServer(r)
	defer ts.Close()

	code := m.Run()
	os.Exit(code)
}

func TestAccount(t *testing.T) {
	t.Parallel()

	_, body = common.TestRequest(t, ts, "GET", "/api/v1/account", nil)
	expected = `{"elapsed":0}`
	if strings.Compare(strings.TrimSpace(body), expected) != 0 {
		t.Fatalf("expected:%s got:%s", expected, body)
	}
}

func loginTestAccount(dbconn *sql.DB) {
	dbconn.Exec(`INSERT INTO users (steam_id, display_name, avatar) VALUES ($1, $2, $3) 
	ON CONFLICT (steam_id) DO UPDATE SET display_name = $2, avatar = $3`,
		00000000000000000, "TestPersonaNamae", "TestAvatar")

}
