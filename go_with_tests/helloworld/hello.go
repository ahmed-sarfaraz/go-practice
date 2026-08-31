package main

import "fmt"

const spanish = "Spanish"
const french = "French"
const englishWordPrefix = "Hello, "
const spanishWordPrefix = "Hola, "
const frenchWordPrefix = "Bonjour, "

func Hello(name string, lang string) string {
	if len(name) == 0 {
		name = "World"
	}

	return greetingPrefix(lang) + name
}

func greetingPrefix(lang string) string {
	prefix := englishWordPrefix
	switch lang {
	case spanish:
		prefix = spanishWordPrefix
	case french:
		prefix = frenchWordPrefix
	}
	return prefix
}

func main() {
	fmt.Println(Hello("Chris", ""))
	fmt.Println(Hello("Elodie", spanish))
}
