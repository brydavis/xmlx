package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
	Nodes   []Node `xml:",any"`
}

func main() {
	n := Nodify("xml/complex.xml")
	n.Write("json/complex.json", 0700)
	n.Pretty(0)
}

func Walk(n Node) interface{} {
	if len(n.Nodes) < 1 {
		return string(n.Content)
	} else {
		x := make(map[string][]interface{})
		y := make(map[string]interface{})

		for _, v := range n.Nodes {
			x[v.XMLName.Local] = append(x[v.XMLName.Local], Walk(v))
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

func Nodify(filename string) Node {
	data, _ := ioutil.ReadFile(filename)
	buf := bytes.NewBuffer(data)
	dec := xml.NewDecoder(buf)

	var n Node
	dec.Decode(&n)

	return n
}

func (n Node) Import() []byte {
	var y []interface{}
	for _, v := range n.Nodes {
		y = append(y, Walk(v))
	}

	b, _ := json.Marshal(y)
	return b
}

func (n Node) Pretty(indent int) {
	var tabs string
	for i := 0; i < indent; i++ {
		tabs += "\t"
	}

	for k, v := range n.Nodes {
		if len(v.Nodes) < 1 {
			fmt.Printf("%s(%d) %s == content node (%s)\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)), string(v.Content))
		} else {
			fmt.Printf("%s(%d) %s == parent node\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)))
			v.Pretty(indent + 1)
		}
	}
}

func (n Node) Write(name string, perm os.FileMode) {
	j := n.Import()
	ioutil.WriteFile(name, j, perm)
}
