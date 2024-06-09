package ast

import (
	"bytes"
	"strings"

	"github.com/fengshux/monkey/token"
)

type MacroLiteral struct {
	Token      token.Token
	Parameters []*Identifer
	Body       *BlockStatement
}

func (m *MacroLiteral) expressionNode() {}

func (m *MacroLiteral) TokenLiteral() string {
	return m.Token.Literal
}

func (m *MacroLiteral) String() string {
	var out bytes.Buffer

	params := make([]string, len(m.Parameters))
	for i, p := range m.Parameters {
		params[i] = p.String()
	}

	out.WriteString(m.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(m.Body.String())

	return out.String()
}
