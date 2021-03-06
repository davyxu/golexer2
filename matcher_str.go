package ulexer

import (
	"strings"
	"unicode"
)

// 匹配标识符
func Identifier() Matcher {
	return (*identifierMatcher)(nil)
}

type identifierMatcher int

func (*identifierMatcher) TokenType() string {
	return "Identifier"
}

func (self *identifierMatcher) Read(lex *Lexer) (tk *Token) {

	var count int
	for {
		c := lex.Peek(count)

		isBasic := unicode.IsLetter(c) || c == '_'

		switch {
		case count == 0 && isBasic:
		case count > 0 && (isBasic || unicode.IsDigit(c)):
		default:
			goto ExitFor
		}

		count++
	}

ExitFor:

	if count == 0 {
		return EmptyToken
	}

	tk = lex.NewToken(count, self)

	lex.Consume(count)

	return
}

// 包含字面量
func Contain(literal interface{}) Matcher {

	self := &literalMatcher{}

	switch v := literal.(type) {
	case string:
		self.literal = []rune(v)
	case rune:
		self.literal = []rune{v}
	default:
		panic("invalid contain")
	}

	return self
}

type literalMatcher struct {
	literal []rune
}

func (*literalMatcher) TokenType() string {
	return "Literal"
}

func (self *literalMatcher) Read(lex *Lexer) (tk *Token) {

	var count int
	for {
		c := lex.Peek(count)

		if count >= len(self.literal) {
			break
		}

		if c != self.literal[count] {
			break
		}

		count++
	}

	if count == 0 {
		return EmptyToken
	}

	tk = lex.NewToken(count, self)

	lex.Consume(count)

	return
}

// 匹配字符串
func String() Matcher {
	return (*stringMatcher)(nil)
}

type stringMatcher int

func (*stringMatcher) TokenType() string {
	return "String"
}

func (self *stringMatcher) Read(lex *Lexer) (tk *Token) {

	beginChar := lex.Peek(0)
	if beginChar != '"' && beginChar != '\'' {
		return EmptyToken
	}

	lex.Consume(1)

	var escaping bool
	var sb strings.Builder

	var count int
	for {
		c := lex.Peek(count)

		if escaping {
			switch c {
			case 'n':
				sb.WriteRune('\n')
			case 'r':
				sb.WriteRune('\r')
			case '"', '\'':
				sb.WriteRune(c)
			default:
				sb.WriteRune('\\')
				sb.WriteRune(c)
			}

			escaping = false
		} else if c != beginChar {
			if c == '\\' {
				escaping = true
			} else {
				sb.WriteRune(c)
			}
		} else {
			break
		}

		if c == '\n' || c == 0 {
			break
		}

		count++
	}

	if count == 0 {
		return EmptyToken
	}

	end := count + 1

	tk = lex.NewTokenLiteral(end, self, sb.String())

	lex.Consume(end)

	return
}
