package rest

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRest_Ping(t *testing.T) {
	ts, _, teardown := startup()
	defer teardown()

	response, code := get(t, ts.URL+"/api/health")
	assert.Equal(t, "", response)
	assert.Equal(t, 200, code)
}

func get(t *testing.T, url string) (string, int) {
	r, err := http.Get(url)
	require.NoError(t, err)

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	require.NoError(t, err)
	return string(body), r.StatusCode
}

func startup() (ts *httptest.Server, server *RestServer, teardown func()) {
	server = &RestServer{
		authenticator: &AuthService{authenticator: &TestAuthenticator{}},
	}

	ts = httptest.NewServer(server.routes())

	teardown = func() {
		ts.Close()
	}

	return ts, server, teardown
}

type TestAuthenticator struct{}

func (a *TestAuthenticator) verify(token string) error {
	return nil
}
