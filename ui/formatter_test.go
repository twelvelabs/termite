package ui

import (
	"testing"

	"github.com/mgutz/ansi" //cspell: disable-line
	"github.com/stretchr/testify/assert"
)

// cspell: ignoreRegExp (Magenta|Cyan|Red|Yellow|Blue|Green|Gray|Bold|Underline)f
func TestFormatter(t *testing.T) {
	formatter := NewFormatter(true)

	var start string
	var stop = ansi.Reset

	start = ansi.ColorCode("default+b")
	assert.Equal(t, (start + "text" + stop), formatter.Bold("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Boldf("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Color("text", "default+b"))

	start = ansi.ColorCode("default+u")
	assert.Equal(t, (start + "text" + stop), formatter.Underline("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Underlinef("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Color("text", "default+u"))

	start = ansi.ColorCode("red")
	assert.Equal(t, (start + "text" + stop), formatter.Red("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Redf("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Color("text", "red"))
	assert.Equal(t, (start + "✖" + stop), formatter.FailureIcon())

	start = ansi.ColorCode("yellow")
	assert.Equal(t, (start + "text" + stop), formatter.Yellow("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Yellowf("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Color("text", "yellow"))
	assert.Equal(t, (start + "!" + stop), formatter.WarningIcon())

	start = ansi.ColorCode("green")
	assert.Equal(t, (start + "text" + stop), formatter.Green("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Greenf("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Color("text", "green"))
	assert.Equal(t, (start + "✓" + stop), formatter.SuccessIcon())

	start = ansi.ColorCode("black+h")
	assert.Equal(t, (start + "text" + stop), formatter.Gray("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Grayf("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Color("text", "black+h"))
	assert.Equal(t, (start + "•" + stop), formatter.InfoIcon())

	start = ansi.ColorCode("magenta")
	assert.Equal(t, (start + "text" + stop), formatter.Magenta("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Magentaf("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Color("text", "magenta"))

	start = ansi.ColorCode("cyan")
	assert.Equal(t, (start + "text" + stop), formatter.Cyan("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Cyanf("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Color("text", "cyan"))

	start = ansi.ColorCode("blue")
	assert.Equal(t, (start + "text" + stop), formatter.Blue("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Bluef("text"))
	assert.Equal(t, (start + "text" + stop), formatter.Color("text", "blue"))
}

func TestFormatterIsNoopWhenDisabled(t *testing.T) {
	formatter := NewFormatter(false)

	assert.Equal(t, "text", formatter.Bold("text"))
	assert.Equal(t, "text", formatter.Boldf("text"))
	assert.Equal(t, "text", formatter.Color("text", "default+b"))

	assert.Equal(t, "text", formatter.Underline("text"))
	assert.Equal(t, "text", formatter.Underlinef("text"))
	assert.Equal(t, "text", formatter.Color("text", "default+u"))

	assert.Equal(t, "text", formatter.Red("text"))
	assert.Equal(t, "text", formatter.Redf("text"))
	assert.Equal(t, "text", formatter.Color("text", "red"))
	assert.Equal(t, "✖", formatter.FailureIcon())

	assert.Equal(t, "text", formatter.Yellow("text"))
	assert.Equal(t, "text", formatter.Yellowf("text"))
	assert.Equal(t, "text", formatter.Color("text", "yellow"))
	assert.Equal(t, "!", formatter.WarningIcon())

	assert.Equal(t, "text", formatter.Green("text"))
	assert.Equal(t, "text", formatter.Greenf("text"))
	assert.Equal(t, "text", formatter.Color("text", "green"))
	assert.Equal(t, "✓", formatter.SuccessIcon())

	assert.Equal(t, "text", formatter.Gray("text"))
	assert.Equal(t, "text", formatter.Grayf("text"))
	assert.Equal(t, "text", formatter.Color("text", "black+h"))
	assert.Equal(t, "•", formatter.InfoIcon())

	assert.Equal(t, "text", formatter.Magenta("text"))
	assert.Equal(t, "text", formatter.Magentaf("text"))
	assert.Equal(t, "text", formatter.Color("text", "magenta"))

	assert.Equal(t, "text", formatter.Cyan("text"))
	assert.Equal(t, "text", formatter.Cyanf("text"))
	assert.Equal(t, "text", formatter.Color("text", "cyan"))

	assert.Equal(t, "text", formatter.Blue("text"))
	assert.Equal(t, "text", formatter.Bluef("text"))
	assert.Equal(t, "text", formatter.Color("text", "blue"))
}
