package documents

import (
	"bufio"
	"bytes"
	"html"
	"net/url"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	htmlrenderer "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// RichNode contains data only. Model output never becomes executable Typst.
type RichNode struct {
	Kind     string     `json:"kind"`
	Text     string     `json:"text"`
	URL      string     `json:"url"`
	Children []RichNode `json:"children"`
}

func markdownBlocks(value string) []RichNode {
	source := []byte(value)
	root := goldmark.New(goldmark.WithExtensions(extension.Table, extension.Strikethrough)).Parser().Parse(text.NewReader(source))
	var convert func(ast.Node, int) RichNode
	convert = func(n ast.Node, depth int) RichNode {
		r := RichNode{Kind: "group", Children: []RichNode{}}
		if depth > 48 {
			r.Kind = "text"
			r.Text = string(n.Text(source))
			return r
		}
		switch v := n.(type) {
		case *ast.Text:
			r.Kind = "text"
			r.Text = string(v.Segment.Value(source))
			if !v.IsRaw() {
				r.Text = markdownText(v.Segment.Value(source))
			}
			if v.SoftLineBreak() || v.HardLineBreak() {
				r.Text += "\n"
			}
		case *ast.String:
			r.Kind = "text"
			r.Text = string(v.Value)
		case *ast.Paragraph, *ast.TextBlock:
			r.Kind = "paragraph"
		case *ast.Heading:
			r.Kind = "heading"
		case *ast.Emphasis:
			if v.Level == 2 {
				r.Kind = "strong"
			} else {
				r.Kind = "emphasis"
			}
		case *ast.CodeSpan:
			r.Kind = "code"
		case *ast.FencedCodeBlock, *ast.CodeBlock:
			r.Kind = "codeblock"
			for i := 0; i < n.Lines().Len(); i++ {
				segment := n.Lines().At(i)
				r.Text += string(segment.Value(source))
			}
		case *ast.List:
			r.Kind = "list"
			if v.IsOrdered() {
				r.Kind = "ordered"
			}
		case *ast.ListItem:
			r.Kind = "item"
		case *ast.Blockquote:
			r.Kind = "quote"
		case *ast.Link:
			r.Kind = "link"
			u, err := url.Parse(string(v.Destination))
			if err == nil && (u.Scheme == "https" || u.Scheme == "http") && u.Host != "" {
				r.URL = u.String()
			}
		case *ast.AutoLink:
			r.Kind = "text"
			r.Text = string(v.Label(source))
		case *ast.Image:
			r.Kind = "group"
		case *ast.RawHTML:
			r.Kind = "text"
			for i := 0; i < v.Segments.Len(); i++ {
				segment := v.Segments.At(i)
				r.Text += string(segment.Value(source))
			}
		case *ast.HTMLBlock:
			r.Kind = "text"
			r.Text = string(v.Text(source))
		case *ast.ThematicBreak:
			r.Kind = "paragraph"
		case *extast.Table:
			r.Kind = "table"
		case *extast.TableHeader, *extast.TableRow:
			r.Kind = "row"
		case *extast.TableCell:
			r.Kind = "cell"
		case *extast.Strikethrough:
			r.Kind = "strike"
		}
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			r.Children = append(r.Children, convert(child, depth+1))
		}
		return r
	}
	result := []RichNode{}
	for n := root.FirstChild(); n != nil; n = n.NextSibling() {
		result = append(result, convert(n, 0))
	}
	return result
}

// Use CommonMark's text decoder, then discard HTML escaping before JSON encoding.
func markdownText(value []byte) string {
	var out bytes.Buffer
	writer := bufio.NewWriter(&out)
	htmlrenderer.DefaultWriter.Write(writer, value)
	_ = writer.Flush()
	return html.UnescapeString(out.String())
}
