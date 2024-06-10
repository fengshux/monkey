package evaluator

import (
	"fmt"

	"github.com/fengshux/monkey/ast"
	"github.com/fengshux/monkey/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, statement := range program.Statements {
		if isMacroDefinition(statement) {
			addMacro(statement, env)
			definitions = append(definitions, i)
		}
	}

	for i := len(definitions) - 1; i >= 0; i-- {
		index := definitions[i]
		program.Statements = append(program.Statements[:index],
			program.Statements[index+1:]...,
		)
	}
}

func isMacroDefinition(node ast.Statement) bool {
	letStatement, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}
	_, ok = letStatement.Value.(*ast.MacroLiteral)
	return ok
}

func addMacro(node ast.Statement, env *object.Environment) {
	letStatement, _ := node.(*ast.LetStatement)
	macroLiteral, _ := letStatement.Value.(*ast.MacroLiteral)

	macro := &object.Marco{
		Parameters: macroLiteral.Parameters,
		Body:       macroLiteral.Body,
		Env:        env,
	}
	env.Set(letStatement.Name.Value, macro)
}

func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpression, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		macro, ok := isMacroCall(callExpression, env)
		if !ok {
			return node
		}

		args := quoteArgs(callExpression)
		evalEnv := extendMacroEnv(macro, args)

		evaluated := Eval(macro.Body, evalEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST-node from macros")
		}

		return quote.Node
	})
}

func isMacroCall(exp *ast.CallExpression, env *object.Environment) (*object.Marco, bool) {
	identifer, ok := exp.Function.(*ast.Identifer)
	if !ok {
		return nil, false
	}

	obj, ok := env.Get(identifer.Value)
	if !ok {
		return nil, false
	}

	macro, ok := obj.(*object.Marco)
	if !ok {
		return nil, false
	}
	return macro, true
}

func quoteArgs(exp *ast.CallExpression) []*object.Quote {
	args := make([]*object.Quote, 0, len(exp.Arguments))

	for _, arg := range exp.Arguments {
		args = append(args, &object.Quote{Node: arg})
	}
	return args
}

func extendMacroEnv(macro *object.Marco, args []*object.Quote) *object.Environment {
	extend := object.NewEnclosedEnvironment(macro.Env)

	for idx, param := range macro.Parameters {
		extend.Set(param.Value, args[idx])
	}
	return extend
}
