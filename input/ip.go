package input

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var Ip string

func init() {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	Ip = strings.TrimSpace(string(data))
}
