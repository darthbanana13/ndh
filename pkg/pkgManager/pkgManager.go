package pkgManager

import (
	"fmt"
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
		fmt.Println(curPkg.Name)
		p, _ := npmPkg.GetDependencies(curPkg.Name, curPkg.Version)
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
				Version: val,
			}
			unprocessed = append(unprocessed, pkg)
		}
		bigMap[name] = val
	}
	return unprocessed
}
