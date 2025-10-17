package utils

import "strings"

func InjectArg(arg string, template string) string {
	return strings.Replace(template, "{arg}", arg, 1)
}
