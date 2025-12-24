package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Prompter handles interactive user prompts
type Prompter struct {
	reader *bufio.Reader
}

// NewPrompter creates a new Prompter instance
func NewPrompter() *Prompter {
	return &Prompter{
		reader: bufio.NewReader(os.Stdin),
	}
}

// PromptForValues prompts the user for variable values
func (p *Prompter) PromptForValues(variables map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	fmt.Println("\n=== Interactive Variable Prompt ===")
	fmt.Println("Please provide values for the following variables:")
	fmt.Println()

	// Convert to slice for ordered prompting
	varKeys := make([]string, 0, len(variables))
	for k := range variables {
		varKeys = append(varKeys, k)
	}

	for i, key := range varKeys {
		defaultValue := variables[key]
		prompt := fmt.Sprintf("[%d/%d] %s", i+1, len(varKeys), key)

		if defaultValue != "" {
			prompt += fmt.Sprintf(" (default: %s)", defaultValue)
		}
		prompt += ": "

		fmt.Print(prompt)
		input, err := p.reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)

		// Use default value if input is empty
		if input == "" && defaultValue != "" {
			input = defaultValue
		}

		result[key] = input
	}

	return result, nil
}

// PromptForConfirmation prompts the user for confirmation
func (p *Prompter) PromptForConfirmation(message string) (bool, error) {
	fmt.Printf("\n%s [y/N]: ", message)

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes", nil
}

// PromptForChoice prompts the user to select from a list of choices
func (p *Prompter) PromptForChoice(message string, choices []string, defaultIndex int) (int, error) {
	fmt.Printf("\n%s\n", message)

	for i, choice := range choices {
		prefix := " "
		if i == defaultIndex {
			prefix = "*"
		}
		fmt.Printf("  %s [%d] %s\n", prefix, i+1, choice)
	}

	fmt.Printf("\nSelect choice [1-%d]", len(choices))
	if defaultIndex >= 0 {
		fmt.Printf(" (default: %d)", defaultIndex+1)
	}
	fmt.Print(": ")

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return -1, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)

	// Use default if input is empty
	if input == "" && defaultIndex >= 0 {
		return defaultIndex, nil
	}

	// Parse input
	choice, err := strconv.Atoi(input)
	if err != nil {
		return -1, fmt.Errorf("invalid input: %s", input)
	}

	if choice < 1 || choice > len(choices) {
		return -1, fmt.Errorf("choice out of range: %d", choice)
	}

	return choice - 1, nil
}

// PromptForString prompts the user for a string input
func (p *Prompter) PromptForString(message, defaultValue string) (string, error) {
	prompt := message
	if defaultValue != "" {
		prompt += fmt.Sprintf(" (default: %s)", defaultValue)
	}
	prompt += ": "

	fmt.Print(prompt)

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)

	if input == "" {
		input = defaultValue
	}

	return input, nil
}
