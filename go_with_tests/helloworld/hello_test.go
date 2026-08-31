package main

import "testing"

func TestHello(t *testing.T) {
	t.Run("Test Hello with a string", func(t *testing.T) {
		got := Hello("Chris", "")
		want := englishWordPrefix + "Chris"
		assertCorrectMessage(t, got, want)
	})

	t.Run("Test Hello without passing a string", func(t *testing.T) {
		got := Hello("", "")
		want := englishWordPrefix + "World"
		assertCorrectMessage(t, got, want)
	})

	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola, Elodie"
		assertCorrectMessage(t, got, want)
	})

	t.Run("in French", func(t *testing.T) {
		got := Hello("Elodie", "French")
		want := "Bonjour, Elodie"
		assertCorrectMessage(t, got, want)
	})
}

func assertCorrectMessage(t testing.TB, got string, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
