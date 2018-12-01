package main

import (
	"fmt"
	"regexp"
)

// func normalize(phone string) string {
// 	var buf bytes.Buffer

// 	// when iterate, string will be []rune
// 	for _, ch := range phone {
// 		// condition to check if ch >= 0 && ch <= 9
// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}
// 	}

// 	fmt.Println(buf.String())

// 	return buf.String()
// }

func normalize(phone string) string {
	// compile regex to make sure is not error
	// if error it will panic
	// \D -> non digit
	re := regexp.MustCompile("\\D")

	// replace non digit with ""
	return re.ReplaceAllString(phone, "")
}

func main() {
	res := normalize("12345678")
	fmt.Println(res)
}
