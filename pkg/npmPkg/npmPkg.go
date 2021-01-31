package npmPkg

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strings"
)

type Pkg struct {
	Name			string				`json:"name"`
	Dependencies	map[string]string	`json:"dependencies"`
}

const endpoint = "https://registry.npmjs.org"

func GetDependencies(name, ver string) (Pkg, error) {
	url := []string{endpoint, name, ver}
	resp, err := http.Get(strings.Join(url, "/"))
	if err != nil {
		return Pkg{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Pkg{}, err
	}
	var p Pkg
	if err := json.Unmarshal(body, &p); err != nil {
		return Pkg{}, err
	}
	return p, nil
}
