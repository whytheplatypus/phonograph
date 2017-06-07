package phonograph

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
)

// Cylinder provides a [RoundTripper](https://golang.org/pkg/net/http/#RoundTripper)
// that records interactions sent through it's `Parent`.
type Cylinder struct {
	Parent http.RoundTripper
	Path   string
}

// RoundTrip records each response in a file
// named with the md5 hash of it's request.
func (c *Cylinder) RoundTrip(req *http.Request) (*http.Response, error) {
	rb, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	h := fmt.Sprintf("%x", md5.Sum(rb))

	resp, err := c.Parent.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	respb, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(filepath.Join(c.Path, string(h)), respb, os.ModePerm); err != nil {
		return nil, err
	}

	return resp, err
}

// Crank provides a [RoundTripper](https://golang.org/pkg/net/http/#RoundTripper)
// that plays back the matching recorded response in `Path`
// for the incoming request.
type Crank struct {
	Path string
}

// RoundTrip matches the md5 hash for a request to a recorded response.
// If a file in `Path` mataches the hash the response recorded within
// is returned.
func (c *Crank) RoundTrip(req *http.Request) (*http.Response, error) {
	rb, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	h := fmt.Sprintf("%x", md5.Sum(rb))
	respb, err := ioutil.ReadFile(filepath.Join(c.Path, string(h)))

	if err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(bytes.NewBuffer(respb)), req)
}

// Record sets the http.DefaultClient transport
// to a Cylinder wrapping http.DefaultTransport.
func Record(path string) {
	http.DefaultClient.Transport = &Cylinder{
		Path:   path,
		Parent: http.DefaultTransport,
	}
}

// Play sets the http.DefaultClient transport
// to a Crank wrapping http.DefaultTransport.
func Play(path string) {
	http.DefaultClient.Transport = &Crank{
		Path: path,
	}
}
