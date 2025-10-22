package tmpl2xml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"sort"
	"strconv"
	"text/template/parse"
)

type nodeDefines struct {
	XMLName xml.Name `xml:"defines"`
	Defines []nodeDefine
}

type nodeDefine struct {
	XMLName  xml.Name `xml:"define"`
	Name     string   `xml:"template,attr,omitempty"`
	Template []any
}

type nodeThen struct {
	XMLName xml.Name `xml:"then"`
	List    []any
}

type nodeElse struct {
	XMLName xml.Name `xml:"else"`
	List    []any
}

type nodeIf struct {
	XMLName xml.Name  `xml:"if"`
	Pipe    string    `xml:"cond,attr,omitempty"`
	Then    nodeThen  `xml:"then"`
	Else    *nodeElse `xml:"else,omitempty"`
}

type nodeAction struct {
	XMLName xml.Name `xml:"action"`
	Pipe    string   `xml:",innerxml"`
}

type nodeList struct {
	List []any `xml:",innerxml"`
}

type nodeRange struct {
	XMLName xml.Name  `xml:"range"`
	Pipe    string    `xml:"pipe,attr,omitempty"`
	List    []any     `xml:",innerxml"`
	Else    *nodeElse `xml:"else,omitempty"`
}

type nodeWith struct {
	XMLName xml.Name  `xml:"with"`
	Pipe    string    `xml:"pipe,attr,omitempty"`
	List    []any     `xml:",innerxml"`
	Else    *nodeElse `xml:"else,omitempty"`
}

type nodeTemplate struct {
	XMLName xml.Name `xml:"template"`
	Name    string   `xml:"name,attr,omitempty"`
}

type nodeText struct {
	XMLName xml.Name `xml:"text"`
	Text    string   `xml:",innerxml"`
}

func listToAny(x *parse.ListNode, mode EscapeMode) []any {
	if x == nil {
		return nil
	}
	kids := make([]any, len(x.Nodes))
	for i, k := range x.Nodes {
		kids[i] = convert(k, mode)
	}
	return kids
}

func convert(n parse.Node, mode EscapeMode) any {
	switch x := n.(type) {
	case *parse.ListNode:
		return listToAny(x, mode)
	case *parse.TextNode:
		switch mode {
		case EscapeModeQuote:
			return nodeText{Text: strconv.Quote(string(x.Text))}
		case EscapeModeXML:
			var buf bytes.Buffer
			xml.EscapeText(&buf, x.Text)
			return nodeText{Text: buf.String()}
		default:
			panic(fmt.Sprintf("unexpected EscapeMode %d", mode))
		}
	case *parse.ActionNode:
		return nodeAction{Pipe: x.Pipe.String()}
	case *parse.RangeNode:
		v := nodeRange{Pipe: x.Pipe.String(), List: listToAny(x.List, mode)}
		if x.ElseList != nil {
			v.Else = &nodeElse{List: listToAny(x.ElseList, mode)}
		}
		return v
	case *parse.IfNode:
		v := nodeIf{Pipe: x.Pipe.String(), Then: nodeThen{List: listToAny(x.List, mode)}}
		if x.List != nil {
			v.Then = nodeThen{List: listToAny(x.List, mode)}
		}
		if x.ElseList != nil {
			v.Else = &nodeElse{List: listToAny(x.ElseList, mode)}
		}
		return v
	case *parse.WithNode:
		v := nodeWith{Pipe: x.Pipe.String(), List: listToAny(x.List, mode)}
		if x.ElseList != nil {
			v.Else = &nodeElse{List: listToAny(x.ElseList, mode)}
		}
		return v
	case *parse.TemplateNode:
		return nodeTemplate{Name: x.Name}
	case *parse.PipeNode:
		return string(x.String())
	default:
		panic(fmt.Sprintf("unexpected %T", n))
	}
}

type EscapeMode uint8

const (
	// EscapeModeQuote uses strconv.Quote to escape text
	EscapeModeQuote EscapeMode = iota
	// EscapeModeXML uses xml.EscapeText to escape text
	EscapeModeXML
)

type Converter struct {
	Encoder *xml.Encoder
	io.Writer
	buf        *bytes.Buffer
	EscapeMode EscapeMode
}

func (c *Converter) ensureEncoder() {
	if c.Encoder != nil {
		return
	}
	if c.Writer == nil {
		c.buf = new(bytes.Buffer)
		c.Writer = c.buf
	}
	c.Encoder = xml.NewEncoder(c.Writer)
	c.Encoder.Indent("", "  ")
}

func (c *Converter) FromTemplate(t *template.Template, text string) error {
	return c.FromTree(t.Tree, text)
}

func (c *Converter) FromTree(tr *parse.Tree, text string) error {
	c.ensureEncoder()
	treeSet := map[string]*parse.Tree{}
	tr, err := tr.Parse(text, "", "", treeSet)
	if err != nil {
		return err
	}
	names := make([]string, 0, len(treeSet))
	for name := range treeSet {
		names = append(names, name)
	}
	sort.Strings(names)

	templates := make([]nodeDefine, len(treeSet))
	for i, name := range names {
		tr := treeSet[name]
		templates[i] = nodeDefine{Name: name, Template: []any{convert(tr.Root, c.EscapeMode)}}
	}
	root := nodeDefines{Defines: templates}
	return c.Encoder.Encode(root)
}

func (c *Converter) FromString(text string) error {
	tr := parse.New("")
	tr.Mode = parse.SkipFuncCheck
	if c.buf != nil {
		c.buf.WriteString(xml.Header)
	}
	return c.FromTree(tr, text)
}

func String(text string) (string, error) {
	c := new(Converter)
	if err := c.FromString(text); err != nil {
		return "", err
	}
	return c.buf.String(), nil
}
