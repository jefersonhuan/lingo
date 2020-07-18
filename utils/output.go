package utils

var colors = map[string]string{
	"reset": "\033[0m",

	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"purple": "\033[35m",
	"cyan":   "\033[36m",
	"white":  "\033[37m",
}

func ColorfulString(color string, content string) string {
	code, ok := colors[color]

	if !ok {
		return content
	} else {
		return code + string(content) + colors["reset"]
	}
}
