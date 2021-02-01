package pkgManager

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/darthrevan13/ndh/pkg/npmPkg"
)

type PkgTreeNode struct {
	Name         string
	Version      string
	Dependencies []*PkgTreeNode
}

type TreeNodeJson struct {
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	Dependencies []TreeNodeJson `json:"dependencies"`
}

func GetAllDependencies(name, ver string) (PkgTreeNode, error) {
	topNode := PkgTreeNode{
		Name:    name,
		Version: santizeVersion(ver),
	}
	//map[packageName][packageVersion]packagePointer
	flatDependencies := map[string]map[string]*PkgTreeNode{
		name: map[string]*PkgTreeNode{
			ver: &topNode,
		},
	}
	//TODO: Convert to channel with a buffer of 1
	pkgsToProcess := []*PkgTreeNode{&topNode}

	for i := 0; i < len(pkgsToProcess); i++ {
		curPkg := pkgsToProcess[i]
		// TODO: Move to subrutine with waitGroup
		p, err := npmPkg.GetDependencies(curPkg.Name, curPkg.Version)
		if err != nil {
			//TODO: Make errors non blocking, put errors on a separate channel and fail subrutine, not the entire process
			return PkgTreeNode{}, err
		}
		node := convertPkgToTreeNode(p)
		flatDependencies[curPkg.Name][curPkg.Version].Dependencies = node.Dependencies
		unprocessed := findUnprocessedDependencies(flatDependencies, node.Dependencies)
		// TODO: Send values to channel instead
		pkgsToProcess = append(pkgsToProcess, unprocessed...)
	}
	return topNode, nil
}

func findUnprocessedDependencies(bigMap map[string]map[string]*PkgTreeNode, smallMap []*PkgTreeNode) []*PkgTreeNode {
	var unprocessed []*PkgTreeNode
	for i := 0; i < len(smallMap); i++ {
		curValue := smallMap[i]
		if _, ok := bigMap[curValue.Name]; !ok {
			unprocessed = append(unprocessed, curValue)
			bigMap[curValue.Name] = map[string]*PkgTreeNode{curValue.Version: curValue}
		} else if _, ok := bigMap[curValue.Name][curValue.Version]; !ok {
			unprocessed = append(unprocessed, curValue)
			bigMap[curValue.Name][curValue.Version] = curValue
		}
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
		if len(matches) > 0 {
			return matches[0]
		}
		//TODO: Handle regex no match
	}
	//TODO: Handle other less common cases
	return ver
}

func convertPkgToTreeNode(pkg npmPkg.Pkg) PkgTreeNode {
	topNode := PkgTreeNode{
		Name:    pkg.Name,
		Version: santizeVersion(pkg.Version),
	}
	for name, ver := range pkg.Dependencies {
		topNode.Dependencies = append(topNode.Dependencies, &PkgTreeNode{
			Name:    name,
			Version: santizeVersion(ver),
		})
	}
	return topNode
}

func (n PkgTreeNode) ToJson() (string, error) {
	displayNode := n.baseConvertToTreeNodeJson()
	displayNode.Dependencies = dereferenceAsTreeNodeJson(n.Dependencies)
	return displayNode.ToPrettyJson()
}

func (t TreeNodeJson) ToPrettyJson() (string, error) {
	jsonByte, err := json.MarshalIndent(t, "", "    ")
	return string(jsonByte), err
}

func (n PkgTreeNode) baseConvertToTreeNodeJson() TreeNodeJson {
	return TreeNodeJson{
		Name:    n.Name,
		Version: n.Version,
	}
}

// TODO: Convert function to non-recursive call
func dereferenceAsTreeNodeJson(slice []*PkgTreeNode) []TreeNodeJson {
	var result []TreeNodeJson
	for _, value := range slice {
		treeNode := value.baseConvertToTreeNodeJson()
		treeNode.Dependencies = dereferenceAsTreeNodeJson(value.Dependencies)
		result = append(result, treeNode)
	}
	return result
}
