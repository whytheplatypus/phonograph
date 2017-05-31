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

type Cylinder struct {
	Parent http.RoundTripper
	Path   string
}

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

type Crank struct {
	Path string
}

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

func Record(path string) {
	http.DefaultClient.Transport = &Cylinder{
		Path:   path,
		Parent: http.DefaultTransport,
	}
}

func Play(path string) {
	http.DefaultClient.Transport = &Crank{
		Path: path,
	}
}
