package main

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

// Helper function to check if a character is a valid file
func isFile(ch byte) bool {
	return ch >= 'a' && ch <= 'h'
}

func isRow(ch byte) bool {
	return ch >= '1' && ch <= '8'
}
