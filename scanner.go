package main

import (
	"bufio"
	"io"
	"strings"
)

// Game represents a complete PGN game as raw text before tokenization
type Game struct {
	Raw string
}

// Scanner handles reading PGN games from a source
type Scanner struct {
	reader *bufio.Reader
	buffer string
}

// NewScanner creates a new PGN scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(r),
	}
}

// ScanGame reads the next complete game from the input
// Returns nil when there are no more games
func (s *Scanner) ScanGame() (*Game, error) {
	var gameBuilder strings.Builder
	inGame := false
	bracketCount := 0
	emptyLineCount := 0

	// Handle buffered content first
	if s.buffer != "" {
		gameBuilder.WriteString(s.buffer)
		s.buffer = ""
		inGame = true
	}

	for {
		line, err := s.reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}

		trimmedLine := strings.TrimSpace(line)

		// Handle empty lines
		if trimmedLine == "" {
			if inGame {
				emptyLineCount++
				if emptyLineCount >= 2 && isGameComplete(gameBuilder.String()) {
					break
				}
			}
			if err == io.EOF {
				break
			}
			continue
		}

		emptyLineCount = 0

		// Start of a new game
		if strings.HasPrefix(trimmedLine, "[") {
			if inGame && isGameComplete(gameBuilder.String()) {
				s.buffer = line
				break
			}
			inGame = true
		}

		if inGame {
			gameBuilder.WriteString(line)
			bracketCount = updateBracketCount(line, bracketCount)

			if bracketCount == 0 && isGameComplete(trimmedLine) {
				if err == io.EOF {
					break
				}
				continue
			}
		}

		if err == io.EOF {
			break
		}
	}

	if gameBuilder.Len() > 0 {
		return &Game{Raw: gameBuilder.String()}, nil
	}

	return nil, io.EOF
}

// HasNext checks if there are more games to read
func (s *Scanner) HasNext() bool {
	if s.buffer != "" {
		return true
	}

	b, err := s.reader.Peek(1)
	if err != nil {
		return false
	}

	// Skip leading whitespace
	for len(b) > 0 && isWhitespace(b[0]) {
		s.reader.ReadByte()
		b, err = s.reader.Peek(1)
		if err != nil {
			return false
		}
	}

	return len(b) > 0
}

// TokenizeGame converts a raw game into tokens using the lexer
func TokenizeGame(game *Game) ([]Token, error) {
	if game == nil {
		return nil, nil
	}

	lexer := NewLexer(game.Raw)
	var tokens []Token

	for {
		token := lexer.NextToken()
		if token.Type == EOF {
			break
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// Helper functions
func isGameComplete(text string) bool {
	return strings.Contains(text, "1-0") ||
		strings.Contains(text, "0-1") ||
		strings.Contains(text, "1/2-1/2") ||
		strings.Contains(text, "*")
}

func updateBracketCount(line string, count int) int {
	for _, ch := range line {
		if ch == '{' {
			count++
		} else if ch == '}' {
			count--
		}
	}
	return count
}
