package tst

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/mgutz/ansi"
)

func Request(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return resp.StatusCode, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return resp.StatusCode, ""
	}
	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

func AssertNoError(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

func AssertError(tb testing.TB, err error) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: expecting error but got nil\033[39m\n\n", filepath.Base(file), line)
		tb.FailNow()
	}
}

func AssertEqual(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		yellow := ansi.ColorFunc("yellow+h")
		green := ansi.ColorCode("green+h:black")
		red := ansi.ColorCode("red+h:black")
		reset := ansi.ColorCode("reset")

		_, file, line, _ := runtime.Caller(1)

		fmt.Println(yellow(fmt.Sprintf("%s:%d", filepath.Base(file), line)))
		fmt.Println(green, "Expected: ", exp, reset)
		fmt.Println(red, "     Got: ", act, reset)

		tb.FailNow()
	}
}

// NotEquals fails the test if exp is not equal to act.
func AssertNotEqual(tb testing.TB, exp, act interface{}) {
	if reflect.DeepEqual(exp, act) {
		yellow := ansi.ColorFunc("yellow+h")
		green := ansi.ColorCode("green+h:black")
		red := ansi.ColorCode("red+h:black")
		reset := ansi.ColorCode("reset")

		_, file, line, _ := runtime.Caller(1)

		fmt.Println(yellow(fmt.Sprintf("%s:%d", filepath.Base(file), line)))
		fmt.Println(green, "Expected: ", exp, reset)
		fmt.Println(red, "Got: ", act, reset)

		tb.FailNow()
	}
}