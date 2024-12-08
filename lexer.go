package main

import (
	"strings"
	"unicode"
)

type TokenType int

const (
	EOF              TokenType = iota
	TAG_START                  // [
	TAG_END                    // ]
	TAG_KEY                    // The key part of a tag (e.g., "Site")
	TAG_VALUE                  // The value part of a tag (e.g., "Internet")
	MOVE_NUMBER                // 1, 2, 3, etc.
	DOT                        // .
	ELLIPSIS                   // ...
	PIECE                      // N, B, R, Q, K
	SQUARE                     // e4, e5, etc.
	COMMENT_START              // {
	COMMENT_END                // }
	COMMENT                    // The comment text
	RESULT                     // 1-0, 0-1, 1/2-1/2
	CAPTURE                    // 'x' in moves
	FILE                       // a-h in moves when used as disambiguation
	RANK                       // 1-8 in moves when used as disambiguation
	KINGSIDE_CASTLE            // 0-0
	QUEENSIDE_CASTLE           // 0-0-0
	PROMOTION                  // = in moves
	PROMOTION_PIECE            // The piece being promoted to (Q, R, B, N)
	CHECK                      // + in moves
	CHECKMATE                  // # in moves
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

func (l *Lexer) readRank() Token {
	rank := string(l.ch)
	l.readChar()
	return Token{Type: RANK, Value: rank}
}

func (l *Lexer) readComment() Token {
	position := l.position
	for l.ch != '}' && l.ch != 0 {
		l.readChar()
	}
	comment := strings.TrimSpace(l.input[position:l.position])
	return Token{Type: COMMENT, Value: comment}
}

// Update readPieceMove to handle piece moves
func (l *Lexer) readPieceMove() Token {
	// Capture just the piece
	piece := string(l.ch)
	l.readChar()

	// Return just the piece - the square or capture will be read in subsequent tokens
	return Token{Type: PIECE, Value: piece}
}

func (l *Lexer) readMove() Token {
	position := l.position

	// For pawn captures
	if isFile(l.ch) {
		file := string(l.ch)
		l.readChar()

		// Check for capture
		if l.ch == 'x' {
			return Token{Type: FILE, Value: file}
		}
	}

	for isFile(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	// Get the total length of what we read
	length := l.position - position

	// If we read 3 characters, first one is disambiguation
	if length == 3 {
		l.readPosition = position + 1
		l.readChar()
		// Return just the first character as disambiguation
		return Token{Type: FILE, Value: string(l.input[position])}
	}

	return Token{Type: SQUARE, Value: l.input[position:l.position]}
}

func (l *Lexer) readPromotionPiece() Token {
	piece := string(l.ch)
	l.readChar()
	return Token{Type: PROMOTION_PIECE, Value: piece}
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

func (l *Lexer) readCastling() (Token, bool) {
	position := l.position

	// First character should be uppercase 'O'
	if l.ch != 'O' {
		return Token{}, false
	}

	// Check if we have enough characters for at least kingside castling (O-O)
	if l.position+2 >= len(l.input) {
		return Token{}, false
	}

	// Check for "O-O" pattern
	if l.peekChar() != '-' {
		return Token{}, false
	}
	l.readChar() // skip O
	l.readChar() // skip -

	if l.ch != 'O' {
		// Reset if pattern doesn't match
		l.position = position
		l.readPosition = position + 1
		l.ch = l.input[position]
		return Token{}, false
	}
	l.readChar() // skip O

	// Look ahead to see if this is queenside castling (O-O-O)
	if l.ch == '-' && l.peekChar() == 'O' {
		l.readChar() // skip -
		l.readChar() // skip O
		return Token{Type: QUEENSIDE_CASTLE, Value: "O-O-O"}, true
	}

	return Token{Type: KINGSIDE_CASTLE, Value: "O-O"}, true
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
	case 'x':
		l.readChar()
		return Token{Type: CAPTURE, Value: "x"}
	case '-':
		return l.readResult()
	case 'O':
		// Check for castling
		if token, isCastling := l.readCastling(); isCastling {
			return token
		}
		// If not castling, treat as a regular piece move
		return l.readPieceMove()
	case '=':
		l.readChar()
		return Token{Type: PROMOTION, Value: "="}
	case '+':
		l.readChar()
		return Token{Type: CHECK, Value: "+"}
	case '#':
		l.readChar()
		return Token{Type: CHECKMATE, Value: "#"}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if l.inTag {
			return l.readTagValue()
		}

		// Look at previous characters to determine context
		if l.position > 0 && unicode.IsUpper(rune(l.input[l.position-1])) {
			// If preceded by a piece, it's a rank disambiguation
			return l.readRank()
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
				// If it follows a promotion token, it's a promotion piece
				if l.position > 0 && l.input[l.position-1] == '=' {
					return l.readPromotionPiece()
				}
				return l.readPieceMove()
			}
			return l.readMove()
		}
	}

	tok := Token{Type: EOF, Value: string(l.ch)}
	l.readChar()
	return tok
}
