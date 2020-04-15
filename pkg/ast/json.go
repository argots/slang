package ast

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// JSON provides support for marshaling and unmarshaling a Node.
type JSON struct {
	LocMap
	Node
}

type jsonNode struct {
	Type   string `json:"type"`
	Val    string `json:"val,omitempty"`
	Op     string `json:"op,omitempty"`
	EndOp  string `json:"endop,omitempty"`
	Loc    string `json:"loc,omitempty"`
	EndLoc string `json:"endloc,omitempty"`
	Nodes  []JSON `json:"nodes,omitempty"`
}

// MarshalJSON marshals a node into JSON
func (j *JSON) MarshalJSON() ([]byte, error) {
	var jn jsonNode

	if j.Node == nil {
		return json.Marshal(nil)
	}

	literal := func(t, val, loc string) {
		jn.Type, jn.Val, jn.Loc = t, val, loc
	}
	expr := func(t, op, loc, endOp, endLoc string, x, y Node) {
		jn.Type, jn.Op, jn.Loc = t, op, loc
		jn.EndOp, jn.EndLoc = endOp, endLoc
		jn.Nodes = []JSON{{j.LocMap, x}, {j.LocMap, y}}
	}

	switch n := j.Node.(type) {
	case Number:
		literal("Number", n.Val, j.formatLoc(n.Loc))
	case Quote:
		literal("Quote", n.Val, j.formatLoc(n.Loc))
	case Ident:
		literal("Ident", n.Val, j.formatLoc(n.Loc))
	case *Expr:
		expr("Expr", n.Op, j.formatLoc(n.Loc), "", "", n.X, n.Y)
	case *Paren:
		startLoc, endLoc := j.formatLoc(n.StartLoc), j.formatLoc(n.EndLoc)
		expr("Paren", n.StartOp, startLoc, n.EndOp, endLoc, n.X, n.Y)
	case *Set:
		startLoc, endLoc := j.formatLoc(n.StartLoc), j.formatLoc(n.EndLoc)
		expr("Set", n.StartOp, startLoc, n.EndOp, endLoc, n.X, n.Y)
	case *Seq:
		startLoc, endLoc := j.formatLoc(n.StartLoc), j.formatLoc(n.EndLoc)
		expr("Seq", n.StartOp, startLoc, n.EndOp, endLoc, n.X, n.Y)
	}

	return json.Marshal(jn)
}

// UnmarshalJSON unmarshals a set of bytes into a node
func (j *JSON) UnmarshalJSON(data []byte) error {
	jn := jsonNode{Nodes: []JSON{{}, {}}}
	jn.Nodes[0].LocMap = j.LocMap
	jn.Nodes[1].LocMap = j.LocMap

	if err := json.Unmarshal(data, &jn); err != nil {
		return err
	}
	loc, endLoc := j.parseLoc(jn.Loc), j.parseLoc(jn.EndLoc)
	switch jn.Type {
	case "Expr":
		j.Node = &Expr{jn.Op, loc, jn.Nodes[0].Node, jn.Nodes[1].Node}
	case "Paren":
		j.Node = &Paren{jn.Op, jn.EndOp, loc, endLoc, jn.Nodes[0].Node, jn.Nodes[1].Node}
	case "Set":
		j.Node = &Set{jn.Op, jn.EndOp, loc, endLoc, jn.Nodes[0].Node, jn.Nodes[1].Node}
	case "Seq":
		j.Node = &Set{jn.Op, jn.EndOp, loc, endLoc, jn.Nodes[0].Node, jn.Nodes[1].Node}
	case "Number":
		j.Node = Number{jn.Val, loc}
	case "Quote":
		j.Node = Quote{jn.Val, loc}
	case "Ident":
		j.Node = Ident{jn.Val, loc}
	}
	return nil
}

func (j *JSON) parseLoc(s string) Loc {
	if j.LocMap == nil || s == "" {
		return Loc(0)
	}
	endIndex := strings.LastIndex(s, ":")
	end, err1 := strconv.Atoi(s[endIndex+1:])
	startIndex := strings.LastIndex(s[:endIndex], ":")
	start, err2 := strconv.Atoi(s[startIndex+1 : endIndex])

	if err1 != nil || err2 != nil {
		return Loc(0)
	}
	return j.LocMap.Add(s[:startIndex], uint32(start), uint32(end))
}

func (j *JSON) formatLoc(loc Loc) string {
	if j.LocMap == nil {
		return ""
	}
	source, start, end := j.LocMap.Get(loc)
	return fmt.Sprintf("%s:%d:%d", source, start, end)
}
