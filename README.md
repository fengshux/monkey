# Monkey

*本工程是根据《用Go语言资质解释器》（Thorsten Ball） 里构建Monkey语言进行构建的。本工程里的代码和技术均引用自该书* [Monkey源代码](https://github.com/fengshux/monkey)

Monkey 语言是为了学习编译原理，自制的一门语言。Monkey是一门解释性语言，本项目实现了Monkey语言的解释器。
Monkey语言具有以下特性：
* 类C语语法
* 变量绑定
* 整形和布尔类型
* 算数表达式
* 内置函数
* 头等函数和高阶函数
* 闭包
* 字符串数据结构
* 数组数据结构
* 哈希数据结构
* 宏系统

## 如何运行
  
```shell
go run main.go


> let a = 1;
> let b = 2;
> a + b;
// 3
> let m = {"a":1, "b":2};
> m["a"]
// 1
> fn add(a, b) {return a+b;}
> add(1, 2)
// 3

```

## 基础数据类型
Monkey语言支持整形和布尔类型字符串和空类型。
```shell
> 1
// 1
> true
// true
> "abcd"
// abcd

```
## 数组和字典
Monkey语言支持数组和字典类型,并支持对数组和字典的索引运算
```shell
> let a = ["a", 1, true]
> a[2]
// true
> let h = {"a":1, "b":false, c:"c"}
> h["c"]
// c                                                                            
```

## 函数和闭包
Monkey语言支持函数和闭包
```shell
> let a  = fn(n){fn(b) {return b+n;};}
> a(1)(2);
// 3
```

## 内置函数
Monkey语言支持内置函数`len`、`first`、 `rest`、`puts`等内置函数, 其中`puts`是输出函数，会向std打印传入`puts`的参数。
```shell
> len([1, 2, 3])
// 3
> first([1, 2, 3])
// 1
> rest([1, 2, 3])
// [2, 3]
> puts("hello world")
// hello world
```  

## 解释器的实现

在解释性语言中，解释器从源代码到得到运行结果经过了，词法分析、语法分析、宏展开、和求值的过程。
1. 词法分析， 将源代码解析成有意义的符号（token）
2. 语法分析， 将词法分析得出的符号（token）解析成抽象语法树（AST）
3. 宏展开，将抽象语法树（AST）中的宏定义删除，并且将抽象语法书中宏调用，进行宏展开，用展开的得到的AST节点，替换原来宏调用的节点。
4. 求值，将AST进行求值



## 词法分析器
解释器的第一个步骤是词法分析，词法分析将源代码转换为词法单元。词法单元是短小易于分类的语法结构，用于之后的语法分析。词法分析器有时候也叫词法单元生成器或扫描器。


```golang
// lexer/lexer.go
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Literal = string(ch) + string(l.ch)
			tok.Type = token.EQ
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
    .... 
    	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
```

词法分析器的主要逻辑在`NextToken`方法里，每次调用一次NextToken方法，就会返回下一段源代码对应的词法单元。

`NextToken`的逻辑即是简单的字符串扫描，首先跳过空白字符，之后扫描当前字符，根据对应的符号，以及当前符号可能的的情况判定具体是哪种词法单元。如，运算符、括号、逗号、变量、字符串、数字等。

## 语法分析器

语法分析器主要作用是将词法分析器获得的一连串的词法单元构建成抽象语法树。在此过程中顺便进行有语法检测。因为构建抽象语法树的过程中，错误的语法导致无法构建抽象语法树。

构建语法树是因为，需要将程序从语义化的字符串转换成结构化的数据，计算机处理程序才方便处理。而前期的词法分析其实上是为了语法分析做的准备工作。将字符串格式的代码转换为标准的词法单元，在语法分析时更具便捷和效率。

对编程语言语法分析时，主要由两种策略，自上而下的分析或自下而上的分析。
Monkey语言使用的是递归向下语法分析器，是基于自上而下的运算符优先级分析法的语法分析器， 是由沃恩·普拉特（Vaughan Pratt）发明的。因此又叫普拉特语法分析器。

在Monkey语言中分为语句和表达式，语句则不会返回值，表达式会返回值。


**ParseProgram**
```golang
// parser/parser.go
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}
```
在语法分析器的最顶层为`ParseProgram`,在此方法中，将循环对程序中的的语句进行逐一扫描，每一行或每一个代码块都视为语句。

**ParseStatement**
```golang
// parser/parser.go
func (p *Parser) ParseStatement() ast.Statement {

	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}
```
在`ParseStatement`中，分为三种类型语句：`letSatement`即赋值语句或者说变量绑定语句。`returnStatment`即return语句。其它的均为表达式。语法解析器的主要逻辑都是在解析表达式，将表达式构建成抽象语法树。当然赋值语句和return语句也很重要，由于篇幅限制这里没有列出，感兴趣可以去看源码。

**parseExpression**
```golang
// parser/parser.go
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParsefns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}
```
在`parseExpression`中分为前缀表达式即`prefixExpression`和中缀表达式即`infixExpression`。根据表达式实际的token或运算符，获取事先注册好的对应的表达式解析方法进行语法解析。

**ExpressionParserRegister**
```golang
// parser/parser.go
	p.registerPrefix(token.IDENT, p.parseIdentifer)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parserHashLiteral)
	p.registerPrefix(token.MACRO, p.parseMacroLiteral)

	// 以下这些符号都用parseInfixExpression
	for _, v := range []token.TokenType{token.PLUS, token.MINUS,
		token.SLASH, token.ASTERISK, token.EQ, token.NOT_EQ,
		token.LT, token.GT,
	} {
		p.registerInfix(v, p.parseInfixExpression)
	}
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
```    
前缀表达式如字面量、变量、函数调用、if表达式、负号、!(非)等都属于前缀表达式。中缀表达式包括算数运算、逻辑运算、索引运算等都属与中缀表达式。每种表达式都注册了自己的语法分析方法。在`parseExpression`根据对应的符号来获取调用。


**优先级**
```golang
// parser/parser.go
const (
	_ int = iota
	LOWEST
	EQUALS     // ==
	LESSGEAGER // < or >
	SUM        // +
	PRODUCT    // *
	PREFIX     // -x or !x
	CALL       // myFunc()
	INDEX
)
```
在解析表达式时，涉及到运算符的优先级问题，事先定义好了各种运算符的有限级，在解析表达式时，可根据优先级来决定运算符在语法树中的层级。优先级越低，在语法树中离根节点越近，即求值时越在后面计算。在求值的过程中，率先对子节点求值。

## 求值
求值是对已经建好的抽象语法树，不断的递归调用`Eval`函数的过程。从根节点开始，递归调用`Eval`函数找到叶子节点，在根据节点的运算符进行求值，然后返回，对上层节点求值，直到根节点。

****
```golang
// evaluator/evaluator.go
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// 语句
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
    ……
    }
}


func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpressoin(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
```
例如叶子节点为中缀运算时，调用`evalInfixExpression`函数进行求值，`evalInfixExpression`函数中根数据类型以及运算符调用对应的求值函数求值。

## 对象系统
语言的对象系统，是在求值过程中用于存储和表示求值结果的对象。例如Monkey语言中的数字字面量`1`,在求值过程中，内存中存储的数据对象是什么呢？

**Object接口**

在Monkey语言的数据类型，在go语言实现的对象系统中，所有对象实现了`Object`接口
```golang
// object/object.go
type Object interface {
	Type() ObjectType
	Inspect() string
}
```

**整形、bool类型、字符串**

整形、bool类型、字符串分别对go语言的int64、bool和字符串进行了封装

```golang
// object/object.go
type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

type String struct {
	Value string
}
```
 
**数组**

数组是对go语言的切片进行了封装，切片的元素为`Object`,由于Monkey语言的数据类型都实现了`Object`接口，因此数组的元素可以是任何类型

```golang
// object/object.go
type Array struct {    
	Elements []Object
}
```

**字典**

字典的实现借助了go语言的`map`,key的类型为`HashKey`, `HashKey`里包含了key的数据类型和key计算所得的hash值两个属性。value的类型为`HashPair`，`HashPair`里包含了key的原始值和实际要存入字典里的原始value。

Monkey语言中，可以作为键的数据类型都实现了`Hashable`接口。通过调用key的`HashKey`方法，得到key对应的`Hashkey`

例子中展示了字符串类型`HashKey`方法

```golang
// object/object.go
type Hash struct {
	Pairs map[HashKey]HashPair
}

type HashPair struct {
	Key   Object
	Value Object
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

// 字符串类型的HashKey方法
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: uint64(h.Sum64())}
}

```

## 环境
程序求值的过程中离不开运行环境，所谓的运行环境，即是代码运行的上下文。自定义的函数，可以在代码中直接调用，全局变量也可以在代码块中直接使用，子代码块中定义了和全局变量重名的局部变量则会遮蔽变量，以上这些都是由于自定义的函数、全局变量、局部变量都存储在“环境”这个上下文中。

`Environment`里面封装了一个`map`用作数据存储。在求值过程中会随着`Eval`方法层层传递下去，在遇到方法定义和变量定义时，通过调用`Environment`的`Set`方法将函数和变量存储到`Environment`中。在调用函数、或使用变量时，则通过`Environment`的`Get`方法将对应的值取出来。

`Environment`同时也实现了对自身类型的嵌套引用，当`Eval`函数对子代码块(if或函数)进行求值时，就会调用`NewEnclosedEnvironment`方法生成新的`Environment`, 外面传入的旧`Environment`存入新`Environment`的outer变量。因此新的`Environment`继承了旧`Environment`的所有数据。在求值过程中，如果遇到变量定义则会调用新的`Environment`的`Set`方法存入新的`Environment`。使用时调用`Environment`的`Get`方法，先从当前的`Environment`中寻找，如果找不到再向`outer`的`Environment`中查找。因此实现了子作用域和变量遮蔽。



```golang
// object/environment.go
type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}


```


## 宏
宏实质上是使用宏代码，在编译时生成monkey代码。这样可以使monkey语言由更强的表现力.
monkey 语言的宏使用 macro关键字实现，在宏里面使用quote关键字和unquote关键字来控制是否对代码进行求值。
```
> let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); }
> reverse(2 + 2, 10 - 5);
> 1
```
`reverse(2 + 2, 10 - 5)`宏展开之后的代码为`(2 + 2) - (10 - 5)`。然后进行求值，结果为`1`。