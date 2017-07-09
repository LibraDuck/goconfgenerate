package main

import "strconv"

// convert：ab_cd_ef_  ->  AbCdEf
func ParseField(s string) string {
	s = CapitalizeCertainLetter(s, 0)
	s = UnderlineConvertCapitalize(s)
	return s
}

// 字母转大写	//
func CapitalizeCertainLetter(s string, index int) string {
	if len(s) == 0 {
		return s
	}
	field := []byte(s)
	if field[index] >= 97 && field[index] <= 122 {
		field[index] -= 32
	}
	return string(field)
}

// convert：ab_cd_ef_  ->  abCdEf
func UnderlineConvertCapitalize(s string) string {
	// 未排除纯下划线字符串的特殊情况
	if len(s) == 0 {
		return s
	}
	// 转大写字母
	for i, v := range []byte(s) {
		if v == '_' && i < len(s) - 1 {
			s = CapitalizeCertainLetter(s, i + 1)
		}
	}

	// 去结尾下划线
	field := []byte(s)
	for {
		if field[len(field) - 1] == '_' {
			field = append(field[:len(field) - 1], []byte{}...)
		} else {
			break
		}
	}

	return string(field[:len(field)])
}

// 提取引号内的[]byte
func ExtractQuotesBytes(slice []byte) []byte {
	var state, begin, end int
	//log.Println(string(slice))
	for i, v := range slice {
		if v == '"' {
			switch state {
			case 0: begin = i
			case 1: end = i
			}
			state = 1
		}
	}
	if state != 1 {
		return []byte{}
	}

	slice = append(slice[:end], []byte{}...)
	slice = append(slice[:0], slice[begin + 1:]...)

	return slice
}

// 将 ascii 1 转换成 "1"	 //
func ByteToString(b byte) string {
	return strconv.Itoa(int(b - 'A' + 1))
}
