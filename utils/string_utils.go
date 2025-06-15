package utils

import "strings"

type ParsedName struct {
	FirstName string
	LastName  string
}

func ParseFullName(fullName string) ParsedName {
	if fullName == "" {
		return ParsedName{FirstName: "", LastName: ""}
	}

	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return ParsedName{FirstName: "", LastName: ""}
	}

	if len(parts) == 1 {
		return ParsedName{FirstName: parts[0], LastName: ""}
	}

	return ParsedName{
		FirstName: parts[0],
		LastName:  strings.Join(parts[1:], " "),
	}
}
