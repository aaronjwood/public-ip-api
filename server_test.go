package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestStartServer(t *testing.T) {
	go main()
	url := fmt.Sprintf("http://localhost:%d", port)
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err.Error())
	}

	if res.StatusCode != http.StatusOK {
		t.Fatal("Should have responded with a 200")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	res.Body.Close()
	if string(data) != "::1" {
		t.Fatal("Wrong response")
	}

	res, err = http.Get(url + "/ip")
	if err != nil {
		t.Fatal(err.Error())
	}

	if res.StatusCode != http.StatusNotFound {
		t.Fatal("Should have responded with a 404")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	ip := "192.168.1.0"
	req.Header.Set("X-Forwarded-For", ip)
	res, err = client.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	res.Body.Close()
	if strings.Trim(string(data), "\n") != ip {
		t.Fatal("Wrong IP returned")
	}
}
