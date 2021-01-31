package pkgManager

import (
	"fmt"
	"github.com/darthrevan13/ndh/pkg/npmPkg"
)

func GetAllDependencies(name, ver string) (map[string]string, error) {
	flatDependencies := map[string]string{name: ver}
	for n, v := range flatDependencies {
		fmt.Println(n)
		p, _ := npmPkg.GetDependencies(n, v)
		addMap(flatDependencies, p.Dependencies)
	}
	return flatDependencies, nil
}

func addMap(bigMap, smallMap map[string]string) {
	for k, v := range smallMap {
		bigMap[k] = v
	}
}
