package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testEnvironment(t *testing.T, host string, path string) string {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal("Failed to create request")
	}

	req.Header.Add("Host", host)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(onRequest)

	handler.ServeHTTP(rr, req)

	status := rr.Code
	bodyLen := rr.Body.Len()

	if status != http.StatusTemporaryRedirect {
		t.Errorf("Bad status code %d", status)
	}
	if bodyLen != 0 {
		t.Errorf("Bad body length %d", bodyLen)
	}

	return rr.Header().Get("Location")
}

func assertTest(
	t *testing.T,
	host string,
	path string,
	expected string,
	) {
	actual := testEnvironment(t, host, path)
	if actual != expected {
		t.Errorf("Expected to get sent to %s actually got %s", expected, actual)
	}
}

// https://www.jrtapsell.co.uk|{p}|{q}
func TestBasic(t *testing.T) {
	assertTest(
		t,
		"jrtapsell.co.uk",
		"/",
		"https://www.jrtapsell.co.uk/",
	)
}

// https://www.jrtapsell.co.uk|{p}|{q}
func TestPath(t *testing.T) {
	assertTest(
		t,
		"jrtapsell.co.uk",
		"/test",
		"https://www.jrtapsell.co.uk/test",
	)
}

// https://monzo.me/jamesrichardtapsell|{p}
func TestPay(t *testing.T) {
	assertTest(
		t,
		"pay.jrtapsell.co.uk",
		"/12.34",
		"https://monzo.me/jamesrichardtapsell/12.34",
	)
}

// https://www.jrtapsell.co.uk|{p}|{q}
func TestQuery(t *testing.T) {
	assertTest(
		t,
		"jrtapsell.co.uk",
		"/test?a=a&b=b",
		"https://www.jrtapsell.co.uk/test?a=a&b=b",
	)
}