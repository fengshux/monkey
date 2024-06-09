package evaluator

import (
	"testing"

	"github.com/fengshux/monkey/ast"
	"github.com/fengshux/monkey/lexer"
	"github.com/fengshux/monkey/object"
	"github.com/fengshux/monkey/parser"
)

func TestDefineMacros(t *testing.T) {
	input := `
		let number = 1;
		let function = fn(x, y) { x + y; };
		let mymacro = macro(x, y) { x + y; };
	`
	env := object.NewEnvironment()
	program := testParseProgram(input)
	DefineMacros(program, env)

	if len(program.Statements) != 2 {
		t.Fatalf("Worng number of statements, got=%d want=%d", len(program.Statements), 2)
	}

	_, ok := env.Get("number")
	if ok {
		t.Fatalf("number should not be defined")
	}

}

func testParseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func TestExpandMacros(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{
			`let infixExpression = macro() { quote(1 + 2); };
			infixExpression();
			`,
			"(1 + 2)",
		},
		{
			`let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
			reverse(2 + 2, 10 - 5);
			`,
			"(10 - 5) - (2 + 2)",
		},
		{
			`
			let unless = macro(condition, consequence, alternative) {
				quote(if(!(unquote(condition))) {
					unquote(consequence);
				} else {
					unquote(alternative);
				});
			};

			unless(10 > 5, puts("no greater"), puts("greater"));
			`,
			`if(!(10 > 5)) { puts("no greater")} else { puts("greater") })`,
		},
	}

	for _, tt := range test {
		expected := testParseProgram(tt.expected)
		program := testParseProgram(tt.input)
		env := object.NewEnvironment()
		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)

		if expanded.String() != expected.String() {
			t.Errorf("not equal, want=%q, got=%q", expected.String(), expanded.String())
		}

	}

}
