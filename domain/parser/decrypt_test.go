package parser

import (
	"testing"
)

func encNextChar(v byte) []byte {
	combined := (^v & 0xff) ^ PW_MAGIC
	return []byte{combined >> 4, combined & 0xf}
}

func encPassword(password string, length byte, toBeDeleted byte, usePWFlag bool) string {
	// Encode password into hex string format expected by Decrypt
	var parts []byte
	if usePWFlag {
		parts = append(parts, encNextChar(PW_FLAG)...)
		parts = append(parts, encNextChar(0)...) // unknown byte
		parts = append(parts, encNextChar(length)...)
	} else {
		parts = append(parts, encNextChar(length)...)
	}
	parts = append(parts, encNextChar(toBeDeleted)...)
	for i := byte(0); i < toBeDeleted; i++ {
		parts = append(parts, encNextChar(0)...)
	}
	for i := byte(0); i < length; i++ {
		parts = append(parts, encNextChar(password[i])...)
	}
	hex := ""
	for _, p := range parts {
		hex += string("0123456789ABCDEF"[p])
	}
	return hex
}

func TestDecrypt_Empty(t *testing.T) {
	if got := Decrypt("h", "u", ""); got != "" {
		t.Errorf("Decrypt() = %q, want empty", got)
	}
}

func TestDecrypt_ShortPassword(t *testing.T) {
	if got := Decrypt("h", "u", "ab"); got != "ab" {
		t.Errorf("Decrypt() = %q, want %q", got, "ab")
	}
}

func TestDecrypt_InvalidHex(t *testing.T) {
	if got := Decrypt("h", "u", "ZZZZZZZZ"); got != "ZZZZZZZZ" {
		t.Errorf("Decrypt() = %q, want %q", got, "ZZZZZZZZ")
	}
}

func TestDecrypt_Simple(t *testing.T) {
	enc := encPassword("ad", 2, 0, false)
	got := Decrypt("host", "user", enc)
	if got != "ad" {
		t.Errorf("Decrypt() = %q, want %q", got, "ad")
	}
}

func TestDecrypt_WithGarbage(t *testing.T) {
	enc := encPassword("hello", 5, 2, false)
	got := Decrypt("h", "u", enc)
	if got != "hello" {
		t.Errorf("Decrypt() = %q, want %q", got, "hello")
	}
}

func TestDecrypt_WithPWFlag(t *testing.T) {
	// PW_FLAG strips the key (username+host) prefix from decrypted result
	// key = "hu" (2 chars), so encode "hu" + "secret" = "husecret" (8 chars)
	enc := encPassword("husecret", 8, 1, true)
	got := Decrypt("h", "u", enc)
	if got != "secret" {
		t.Errorf("Decrypt() = %q, want %q", got, "secret")
	}
}

func TestDecrypt_PWFlagStripsKey(t *testing.T) {
	// With PW_FLAG, key = "user" + "host" = "userhost" (8 chars)
	// "userhostmypass" (14 chars) → strip "userhost" → "mypass" (6 chars)
	enc := encPassword("userhostmypass", 14, 0, true)
	got := Decrypt("host", "user", enc)
	if got != "mypass" {
		t.Errorf("Decrypt() = %q, want %q", got, "mypass")
	}
}

func TestDecrypt_NoPWFlagNoKeyStripping(t *testing.T) {
	// Without PW_FLAG, key stripping does NOT happen
	enc := encPassword("userhostmypass", 14, 0, false)
	got := Decrypt("host", "user", enc)
	if got != "userhostmypass" {
		t.Errorf("Decrypt() = %q, want %q", got, "userhostmypass")
	}
}

func TestDecrypt_ShortPassBytes(t *testing.T) {
	// fewer than 4 bytes after hex parsing
	got := Decrypt("h", "u", "A")
	if got != "A" {
		t.Errorf("Decrypt() = %q, want %q", got, "A")
	}
}

func TestDecrypt_TruncatedAfterFlag(t *testing.T) {
	// Has PW_FLAG but not enough bytes for length
	got := Decrypt("h", "u", "FF00")
	if got != "FF00" {
		t.Errorf("Decrypt() = %q, want %q", got, "FF00")
	}
}

func TestDecrypt_TruncatedAfterLength(t *testing.T) {
	// Has flag+length but no more bytes
	got := Decrypt("h", "u", "5E")
	if got != "5E" {
		t.Errorf("Decrypt() = %q, want %q", got, "5E")
	}
}

func TestDecrypt_ToBeDeletedTooLarge(t *testing.T) {
	// toBeDeleted > len(passbytes)/2 should return original
	// Manually construct hex where flag=length=1, toBeDeleted=3, only 1 passbyte pair left
	// flag=1: encNextChar(1) = [5, 13] → "5D"
	// toBeDeleted=3: encNextChar(3) → (^3&0xFF)^0xA3 = 0xFC^0xA3 = 0x5F → [5,15] → "5F"
	// password "x": encNextChar(120) → (^120&0xFF)^0xA3 = 0x87^0xA3 = 0x24 → [2,4] → "24"
	hex := "5D5F24"
	got := Decrypt("h", "u", hex)
	if got != hex {
		t.Errorf("Decrypt() = %q, want original %q", got, hex)
	}
}

func TestDecrypt_EmptyLength(t *testing.T) {
	enc := encPassword("", 0, 0, false)
	got := Decrypt("h", "u", enc)
	if got != "" {
		t.Errorf("Decrypt() = %q, want empty", got)
	}
}

func TestDecrypt_SpecialChars(t *testing.T) {
	password := "hello.world!"
	enc := encPassword(password, byte(len(password)), 1, false)
	got := Decrypt("h", "u", enc)
	if got != password {
		t.Errorf("Decrypt() = %q, want %q", got, password)
	}
}

func TestDecrypt_Numbers(t *testing.T) {
	password := "abc123xyz"
	enc := encPassword(password, byte(len(password)), 2, false)
	got := Decrypt("h", "u", enc)
	if got != password {
		t.Errorf("Decrypt() = %q, want %q", got, password)
	}
}
