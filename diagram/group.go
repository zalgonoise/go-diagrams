package diagram

import (
	"strconv"

	graphviz "github.com/awalterschulze/gographviz"
)

type Background int

const (
	BackgroundBlue Background = iota
	BackgroundGreen
	BackgroundPurple
	BackgroundYellow
)

func (bg Background) String() string {
	switch bg {
	case BackgroundBlue:
		return "#E5F5FD"
	case BackgroundGreen:
		return "#EBF3E7"
	case BackgroundPurple:
		return "#ECE8F6"
	case BackgroundYellow:
		return "#FDF7E3"
	default:
		return BackgroundBlue.String()
	}
}

type Group struct {
	bg       Background
	id       string
	options  GroupOptions
	parent   *Group
	children map[string]*Group

	nodes map[string]*Node
	edges map[string]*Edge
}

func NewGroup(name string, opts ...GroupOption) *Group {
	return newGroup("cluster_"+name, BackgroundBlue, nil, opts...)

}

func newGroup(name string, bg Background, parent *Group, opts ...GroupOption) *Group {
	options := defaultGroupOptions(bg, opts...)

	return &Group{
		id:       name,
		bg:       bg,
		options:  options,
		parent:   parent,
		children: make(map[string]*Group),
		nodes:    make(map[string]*Node),
		edges:    make(map[string]*Edge),
	}
}

func (g *Group) ID() string {
	return g.id
}

func (g *Group) Nodes() []*Node {
	nodes := make([]*Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		nodes = append(nodes, n)
	}

	return nodes
}

func (g *Group) Edges() []*Edge {
	edges := make([]*Edge, 0, len(g.edges))
	for _, e := range g.edges {
		edges = append(edges, e)
	}

	return edges
}

func (g *Group) Children() []*Group {
	gs := make([]*Group, 0, len(g.children))

	for _, c := range g.children {
		gs = append(gs, c)
	}

	return gs
}

func (g *Group) Add(nodes ...*Node) *Group {
	for _, n := range nodes {
		g.nodes[n.ID()] = n
	}

	return g
}

func (g *Group) Connect(start, end *Node, opts ...EdgeOption) *Group {
	g.Add(start, end)
	return g.ConnectByID(start.ID(), end.ID(), opts...)
}

func (g *Group) ConnectByID(start, end string, opts ...EdgeOption) *Group {
	e := NewEdge(start, end, opts...)
	g.edges[e.ID()] = e

	return g
}

func (g *Group) ConnectAllTo(end string, opts ...EdgeOption) *Group {
	for id := range g.nodes {
		g.ConnectByID(id, end, opts...)
	}

	return g
}

func (g *Group) ConnectAllFrom(start string, opts ...EdgeOption) *Group {
	for id := range g.nodes {
		g.ConnectByID(start, id, opts...)
	}

	return g
}

func (g *Group) attrs() map[string]string {
	attrs := map[string]string{
		"label":     g.options.Label,
		"labeljust": g.options.LabelJustify,
		"pencolor":  g.options.PenColor,
		"bgcolor":   g.options.BackgroundColor,
		"shape":     g.options.Shape,
		"style":     g.options.Style,
		"fontname":  g.options.Font.Name,
		"fontsize":  strconv.FormatInt(int64(g.options.Font.Size), 10),
		"fontcolor": g.options.Font.Color,
	}

	for k, v := range g.options.Attributes {
		attrs[k] = v
	}

	return trimAttrs(attrs)
}

func (g *Group) Group(ng *Group) *Group {
	g.children[ng.id] = ng
	ng.parent = g

	return ng
}

func (g *Group) NewGroup(name string, opts ...GroupOption) *Group {
	ng := newGroup("cluster"+name, g.bg+1, g, opts...)
	g.children[ng.id] = ng
	ng.parent = g

	return ng
}

func (g *Group) Label(l string) *Group {
	g.options.Label = l
	return g
}

func (g *Group) BackgroundColor(c string) *Group {
	g.options.BackgroundColor = c
	return g
}

func (g *Group) render(outdir string, graph *graphviz.Escape) error {
	if err := graph.AddSubGraph(g.parent.id, g.id, g.attrs()); err != nil {
		return err
	}

	for _, n := range g.nodes {
		if err := n.render(g.id, outdir, graph); err != nil {
			return err
		}
	}

	for _, e := range g.edges {
		if err := e.render(e.start, e.end, graph); err != nil {
			return err
		}
	}

	for _, child := range g.children {
		if err := child.render(outdir, graph); err != nil {
			return err
		}
	}

	return nil
}

type GroupOptions struct {
	Label           string
	LabelJustify    string
	Direction       string
	PenColor        string
	BackgroundColor string
	Shape           string
	Style           string
	Font            Font
	Attributes      map[string]string
}

func DefaultGroupOptions(opts ...GroupOption) GroupOptions {
	return defaultGroupOptions(0, opts...)
}

func defaultGroupOptions(bg Background, opts ...GroupOption) GroupOptions {
	options := GroupOptions{
		LabelJustify: "l",
		Direction:    string(LeftToRight),
		PenColor:     "#AEB6BE",
		Shape:        "box",
		Style:        "rounded",
		Font: Font{
			Name:  "Sans-Serif",
			Size:  12,
			Color: "#2D3436",
		},
		Attributes: make(map[string]string),
	}

	WithBackground(bg)(&options)

	for _, o := range opts {
		o(&options)
	}

	return options
}

type GroupOption func(*GroupOptions)

func BackgroundColor(c string) GroupOption {
	return func(o *GroupOptions) {
		o.BackgroundColor = c
	}
}

func WithBackground(bg Background) GroupOption {
	return BackgroundColor(bg.String())
}

func GroupLabel(l string) GroupOption {
	return func(o *GroupOptions) {
		o.Label = l
	}
}
