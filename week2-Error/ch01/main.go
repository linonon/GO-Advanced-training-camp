package main

import (
	"errors"
	"fmt"
)

type errorString string

func (e errorString) Error() string {
	return string(e)
}

func New(text string) error {
	return errorString(text)
}

var ErrNameType = New("EOF")
var ErrStructType = errors.New("EOF")

func main() {
	if ErrNameType == New("EOF") {
		fmt.Println("Named Type Error") // print
		// 不要使用判定字符串是否相等的方法來判斷error
	}

	if ErrStructType == errors.New("EOF") {
		fmt.Println("Struct Type Error") // not print
	}
}
