package main

import (
	"bytes"
	// "encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
)

type Node struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
	Nodes   []Node `xml:",any"`
}

type M map[string]interface{}

func main() {
	filename := "simple.xml"
	b, _ := ioutil.ReadFile(filename)
	ImportXML(b)
}

func Walk(n Node, indent int) bool {
	var tabs string
	for i := 0; i < indent; i++ {
		tabs += "\t"
	}

	for k, v := range n.Nodes {
		if len(v.Nodes) < 1 {
			fmt.Printf("%s(%d) %s == content node (%s)\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)), string(v.Content))

		} else {
			fmt.Printf("%s(%d) %s == parent node\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)))
			Walk(v, indent+1)
		}
	}

	return true
}

func ImportXML(data []byte) {
	buf := bytes.NewBuffer(data)
	dec := xml.NewDecoder(buf)

	var n Node
	dec.Decode(&n)

	Walk(n, 0)
}
