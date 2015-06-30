package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"strconv"
)

type Node struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
	Nodes   []Node `xml:",any"`
}

func main() {
	// filename := "simple.xml"
	filename := "xml/simple.xml"

	b, _ := ioutil.ReadFile(filename)
	j := ImportXML(b)
	ioutil.WriteFile("json/simple.json", j, 0700)
}

func Walk(n Node, indent int) interface{} {

	var tabs string
	for i := 0; i < indent; i++ {
		tabs += "\t"
	}

	if len(n.Nodes) < 1 {
		// m[n.XMLName.Local] = string(n.Content)
		return string(n.Content)
	} else {
		// m[n.XMLName.Local] = make([]interface{}, len(n.Nodes))
		// m[n.XMLName.Local] = make(map[int]interface{})

		// var x []interface{}
		x := make(map[string][]interface{})
		y := make(map[string]interface{})

		for _, v := range n.Nodes {
			// x = append(x, Walk(v, indent+1))

			x[v.XMLName.Local] = append(x[v.XMLName.Local], Walk(v, indent+1))

			// if len(v.Nodes) < 1 {
			// 	fmt.Printf("%s(%d) %s == content node (%s)\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)), string(v.Content))
			// 	// m[n.XMLName.Local] = append(m[n.XMLName.Local].([]interface{}), string(v.Content))

			// } else {
			// 	fmt.Printf("%s(%d) %s == parent node\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)))
			// 	// m[n.XMLName.Local] = append(m[n.XMLName.Local].([]interface{}), Walk(v, indent+1))
			// 	m[n.XMLName.Local] = Walk(v, indent+1)

			// }
		}

		for k, v := range x {
			if len(x[k]) == 1 {
				switch v[0].(type) {
				case int, int32, int64:
					i, _ := strconv.Atoi(v[0].(string))
					y[k] = i //v[0].(int)
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

	// for k, v := range n.Nodes {
	// 	if len(v.Nodes) < 1 {
	// 		fmt.Printf("%s(%d) %s == content node (%s)\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)), string(v.Content))
	// 		m[v.XMLName.Local] = string(v.Content)
	// 	} else {
	// 		fmt.Printf("%s(%d) %s == parent node\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)))
	// 		m[v.XMLName.Local] = Walk(v, indent+1)
	// 	}
	// }

	// return m
}

func ImportXML(data []byte) []byte {
	buf := bytes.NewBuffer(data)
	dec := xml.NewDecoder(buf)

	var n Node
	dec.Decode(&n)

	var y []interface{}
	for _, v := range n.Nodes {
		y = append(y, Walk(v, 0))

	}

	// m := Walk(n, 0)
	b, _ := json.Marshal(y)
	// fmt.Println(string(b))
	return b

}
