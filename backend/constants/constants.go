package constants

import "os"

var LanguageMap = map[string]string{
	"c++":     "54",
	"c":       "50",
	"golang":  "60",
	"python3": "71",
}

const Judge0Endpoint = "http://localhost:2358"
const PathSeparatorStr = string(os.PathSeparator)

const DefaultSourceCode = `
#include<iostream>

int main() {
	return 0;
}
`
