package main

import "testing"

var validCronLines = []string{
	"*/15 0 1,15 * 1-5 /usr/bin/find",
	"* * * * * /usr/bin/find",
	"0 0 1 1 1 /usr/bin/find",
	"*/3 */4 */5 */6 */2 /usr/bin/find",
}

var invalidCronLines = []string{
	"/15 0 1,15 * 1-5 /usr/bin/find",
	"* * * * a /usr/bin/find",
	"0 1000 0 0 0 /usr/bin/find",
	"0 0 0 0 0 /usr/bin/find",
}

func TestValid(t *testing.T) {

	for _, l := range validCronLines {
		_, err := parseCronLine(l)
		if err != nil {
			t.Errorf("Unexpected error for input: %s, err: %v", l, err)
		}
	}
}

func TestInvalid(t *testing.T) {

	for _, l := range invalidCronLines {
		_, err := parseCronLine(l)
		if err == nil {
			t.Errorf("Expected error for input: %s, err: %v", l, err)
		}
	}
}
