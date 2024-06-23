package parser

import (
	"fmt"
	"mokey-type/ast"
	"mokey-type/lexer"
	"mokey-type/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type Parser struct {
	l              *lexer.Lexer
	currentToken   token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}
	//prefix
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.DECREMENT, p.parsePrefixExpression)
	p.registerPrefix(token.INCREMENT, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FOR, p.parseForLoop)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	//infix
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	p.NextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) NextToken() {
	if p.peekToken.Type == token.EOF {
		p.currentToken = p.peekToken
		p.peekToken.Type = token.EOF
	} else {
		p.currentToken = p.peekToken
		p.peekToken = p.l.NextToken()
	}
}

func (p *Parser) ParseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.ParseLetStatement()
	case token.RETURN:
		return p.ParseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) ParseLetStatement() *ast.LetStatement {
	ls := &ast.LetStatement{Token: p.currentToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	ls.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.NextToken()
	ls.Value = p.parseExpression(LOWEST)
	if fl, ok := ls.Value.(*ast.FunctionLiteral); ok {
		fl.Name = ls.Name.Value
	}
	for p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return ls
}

func (p *Parser) ParseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}
	p.NextToken()
	statement.ReturnValue = p.parseExpression(LOWEST)
	for p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return statement
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit

}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.currentTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.NextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	es := &ast.ExpressionStatement{Token: p.currentToken}
	es.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return es
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if p.currentToken.Type == "" {
		return nil
	}
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.NextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) noPrefixParseFnError(t token.Token) {
	msg := fmt.Sprintf("no prefix parse function for %s found %+v", t.Type, t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}
	p.NextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}
	precedence := p.currentPrecedence()
	p.NextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}
	p.NextToken()

	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {

		statement := p.ParseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.NextToken()
	}

	return block
}

func (p *Parser) parseIfExpression() ast.Expression {
	expresssion := &ast.IfExpression{Token: p.currentToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.NextToken()
	expresssion.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expresssion.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.NextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expresssion.Alternative = p.parseBlockStatement()
	}
	return expresssion
}

func (p *Parser) parseParameters() []*ast.Identifier {
	parameters := []*ast.Identifier{}

	if p.currentTokenIs(token.RPAREN) {
		return parameters
	}

	ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	parameters = append(parameters, ident)
	for p.peekTokenIs(token.COMMA) {
		p.NextToken()
		p.NextToken()
		ident = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		parameters = append(parameters, ident)
	}

	if !p.peekTokenIs(token.RPAREN) {
		return nil
	}

	p.NextToken()

	return parameters
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: p.currentToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.NextToken()
	literal.Parameters = p.parseParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	literal.Body = p.parseBlockStatement()
	return literal
}

func (p *Parser) parseForLoop() ast.Expression {
	literal := &ast.ForLoop{Token: p.currentToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	if !p.expectPeek(token.LET) {
		return nil
	}

	letStatement := p.ParseLetStatement()
	if letStatement == nil {
		return nil
	}
	literal.Declaration = *letStatement
	p.NextToken()

	condition := p.parseExpression(LOWEST)
	if condition == nil {
		return nil
	}
	literal.Condition = condition
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	p.NextToken()

	consequence := p.parseExpression(LOWEST)
	if consequence == nil {
		return nil
	}
	literal.Consequence = consequence

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	literal.Body = p.parseBlockStatement()

	return literal
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.NextToken()
		return list
	}
	p.NextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.NextToken()
		p.NextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currentToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: p.currentToken, Function: function}
	expression.Arguments = p.parseExpressionList(token.RPAREN)
	return expression
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{Token: p.currentToken, Left: left}
	p.NextToken()
	expression.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return expression
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.currentToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)
	for !p.peekTokenIs(token.RBRACE) {
		p.NextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.NextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.currentToken.Type != token.EOF {
		var statement ast.Statement
		if p.currentToken.Type != "" {
			statement = p.ParseStatement()
		}
		if statement != nil && statement.TokenLiteral() != "" {
			program.Statements = append(program.Statements, statement)
		}
		p.NextToken()
	}
	return program
}
