package main

import (
	"strings"
	"unicode"
)

type TokenType int

const (
	EOF           TokenType = iota
	TAG_START               // [
	TAG_END                 // ]
	TAG_KEY                 // The key part of a tag (e.g., "Site")
	TAG_VALUE               // The value part of a tag (e.g., "Internet")
	MOVE_NUMBER             // 1, 2, 3, etc.
	DOT                     // .
	ELLIPSIS                // ...
	PIECE                   // N, B, R, Q, K
	SQUARE                  // e4, e5, etc.
	COMMENT_START           // {
	COMMENT_END             // }
	COMMENT                 // The comment text
	RESULT                  // 1-0, 0-1, 1/2-1/2
)

type Token struct {
	Value string
	Type  TokenType
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	inTag        bool
	inComment    bool
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.inComment {
		if l.ch == '}' {
			l.inComment = false
			l.readChar()
			return Token{Type: COMMENT_END, Value: "}"}
		}
		return l.readComment()
	}

	switch l.ch {
	case '[':
		l.inTag = true
		l.readChar()
		return Token{Type: TAG_START, Value: "["}
	case ']':
		l.inTag = false
		l.readChar()
		return Token{Type: TAG_END, Value: "]"}
	case '"':
		return l.readTagValue()
	case '{':
		l.readChar()
		l.inComment = true
		return Token{Type: COMMENT_START, Value: "{"}
	case '}':
		l.readChar()
		return Token{Type: COMMENT_END, Value: "}"}
	case '.':
		if l.peekChar() == '.' && l.readPosition+1 < len(l.input) && l.input[l.readPosition+1] == '.' {
			l.readChar()
			l.readChar()
			l.readChar()
			return Token{Type: ELLIPSIS, Value: "..."}
		}
		l.readChar()
		return Token{Type: DOT, Value: "."}
	case '-':
		return l.readResult()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if l.inTag {
			return l.readTagValue()
		}
		// Look ahead to see if this number is followed by a dot or hyphen
		position := l.position
		for l.ch != 0 && isDigit(l.ch) {
			l.readChar()
		}
		if l.ch == '.' {
			return Token{Type: MOVE_NUMBER, Value: l.input[position:l.position]}
		} else if l.ch == '-' {
			l.position = position
			l.readPosition = position + 1
			l.ch = l.input[position]
			return l.readResult()
		} else {
			// Reset position and try again as a regular number
			l.position = position
			l.readPosition = position + 1
			l.ch = l.input[position]
			return l.readNumber()
		}
	case 0:
		return Token{Type: EOF, Value: ""}
	default:
		if l.inTag && isLetter(l.ch) {
			return l.readTagKey()
		} else if isLetter(l.ch) {
			if unicode.IsUpper(rune(l.ch)) {
				return l.readPieceMove()
			}
			return l.readMove()
		}
	}

	tok := Token{Type: EOF, Value: string(l.ch)}
	l.readChar()
	return tok
}

func (l *Lexer) readNumber() Token {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return Token{Type: MOVE_NUMBER, Value: l.input[position:l.position]}
}

func (l *Lexer) readResult() Token {
	position := l.position
	for !isWhitespace(l.ch) && l.ch != 0 {
		l.readChar()
	}
	result := l.input[position:l.position]
	if isResult(result) {
		return Token{Type: RESULT, Value: result}
	}
	return Token{Type: MOVE_NUMBER, Value: result}
}

func (l *Lexer) readComment() Token {
	position := l.position
	for l.ch != '}' && l.ch != 0 {
		l.readChar()
	}
	comment := strings.TrimSpace(l.input[position:l.position])
	return Token{Type: COMMENT, Value: comment}
}

func (l *Lexer) readPieceMove() Token {
	// First capture the piece
	piece := string(l.ch)
	l.readChar()

	// If there's a capture or additional piece qualifier, we need to handle that
	// (not implemented in this version but would go here)

	// Return just the piece - the square will be read in the next token
	return Token{Type: PIECE, Value: piece}
}

func (l *Lexer) readMove() Token {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '-' {
		l.readChar()
	}

	word := l.input[position:l.position]

	// Check if it's a result
	if isResult(word) {
		return Token{Type: RESULT, Value: word}
	}

	// Check if it's a square
	if isSquare(word) {
		return Token{Type: SQUARE, Value: word}
	}

	// Must be some other identifier
	return Token{Type: SQUARE, Value: word}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readTagValue() Token {
	l.readChar() // skip opening quote
	position := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	value := l.input[position:l.position]
	l.readChar() // skip closing quote
	return Token{Type: TAG_VALUE, Value: value}
}

func (l *Lexer) readTagKey() Token {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return Token{Type: TAG_KEY, Value: l.input[position:l.position]}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isSquare(s string) bool {
	if len(s) != 2 {
		return false
	}
	return 'a' <= s[0] && s[0] <= 'h' && '1' <= s[1] && s[1] <= '8'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isResult(s string) bool {
	return s == "1-0" || s == "0-1" || s == "1/2-1/2" || s == "*"
}
