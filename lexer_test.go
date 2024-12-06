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