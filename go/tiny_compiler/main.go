package main

import (
	"fmt"
	"reflect"
)

type token struct {
	typ   string
	value string
}

// Tokenizer 用作词法分析，将输入的字符拆分成词，得到一个标记数组
func Tokenizer(input string) []token {
	var current int = 0
	var tokens []token

	for current < len(input) {
		char := input[current]
		if char == '(' {
			tokens = append(tokens, token{typ: "paren", value: "("})
			current++
			continue
		}

		if char == ')' {
			tokens = append(tokens, token{typ: "paren", value: ")"})
			current++
			continue
		}

		if char == ' ' {
			current++
			continue
		}

		if char >= '0' && char <= '9' {
			var numChar string
			for char >= '0' && char <= '9' {
				numChar += string(char)
				current++
				char = input[current]
			}
			tokens = append(tokens, token{typ: "number", value: numChar})
			continue
		}

		if char == '"' {
			current++
			char = input[current]

			var strChar string
			for char != '"' {
				strChar += string(char)
				current++
				char = input[current]
			}
			tokens = append(tokens, token{typ: "string", value: strChar})
			continue
		}

		if char >= 'a' && char <= 'z' {
			var chars string
			for char >= 'a' && char <= 'z' {
				chars += string(char)
				current++
				char = input[current]
			}
			tokens = append(tokens, token{typ: "name", value: chars})
			continue
		}

		panic("unknown character this character :" + string(char))
	}
	return tokens
}

type Leaf struct {
	typ   string
	value string
}

type Node struct {
	typ    string
	value  string
	params []interface{}
}

type AST struct {
	typ  string
	body []interface{}
}

func walk(current *int, tokens []token) interface{} {
	tok := tokens[*current]
	if tok.typ == "number" {
		*current++
		return Leaf{
			typ:   "NumberLiteral",
			value: tok.value,
		}
	}
	if tok.typ == "string" {
		*current++
		return Leaf{
			typ:   "StringLiteral",
			value: tok.value,
		}
	}
	if tok.typ == "paren" && tok.value == "(" {
		*current++
		tok = tokens[*current]
		var node Node = Node{
			typ:   "CallExpression",
			value: tok.value,
		}
		*current++
		tok = tokens[*current]
		for tok.typ != "paren" || tok.typ == "paren" && tok.value != ")" {
			node.params = append(node.params, walk(current, tokens))
			tok = tokens[*current]
		}
		*current++
		return node
	}
	return nil
}

// Parser 将词法分析的结果转变成 AST Abstract Syntax Tree 抽象语法树
func Parser(tokens []token) AST {
	var current = int(0)
	var ast AST = AST{
		typ: "Program",
	}
	for current < len(tokens) {
		ast.body = append(ast.body, walk(&current, tokens))
	}
	return ast
}

func traverseArray(array []interface{}, parent interface{}) {
	for i := range array {
		traverseNode(array[i], parent)
	}
}

var visitor map[string]interface{}

func traverseNode(node interface{}, parent interface{}) {
	typ := reflect.ValueOf(node).FieldByName("typ").Interface().(string)
	methods, ok := visitor[typ]
	if ok && reflect.ValueOf(methods).MethodByName("enter").IsValid() {
		var args []reflect.Value
		args = append(args, reflect.ValueOf(node))
		args = append(args, reflect.ValueOf(parent))
		reflect.ValueOf(methods).MethodByName("enter").Call(args)
	}
	switch typ {
	case "Program":
		traverseArray(node.(AST).body, node)
	case "CallExpression":
		traverseArray(node.(Node).params, node)
	case "NumberLiteral":
	case "StringLiteral":
	default:
		panic("type error" + typ)
	}
	if ok && reflect.ValueOf(methods).MethodByName("exit").IsValid() {
		var args []reflect.Value
		args = append(args, reflect.ValueOf(node))
		args = append(args, reflect.ValueOf(parent))
		reflect.ValueOf(methods).MethodByName("exit").Call(args)
	}
}

func traverser(ast AST, visitor interface{}) {
	traverseNode(ast, nil)
}

type expression struct {
	typ    string
	callee struct {
		typ  string
		name string
	}
	arguments []interface{}
}

// 啊啊啊~ 搞不定，js都不用定义变量类型，也没有要求定义返回值，就随便定义变量和返回值，go是强类型的，无法实现

func Transformer(ast AST) {
	//var newAst AST = AST{
	//	typ: "Program",
	//}
	visitor = make(map[string]interface{})
	visitor["NumberLiteral"] = func(node Node, parent Node) {
		parent.params = append(parent.params, Leaf{
			typ:   "NumberLiteral",
			value: node.value,
		})
	}
	visitor["StringLiteral"] = func(node Node, parent Node) {
		parent.params = append(parent.params, Leaf{
			typ:   "StringLiteral",
			value: node.value,
		})
	}
	visitor["CallExpression"] = func(node Node, parent Node) {
		var express expression = expression{
			typ: "CallExpression",
			callee: struct {
				typ  string
				name string
			}{typ: "Identifier", name: node.typ},
		}
		node.params = express.arguments
		if parent.typ != "CallExpression" {

		}
	}
	//traverser(ast)
}

func main() {
	tokens := Tokenizer("(add 2 2)")
	fmt.Println(tokens)
	ast := Parser(tokens)
	fmt.Println(ast)
}
