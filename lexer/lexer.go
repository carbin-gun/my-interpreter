package lexer

import "my-interpreter/token"

type Lexer struct {
	input        string
	position     int  // 所输入字符串的当前位置（指向当前字符串）
	readPosition int  // 所输入字符串的当前读取位置 -- 指向当前字符的后一个字符
	ch           byte // 当前正在查看的字符
}

func New(input string) *Lexer {
	var l = &Lexer{input: input}
	l.readChar() // 初始化整个 lexer(包括 ch,position,readPosition变量)
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '=':
		if l.peekChar() == '=' { // 窥视下一个是 =， 就能组成一个特殊的 token 是 ==
			ch := l.ch
			l.readChar()                            // peekChar 是窥视，readChar 是真实移动指针
			var literal = string(ch) + string(l.ch) // 拿到实际的 == token，这时候指针都移动完毕
			tok = token.Token{                      // 因为 newToken 只接受第二个是 byte，所以直接用 token.Token的构造器
				Type:    token.EQ,
				Literal: literal,
			}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' { // 窥视下一个是 =， 就能组成一个特殊的 token 是 !=
			ch := l.ch
			l.readChar()                            // peekChar 是窥视，readChar 是真实移动指针
			var literal = string(ch) + string(l.ch) // 拿到实际的 != token，这时候指针都移动完毕
			tok = token.Token{                      // 因为 newToken 只接受第二个是 byte，所以直接用 token.Token的构造器
				Type:    token.NOT_EQ,
				Literal: literal,
			}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case 0:
		//tok = newToken(token.EOF, l.ch) // 不能用l.ch来赋值，因为实际这个 literal 不是空字符串的 byte 形式
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal) // 移动了 position 和 readPosition，所以需要提前 return，不能再用外部的 readChar 去移动两个标识位
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber() // 移动了 position 和 readPosition，所以需要提前 return，不能再用外部的 readChar 去移动两个标识位
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}
func (l *Lexer) readIdentifier() string {
	var position = l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || b == '_'
}
func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func newToken(tokenType token.TokenType, literal byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(literal),
	}
}
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 0 是 ASCII 码的 NULL --> 这也表明这个lexer 只支持 ASCII 字符，不能支持所有的 unicode 字符
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// peekChar 「窥视」函数，不会前移 position 和 readPosition，仅仅是「窥视」一下下一个字符，不移动指针。 窥视下一个字符，所以用 readPosition
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	var position = l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}
