package parser

import (
	"strconv"
)

const (
	PW_MAGIC = 0xA3
	PW_FLAG  = 0xFF
)

func Decrypt(host, username, password string) string {
	if password == "" {
		return ""
	}

	if len(password) < 4 {
		return password
	}

	key := username + host
	passbytes := []byte{}
	for i := 0; i < len(password); i++ {
		val, err := strconv.ParseInt(string(password[i]), 16, 8)
		if err != nil {
			return password
		}
		passbytes = append(passbytes, byte(val))
	}

	if len(passbytes) < 4 {
		return password
	}

	var flag byte
	flag, passbytes = decNextChar(passbytes)
	var length byte = 0
	if flag == PW_FLAG {
		if len(passbytes) < 2 {
			return password
		}
		_, passbytes = decNextChar(passbytes)
		if len(passbytes) < 1 {
			return password
		}
		length, passbytes = decNextChar(passbytes)
	} else {
		length = flag
	}

	if len(passbytes) < 1 {
		return password
	}

	toBeDeleted, passbytes := decNextChar(passbytes)
	if int(toBeDeleted) > len(passbytes)/2 {
		return password
	}
	passbytes = passbytes[toBeDeleted*2:]

	clearpass := ""
	var (
		i   byte
		val byte
	)
	for i = 0; i < length; i++ {
		if len(passbytes) < 2 {
			break
		}
		val, passbytes = decNextChar(passbytes)
		clearpass += string(val)
	}

	if flag == PW_FLAG && len(clearpass) > len(key) {
		clearpass = clearpass[len(key):]
	}
	return clearpass
}

func decNextChar(passbytes []byte) (byte, []byte) {
	if len(passbytes) <= 0 {
		return 0, passbytes
	}
	a := passbytes[0]
	b := passbytes[1]
	passbytes = passbytes[2:]
	return ^(((a << 4) + b) ^ PW_MAGIC) & 0xff, passbytes
}