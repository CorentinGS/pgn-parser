package main

import (
	"testing"
)

func TestLexer(t *testing.T) {
	input := `[Event "Example"]
[Site "Internet"]
[Date "2023.12.06"]
[Round "1"]
[White "Player1"]
[Black "Player2"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 {This is the Ruy Lopez.} 3... a6 1-0`

	lexer := NewLexer(input)
	expectedTokens := []struct {
		typ   TokenType
		value string
	}{
		{TAG_START, "["},
		{TAG_KEY, "Event"},
		{TAG_VALUE, "Example"},
		{TAG_END, "]"},
		{TAG_START, "["},
		{TAG_KEY, "Site"},
		{TAG_VALUE, "Internet"},
		{TAG_END, "]"},
		{TAG_START, "["},
		{TAG_KEY, "Date"},
		{TAG_VALUE, "2023.12.06"},
		{TAG_END, "]"},
		{TAG_START, "["},
		{TAG_KEY, "Round"},
		{TAG_VALUE, "1"},
		{TAG_END, "]"},
		{TAG_START, "["},
		{TAG_KEY, "White"},
		{TAG_VALUE, "Player1"},
		{TAG_END, "]"},
		{TAG_START, "["},
		{TAG_KEY, "Black"},
		{TAG_VALUE, "Player2"},
		{TAG_END, "]"},
		{TAG_START, "["},
		{TAG_KEY, "Result"},
		{TAG_VALUE, "1-0"},
		{TAG_END, "]"},
		{MOVE_NUMBER, "1"},
		{DOT, "."},
		{SQUARE, "e4"},
		{SQUARE, "e5"},
		{MOVE_NUMBER, "2"},
		{DOT, "."},
		{PIECE, "N"},
		{SQUARE, "f3"},
		{PIECE, "N"},
		{SQUARE, "c6"},
		{MOVE_NUMBER, "3"},
		{DOT, "."},
		{PIECE, "B"},
		{SQUARE, "b5"},
		{COMMENT_START, "{"},
		{COMMENT, "This is the Ruy Lopez."},
		{COMMENT_END, "}"},
		{MOVE_NUMBER, "3"},
		{ELLIPSIS, "..."},
		{SQUARE, "a6"},
		{RESULT, "1-0"},
	}

	for i, expected := range expectedTokens {
		token := lexer.NextToken()
		if token.Type != expected.typ || token.Value != expected.value {
			t.Errorf("Token %d - Expected {%v, %q}, got {%v, %q}",
				i, expected.typ, expected.value, token.Type, token.Value)
		}
	}

	// Test EOF
	token := lexer.NextToken()
	if token.Type != EOF {
		t.Errorf("Expected EOF token, got %v", token.Type)
	}
}

func TestCaptures(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "Pawn capture",
			input: "exf3",
			expected: []Token{
				{Type: FILE, Value: "e"},
				{Type: CAPTURE, Value: "x"},
				{Type: SQUARE, Value: "f3"},
			},
		},
		{
			name:  "Piece capture",
			input: "Nxc6",
			expected: []Token{
				{Type: PIECE, Value: "N"},
				{Type: CAPTURE, Value: "x"},
				{Type: SQUARE, Value: "c6"},
			},
		},
		{
			name:  "Piece capture with file disambiguation",
			input: "Nbxc6",
			expected: []Token{
				{Type: PIECE, Value: "N"},
				{Type: FILE, Value: "b"},
				{Type: CAPTURE, Value: "x"},
				{Type: SQUARE, Value: "c6"},
			},
		},
		{
			name:  "Piece capture with rank disambiguation",
			input: "N4xd5",
			expected: []Token{
				{Type: PIECE, Value: "N"},
				{Type: RANK, Value: "4"},
				{Type: CAPTURE, Value: "x"},
				{Type: SQUARE, Value: "d5"},
			},
		},
		{
			name:  "Complex position with captures",
			input: "1. e4 d5 2. Nf3 Nc6 3. Nbxd5",
			expected: []Token{
				{Type: MOVE_NUMBER, Value: "1"},
				{Type: DOT, Value: "."},
				{Type: SQUARE, Value: "e4"},
				{Type: SQUARE, Value: "d5"},
				{Type: MOVE_NUMBER, Value: "2"},
				{Type: DOT, Value: "."},
				{Type: PIECE, Value: "N"},
				{Type: SQUARE, Value: "f3"},
				{Type: PIECE, Value: "N"},
				{Type: SQUARE, Value: "c6"},
				{Type: MOVE_NUMBER, Value: "3"},
				{Type: DOT, Value: "."},
				{Type: PIECE, Value: "N"},
				{Type: FILE, Value: "b"},
				{Type: CAPTURE, Value: "x"},
				{Type: SQUARE, Value: "d5"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)

			for i, expected := range tt.expected {
				token := lexer.NextToken()
				if token.Type != expected.Type || token.Value != expected.Value {
					t.Errorf("Token %d - Expected {%v, %q}, got {%v, %q}",
						i, expected.Type, expected.Value, token.Type, token.Value)
				}
			}

			// Verify we get EOF after all tokens
			token := lexer.NextToken()
			if token.Type != EOF {
				t.Errorf("Expected EOF token after capture, got %v", token.Type)
			}
		})
	}
}

// Test captures in context
func TestCapturesInGame(t *testing.T) {
	input := "1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Bxc6"

	expectedTokens := []struct {
		typ   TokenType
		value string
	}{
		{MOVE_NUMBER, "1"},
		{DOT, "."},
		{SQUARE, "e4"},
		{SQUARE, "e5"},
		{MOVE_NUMBER, "2"},
		{DOT, "."},
		{PIECE, "N"},
		{SQUARE, "f3"},
		{PIECE, "N"},
		{SQUARE, "c6"},
		{MOVE_NUMBER, "3"},
		{DOT, "."},
		{PIECE, "B"},
		{SQUARE, "b5"},
		{SQUARE, "a6"},
		{MOVE_NUMBER, "4"},
		{DOT, "."},
		{PIECE, "B"},
		{CAPTURE, "x"},
		{SQUARE, "c6"},
	}

	lexer := NewLexer(input)

	for i, expected := range expectedTokens {
		token := lexer.NextToken()
		if token.Type != expected.typ || token.Value != expected.value {
			t.Errorf("Token %d - Expected {%v, %q}, got {%v, %q}",
				i, expected.typ, expected.value, token.Type, token.Value)
		}
	}
}

func TestCaslting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "Short castle",
			input: "O-O",
			expected: []Token{
				{Type: KINGSIDE_CASTLE, Value: "O-O"},
			},
		},
		{
			name:  "Long castle",
			input: "O-O-O",
			expected: []Token{
				{Type: QUEENSIDE_CASTLE, Value: "O-O-O"},
			},
		},
		{
			name:  "Short castle in game",
			input: "1. e4 e5 2. Nf3 Nc6 3. Bc4 Nf6 4. O-O",
			expected: []Token{
				{Type: MOVE_NUMBER, Value: "1"},
				{Type: DOT, Value: "."},
				{Type: SQUARE, Value: "e4"},
				{Type: SQUARE, Value: "e5"},
				{Type: MOVE_NUMBER, Value: "2"},
				{Type: DOT, Value: "."},
				{Type: PIECE, Value: "N"},
				{Type: SQUARE, Value: "f3"},
				{Type: PIECE, Value: "N"},
				{Type: SQUARE, Value: "c6"},
				{Type: MOVE_NUMBER, Value: "3"},
				{Type: DOT, Value: "."},
				{Type: PIECE, Value: "B"},
				{Type: SQUARE, Value: "c4"},
				{Type: PIECE, Value: "N"},
				{Type: SQUARE, Value: "f6"},
				{Type: MOVE_NUMBER, Value: "4"},
				{Type: DOT, Value: "."},
				{Type: KINGSIDE_CASTLE, Value: "O-O"},
			},
		},
		{
			name:  "Long castle in game",
			input: "1. e4 e5 2. Nf3 Nc6 3. Bc4 Nf6 4. O-O-O",
			expected: []Token{
				{Type: MOVE_NUMBER, Value: "1"},
				{Type: DOT, Value: "."},
				{Type: SQUARE, Value: "e4"},
				{Type: SQUARE, Value: "e5"},
				{Type: MOVE_NUMBER, Value: "2"},
				{Type: DOT, Value: "."},
				{Type: PIECE, Value: "N"},
				{Type: SQUARE, Value: "f3"},
				{Type: PIECE, Value: "N"},
				{Type: SQUARE, Value: "c6"},
				{Type: MOVE_NUMBER, Value: "3"},
				{Type: DOT, Value: "."},
				{Type: PIECE, Value: "B"},
				{Type: SQUARE, Value: "c4"},
				{Type: PIECE, Value: "N"},
				{Type: SQUARE, Value: "f6"},
				{Type: MOVE_NUMBER, Value: "4"},
				{Type: DOT, Value: "."},
				{Type: QUEENSIDE_CASTLE, Value: "O-O-O"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)

			for i, expected := range tt.expected {
				token := lexer.NextToken()
				if token.Type != expected.Type || token.Value != expected.Value {
					t.Errorf("Token %d - Expected {%v, %q}, got {%v, %q}",
						i, expected.Type, expected.Value, token.Type, token.Value)
				}
			}

			// Verify we get EOF after all tokens
			token := lexer.NextToken()
			if token.Type != EOF {
				t.Errorf("Expected EOF token after capture, got %v", token.Type)
			}
		})
	}
}
