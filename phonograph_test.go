package phonograph_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/whytheplatypus/phonograph"
)

func Test_Phonograph(t *testing.T) {
	p := "testdata"
	phonograph.Record(p)
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("test"))
	}))

	resp, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)

	phonograph.Play(p)

	rec, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rec)

	_, err = http.Get("test")
	if err == nil {
		t.Fatal("no test to find")
	}
}
