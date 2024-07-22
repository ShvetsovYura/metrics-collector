package main

import (
	"fmt"
	"os"
)

func mulfunc(i int) (int, error) {
	return i * 2, nil
}

func main() {
	fmt.Println("main func")
	os.Exit(0) // want "Вызов os.Exit в функции main пакета main"
}

func exitFunc() {
	fmt.Println("called exitFunc")
	os.Exit(0)
}
