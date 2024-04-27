package ast

import (
	"bytes"
	"strings"

	"github.com/fengshux/monkey/token"
)

type ExpressionStatement struct {
	Token      token.Token // 该表达式中的第一个词法单元
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token // 前缀词法单元
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteByte('(')
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteByte(')')

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // 运算符词法单元，如+
	Left     Expression
	Operator string
	Right    Expression
}

func (i InfixExpression) expressionNode() {}

func (i InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteByte('(')
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteByte(')')

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) statementNode() {}

func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range b.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifer
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode() {}

func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	param := []string{}
	for _, p := range f.Parameters {
		param = append(param, p.String())
	}

	out.WriteString(f.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(param, ", "))
	out.WriteString(")")
	out.WriteString(f.Body.String())
	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
