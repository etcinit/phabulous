package messages

import (
	"math/rand"

	"github.com/jacobstr/confer"
)

// GetExcuse returns an "excuse" to use for not being able to fulfill a
// request due to an unexpected error.
func GetExcuse(config *confer.Config) string {
	if config.GetBool("server.serious") {
		return "An error occurred and I was unable to fulfill your request."
	}

	excuses := []string{
		"There is some interference right now and I can't fulfill your " +
			"request.",
		"Oh noes. I messed up.",
		"Whoops. Something went wrong.",
		"1000s lines of code and I still cant get some things right.",
		"[explodes]",
		"Error: WHY U NO WORK?!",
		"OMG! It failed.",
		"such failure. such request.",
		"Oops I did it again...",
		"A cat is walking over my keywpdfahsgasgdadfk kj h",
	}

	return excuses[rand.Intn(len(excuses))]
}
