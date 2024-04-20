package ast

import (
	"bytes"

	"github.com/fengshux/monkey/token"
)

type ReturnStatment struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatment) statementNode() {}

func (r *ReturnStatment) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatment) String() string {
	var out bytes.Buffer
	out.WriteString(r.TokenLiteral() + " ")

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}
