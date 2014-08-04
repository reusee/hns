package hns

import (
	"fmt"
	"regexp"
	"strings"
)

var elemPattern = regexp.MustCompile(strings.Join([]string{
	`([a-zA-Z][a-zA-Z0-9_-]*)?`,     // tag
	`(#[a-zA-Z][a-zA-Z0-9_.:-]*)?`,  // id
	`(\.[a-zA-Z][a-zA-Z0-9_.:-]*)?`, // class
}, ""))

func Css(desc string, cont WalkFunc) WalkFunc {
	parts := strings.Split(desc, " ")
	for i := len(parts) - 1; i >= 0; i-- {
		part := strings.TrimSpace(parts[i])
		if len(part) == 0 {
			continue
		}
		matches := elemPattern.FindStringSubmatch(part)
		if matches[0] == "" {
			panic(fmt.Sprintf("bad css selector %s", desc))
		}
		var predicts []func(*Node) bool
		if matches[1] != "" { // tag
			tag := matches[1]
			predicts = append(predicts, func(n *Node) bool {
				return n.Tag == tag
			})
		}
		if matches[2] != "" { // id
			id := matches[2][1:]
			predicts = append(predicts, func(n *Node) bool {
				return n.Id == id
			})
		}
		if matches[3] != "" { // class
			class := matches[3][1:]
			predicts = append(predicts, func(n *Node) bool {
				for _, c := range n.Class {
					if c == class {
						return true
					}
				}
				return false
			})
		}
		cont = AllDescendant(func(ctx *WalkCtx, node *Node) bool {
			for _, predict := range predicts {
				if !predict(node) {
					return false
				}
			}
			return true
		}, cont)
	}
	return cont
}
