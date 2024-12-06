package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to read and tokenize a game
func readAndTokenize(scanner *Scanner) ([]Token, error) {
	game, err := scanner.ScanGame()
	if err != nil {
		return nil, err
	}
	return TokenizeGame(game)
}

func TestScanner(t *testing.T) {
	// Open the fixture file
	file, err := os.Open(filepath.Join("fixtures", "multi_game.pgn"))
	if err != nil {
		t.Fatalf("Failed to open fixture file: %v", err)
	}
	defer file.Close()

	scanner := NewScanner(file)

	// Test first game
	tokens, err := readAndTokenize(scanner)
	if err != nil {
		t.Fatalf("Failed to read first game: %v", err)
	}

	expectedFirstGame := []struct {
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

	if len(tokens) != len(expectedFirstGame) {
		t.Errorf("Expected %d tokens, got %d", len(expectedFirstGame), len(tokens))
		return
	}

	for i, expected := range expectedFirstGame {
		if i >= len(tokens) {
			t.Errorf("Missing token at position %d, expected {%v, %q}", i, expected.typ, expected.value)
			continue
		}

		if tokens[i].Type != expected.typ || tokens[i].Value != expected.value {
			t.Errorf("Token %d - Expected {%v, %q}, got {%v, %q}",
				i, expected.typ, expected.value, tokens[i].Type, tokens[i].Value)
		}
	}

	// Test second game (if exists)
	tokens, err = readAndTokenize(scanner)
	if err != nil && err != io.EOF {
		t.Errorf("Unexpected error reading second game: %v", err)
	}

	// Test HasNext functionality
	file.Seek(0, 0) // Reset file to beginning
	scanner = NewScanner(file)

	gameCount := 0
	for scanner.HasNext() {
		game, err := scanner.ScanGame()
		if err != nil {
			t.Errorf("Error reading game %d: %v", gameCount+1, err)
			break
		}

		_, err = TokenizeGame(game)
		if err != nil {
			t.Errorf("Error tokenizing game %d: %v", gameCount+1, err)
			break
		}

		gameCount++
	}

	expectedGames := 4 // Update this based on your fixture file
	if gameCount != expectedGames {
		t.Errorf("Expected %d games, got %d", expectedGames, gameCount)
	}
}

func TestScannerEmptyFile(t *testing.T) {
	// Create a temporary empty file
	tmpfile, err := os.CreateTemp("", "empty.pgn")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	scanner := NewScanner(tmpfile)

	if scanner.HasNext() {
		t.Error("Expected HasNext() to return false for empty file")
	}

	game, err := scanner.ScanGame()
	if err != io.EOF {
		t.Errorf("Expected EOF error for empty file, got %v", err)
	}
	if game != nil {
		t.Error("Expected nil game for empty file")
	}
}

// Additional test for sequential processing
func TestSequentialProcessing(t *testing.T) {
	file, err := os.Open(filepath.Join("fixtures", "multi_game.pgn"))
	if err != nil {
		t.Fatalf("Failed to open fixture file: %v", err)
	}
	defer file.Close()

	scanner := NewScanner(file)
	var games []*Game
	var allTokens [][]Token

	// Read all games sequentially
	for {
		game, err := scanner.ScanGame()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Failed to scan game: %v", err)
		}
		games = append(games, game)

		// Tokenize each game
		tokens, err := TokenizeGame(game)
		if err != nil {
			t.Fatalf("Failed to tokenize game: %v", err)
		}
		allTokens = append(allTokens, tokens)
	}

	// Verify we got the expected number of games
	if len(games) != 4 { // Update this based on your fixture file
		t.Errorf("Expected 4 games, got %d", len(games))
	}

	// Verify each game has tokens
	for i, tokens := range allTokens {
		if len(tokens) == 0 {
			t.Errorf("Game %d has no tokens", i+1)
		}
	}
}
