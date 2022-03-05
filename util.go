package main

import (
	"flag"
	"fmt"
	"strings"
)

var Debug = flag.Bool("debug", false, "Enables debug mode.")

func ModDebugf(module string, format string, a ...interface{}) {
	if !*Debug {
		return
	}

	info := fmt.Sprintf(format, a...)

	fmt.Printf("[DEBUG] [%s] %v\n", strings.ToUpper(module), info)
}

func ModPrintf(module string, format string, a ...interface{}) {
	info := fmt.Sprintf(format, a...)

	fmt.Printf("[%s] %v\n", strings.ToUpper(module), info)
}

func normalizeTagsSlice(t []string) string {
	return normalizeTags(strings.Join(t, " "))
}

func normalizeTags(t string) string {
	t = strings.ReplaceAll(t, " ", "+")
	t = strings.ReplaceAll(t, ":", "%3A")
	return t
}
