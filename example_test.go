package hterrors_test

import (
	"fmt"
	"net/http"

	"snai.pe/go-hterrors"
)

func ExampleCheck() {
	resp, err := hterrors.Check(http.Get("http://google.com"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(http.StatusText(resp.StatusCode))
	// Output: OK
}

func ExampleCheck_notfound() {
	resp, err := hterrors.Check(http.Get("http://google.com/invalid"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(http.StatusText(resp.StatusCode))
	// Output: GET "http://google.com/invalid": Error 404 (Not Found)!!1 : 404. That’s an error. : The requested URL /invalid was not found on this server. That’s all we know.
}
