package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

//var theTests_old []struct{
//	name       string
//	url        string
//	method     string
//	params     []postData
//	statusCode int
//}{
//{"home", "/", "GET", []postData{}, http.StatusOK},
//}

var theTests = []struct {
	name       string
	url        string
	method     string
	params     []postData
	StatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
}

func testHandlers(t *testing.T) {

	routes := getRoutes()
	testServ := httptest.NewTLSServer(routes)

	defer testServ.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := testServ.Client().Get(testServ.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.StatusCode {
				t.Errorf("for %s, expected %d, got %d", e.name, e.StatusCode, resp.StatusCode)
			}
		} else {

		}
	}
}
