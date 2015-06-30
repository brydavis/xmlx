package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	// "strings"
)

/* TODO
- Transform xml "Attr" to nested elements
- Generalize API for other formats
- Transform data to actual data types
- Integrate into "xql" package
*/

type Node struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
	Nodes   []Node `xml:",any"`
	// Attr    xml.Attr
}

func main() {
	n := Nodify("xml/complex.xml")
	n.Write("json/complex.json", 0700)
	// n.Pretty(0)

	j, _ := Pretty("json/complex.json")
	fmt.Printf("%s\n", j)

}

func Nodify(filename string) Node {
	data, _ := ioutil.ReadFile(filename)
	buf := bytes.NewBuffer(data)
	dec := xml.NewDecoder(buf)

	var n Node
	dec.Decode(&n)

	return n
}

func (n Node) Walk() interface{} {
	if len(n.Nodes) < 1 {
		return string(n.Content)
	} else {
		x := make(map[string][]interface{})
		y := make(map[string]interface{})

		for _, v := range n.Nodes {
			x[v.XMLName.Local] = append(x[v.XMLName.Local], v.Walk())
		}

		for k, v := range x {
			if len(x[k]) == 1 {
				switch v[0].(type) {
				case float32, float64:
					f, _ := strconv.ParseFloat(v[0].(string), 64)
					y[k] = f // v[0].(float64)
				case bool:
					y[k] = v[0].(bool)
				default:
					y[k] = v[0]
				}
			} else {
				y[k] = v
			}
		}
		return y
	}
}

func (n Node) Import() []byte {
	var y []interface{}
	for _, v := range n.Nodes {
		y = append(y, v.Walk())
	}

	b, _ := json.Marshal(y)
	return b
}

// func (n Node) Pretty(indent int) {
// 	var tabs string
// 	for i := 0; i < indent; i++ {
// 		tabs += "\t"
// 	}

// 	for k, v := range n.Nodes {
// 		if len(v.Nodes) < 1 {
// 			fmt.Printf("%s(%d) %s == content node (%s)\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)), string(v.Content))
// 		} else {
// 			fmt.Printf("%s(%d) %s == parent node\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)))
// 			v.Pretty(indent + 1)
// 		}
// 	}
// }

func Pretty(name string) (string, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}

	var data interface{}

	if err := json.Unmarshal(b, &data); err != nil {
		return "", err
	}
	b, err = json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (n Node) Write(name string, perm os.FileMode) {
	j := n.Import()
	ioutil.WriteFile(name, j, perm)
}
