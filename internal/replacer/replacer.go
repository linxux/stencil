package replacer

import (
	"bytes"
	"os"
	"regexp"
	"strings"
)

// Replacer handles keyword replacement in content and paths
type Replacer struct {
	variables map[string]string
}

// NewReplacer creates a new Replacer with the given variables
func NewReplacer(variables map[string]string) *Replacer {
	return &Replacer{
		variables: variables,
	}
}

// ReplaceInContent replaces variables in file content
func (r *Replacer) ReplaceInContent(content []byte) []byte {
	result := content
	for key, value := range r.variables {
		// Replace {{key}} format
		pattern := []byte("{{" + key + "}}")
		result = bytes.ReplaceAll(result, pattern, []byte(value))

		// Replace <<key>> format
		pattern2 := []byte("<<" + key + ">>")
		result = bytes.ReplaceAll(result, pattern2, []byte(value))

		// Replace __key__ format
		pattern3 := []byte("__" + key + "__")
		result = bytes.ReplaceAll(result, pattern3, []byte(value))

		// Replace %key% format
		pattern4 := []byte("%" + key + "%")
		result = bytes.ReplaceAll(result, pattern4, []byte(value))
	}
	return result
}

// ReplaceInPath replaces variables in file or directory paths
func (r *Replacer) ReplaceInPath(path string) string {
	result := path
	for key, value := range r.variables {
		// Replace {{key}} format
		result = strings.ReplaceAll(result, "{{"+key+"}}", value)

		// Replace <<key>> format
		result = strings.ReplaceAll(result, "<<"+key+">>", value)

		// Replace __key__ format (common in folder names)
		result = strings.ReplaceAll(result, "__"+key+"__", value)

		// Replace %key% format
		result = strings.ReplaceAll(result, "%"+key+"%", value)
	}
	return result
}

// ExtractVariablesFromFile extracts variables from file content
func ExtractVariablesFromFile(content []byte) []string {
	variables := make(map[string]bool)

	// Pattern for {{var}}
	pattern1 := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches := pattern1.FindAllSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			variables[string(match[1])] = true
		}
	}

	// Pattern for <<var>>
	pattern2 := regexp.MustCompile(`<<([^>]+)>>`)
	matches = pattern2.FindAllSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			variables[string(match[1])] = true
		}
	}

	// Pattern for __var__
	pattern3 := regexp.MustCompile(`__([A-Za-z0-9_]+)__`)
	matches = pattern3.FindAllSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			variables[string(match[1])] = true
		}
	}

	// Pattern for %var%
	pattern4 := regexp.MustCompile(`%([A-Za-z0-9_]+)%`)
	matches = pattern4.FindAllSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			variables[string(match[1])] = true
		}
	}

	result := make([]string, 0, len(variables))
	for v := range variables {
		result = append(result, v)
	}
	return result
}

// ExtractVariablesFromPath extracts variables from a path
func ExtractVariablesFromPath(path string) []string {
	variables := make(map[string]bool)

	// Pattern for {{var}}
	pattern1 := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches := pattern1.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		if len(match) > 1 {
			variables[match[1]] = true
		}
	}

	// Pattern for <<var>>
	pattern2 := regexp.MustCompile(`<<([^>]+)>>`)
	matches = pattern2.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		if len(match) > 1 {
			variables[match[1]] = true
		}
	}

	// Pattern for __var__
	pattern3 := regexp.MustCompile(`__([A-Za-z0-9_]+)__`)
	matches = pattern3.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		if len(match) > 1 {
			variables[match[1]] = true
		}
	}

	// Pattern for %var%
	pattern4 := regexp.MustCompile(`%([A-Za-z0-9_]+)%`)
	matches = pattern4.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		if len(match) > 1 {
			variables[match[1]] = true
		}
	}

	result := make([]string, 0, len(variables))
	for v := range variables {
		result = append(result, v)
	}
	return result
}

// IsBinaryFile checks if a file is binary (should skip content replacement)
func IsBinaryFile(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read first 512 bytes to determine file type
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return false
	}
	buffer = buffer[:n]

	// Check for null byte (common in binary files)
	for _, b := range buffer {
		if b == 0 {
			return true
		}
	}

	return false
}
