package main

import "testing"

func TestParseSuccess(t *testing.T) {
	data := `
	{
		"access_token": "aaaa",
		"expires_in": 3600,
		"refresh_token": "refresh",
		"scope": "test_scope",
		"token_type": "Bearer"
	}
	`

	_, err := ParseToken([]byte(data))
	if err != nil {
		t.Error("parse error")
	}
}

func TestParseFailure(t *testing.T) {
	// エラーになることのみ確認
	// 細かいケースは確認しない
	data := `
	{
		"expires_in": "aaa"
	}
	`

	_, err := ParseToken([]byte(data))
	if err == nil {
		t.Error("parse success ;want parse error")
	}
}
