package pkgManager

import (
	"fmt"
	"strings"
	"regexp"

	"github.com/darthrevan13/ndh/pkg/npmPkg"
)

type pkgToProcess struct {
	Name	string
	Version string
}

func GetAllDependencies(name, ver string) (map[string]string, error) {
	flatDependencies := map[string]string{name: ver}
	pkg := pkgToProcess{
		Name: name,
		Version: ver,
	}
	pkgsToProcess := []pkgToProcess{pkg}

	for i := 0; i < len(pkgsToProcess); i++ {
		curPkg := pkgsToProcess[i]
		p, err := npmPkg.GetDependencies(curPkg.Name, curPkg.Version)
		if err != nil {
			//TODO: Make errors non blocking, handle in separate chanel with concurrency
			return map[string]string{}, err
		}
		unprocessed := refreshDependencies(flatDependencies, p.Dependencies)
		pkgsToProcess = append(pkgsToProcess, unprocessed...)
	}
	return flatDependencies, nil
}

func refreshDependencies(bigMap, smallMap map[string]string) []pkgToProcess {
	var unprocessed []pkgToProcess
	for name, val := range smallMap {
		if _, ok := bigMap[name]; !ok {
			pkg := pkgToProcess{
				Name: name,
				Version: santizeVersion(val),
			}
			unprocessed = append(unprocessed, pkg)
		}
		bigMap[name] = val
	}
	return unprocessed
}

func santizeVersion(ver string) string {
	if strings.HasPrefix(ver, "~") || strings.HasPrefix(ver, "^") {
		return strings.TrimLeft(ver, "~^")
	} else if ver == "*" {
		return "latest"
	} else if strings.HasPrefix(ver, ">=") {
		reg := regexp.MustCompile(`[0-9\.]+`)
		matches := reg.FindStringSubmatch(ver)
		if len(matches) > 1 {
			return matches[1]
		}
		//TODO: Handle regex no match
	}
	//TODO: Handle other cases
	return ver
}
