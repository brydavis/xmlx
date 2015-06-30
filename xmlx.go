package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Node struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
	Nodes   []Node `xml:",any"`
}

func main() {
	filename := "simple.xml"
	b, _ := ioutil.ReadFile(filename)
	ImportXML(b)
}

// func Walk(n Node, indent int) interface{} {
// 	var tabs string
// 	for i := 0; i < indent; i++ {
// 		tabs += "\t"
// 	}

// 	vex := make(map[string]interface{})

// 	for key, val := range n.Nodes {
// 		if len(val.Nodes) < 1 {
// 			fmt.Printf("%s(%d) %s == content node (%s)\n", tabs, key, strings.Title(strings.ToLower(val.XMLName.Local)), string(val.Content))
// 			vex[val.XMLName.Local] = string(val.Content)
// 		} else {
// 			fmt.Printf("%s(%d) %s == parent node\n", tabs, key, strings.Title(strings.ToLower(val.XMLName.Local)))
// 			// y := Walk(val, indent+1)
// 			// fmt.Println(val.XMLName.Local, y, "\n")

// 			// var x []interface{}
// 			// for k, _ := range val.Nodes {
// 			// 	vex[val.XMLName.Local] = append(x, Walk(vval, indent+1))
// 			// y := Walk(val.Nodes[k], indent+1)
// 			// fmt.Println(val.XMLName.Local, y, "\n")
// 			// }
// 			// vex[val.XMLName.Local] = x

// 		}
// 	}

// 	return vex

func Walk(n Node, indent int) interface{} {
	m := make(map[string]interface{})

	var tabs string
	for i := 0; i < indent; i++ {
		tabs += "\t"
	}

	if len(n.Nodes) < 1 {
		m[n.XMLName.Local] = string(n.Content)
	} else {
		// m[n.XMLName.Local] = make([]interface{}, len(n.Nodes))
		// m[n.XMLName.Local] = make(map[int]interface{})

		var x []interface{}

		for _, v := range n.Nodes {
			x = append(x, Walk(v, indent+1))

			// if len(v.Nodes) < 1 {
			// 	fmt.Printf("%s(%d) %s == content node (%s)\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)), string(v.Content))
			// 	// m[n.XMLName.Local] = append(m[n.XMLName.Local].([]interface{}), string(v.Content))

			// } else {
			// 	fmt.Printf("%s(%d) %s == parent node\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)))
			// 	// m[n.XMLName.Local] = append(m[n.XMLName.Local].([]interface{}), Walk(v, indent+1))
			// 	m[n.XMLName.Local] = Walk(v, indent+1)

			// }
		}

		m[n.XMLName.Local] = x
	}

	// for k, v := range n.Nodes {
	// 	if len(v.Nodes) < 1 {
	// 		fmt.Printf("%s(%d) %s == content node (%s)\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)), string(v.Content))
	// 		m[v.XMLName.Local] = string(v.Content)
	// 	} else {
	// 		fmt.Printf("%s(%d) %s == parent node\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)))
	// 		m[v.XMLName.Local] = Walk(v, indent+1)
	// 	}
	// }

	return m
}

func ImportXML(data []byte) {
	buf := bytes.NewBuffer(data)
	dec := xml.NewDecoder(buf)

	var n Node
	dec.Decode(&n)

	// fmt.Println(Walk(n, 0))

	m := Walk(n, 0)
	b, _ := json.Marshal(m)
	fmt.Println(string(b))
}
