package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/linxux/stencil/config"
	"github.com/linxux/stencil/internal/replacer"
)

// Generator handles the template generation process
type Generator struct {
	cfg      *config.Config
	replacer *replacer.Replacer
}

// NewGenerator creates a new Generator instance
func NewGenerator(cfg *config.Config) *Generator {
	return &Generator{
		cfg:      cfg,
		replacer: replacer.NewReplacer(cfg.Variables, cfg.Formats),
	}
}

// Generate generates the project from template
func (g *Generator) Generate() error {
	// Validate template directory
	if _, err := os.Stat(g.cfg.TemplateDir); os.IsNotExist(err) {
		return fmt.Errorf("template directory does not exist: %s", g.cfg.TemplateDir)
	}

	// Create output directory
	if err := os.MkdirAll(g.cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Walk through template directory
	return filepath.Walk(g.cfg.TemplateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from template directory
		relPath, err := filepath.Rel(g.cfg.TemplateDir, path)
		if err != nil {
			return err
		}

		// Skip the template directory itself
		if relPath == "." {
			return nil
		}

		// Replace variables in path
		targetPath := filepath.Join(g.cfg.OutputDir, g.replacer.ReplaceInPath(relPath))

		if info.IsDir() {
			// Create directory
			if g.cfg.DryRun {
				fmt.Printf("[DRY RUN] Would create directory: %s\n", targetPath)
				return nil
			}
			return os.MkdirAll(targetPath, info.Mode())
		}

		// Process file
		return g.processFile(path, targetPath, info)
	})
}

// processFile processes a single template file
func (g *Generator) processFile(sourcePath, targetPath string, info os.FileInfo) error {
	// Read source file
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Check if file is binary
	isBinary := replacer.IsBinaryFile(sourcePath)

	if isBinary {
		// Copy binary file as-is
		if g.cfg.DryRun {
			fmt.Printf("[DRY RUN] Would copy binary file: %s -> %s\n", sourcePath, targetPath)
			return nil
		}

		// Ensure target directory exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		return g.copyFile(sourcePath, targetPath)
	}

	// Read content
	content, err := io.ReadAll(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	// Replace variables in content
	newContent := g.replacer.ReplaceInContent(content)

	// Write target file
	if g.cfg.DryRun {
		fmt.Printf("[DRY RUN] Would create file: %s\n", targetPath)
		fmt.Printf("[DRY RUN] Content preview (first 200 chars): %s\n",
			truncateString(string(newContent), 200))
		return nil
	}

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer targetFile.Close()

	_, err = targetFile.Write(newContent)
	if err != nil {
		return fmt.Errorf("failed to write target file: %w", err)
	}

	return nil
}

// copyFile copies a file from source to destination
func (g *Generator) copyFile(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// ExtractVariables extracts all variables from the template
func (g *Generator) ExtractVariables() (map[string]string, error) {
	variables := make(map[string]bool)

	err := filepath.Walk(g.cfg.TemplateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Extract variables from directory names
			relPath, err := filepath.Rel(g.cfg.TemplateDir, path)
			if err != nil {
				return err
			}
			if relPath != "." {
				for _, v := range replacer.ExtractVariablesFromPath(relPath, g.cfg.Formats) {
					variables[v] = true
				}
			}
			return nil
		}

		// Extract variables from file names
		relPath, err := filepath.Rel(g.cfg.TemplateDir, path)
		if err != nil {
			return err
		}
		for _, v := range replacer.ExtractVariablesFromPath(relPath, g.cfg.Formats) {
			variables[v] = true
		}

		// Extract variables from file content
		if !replacer.IsBinaryFile(path) {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			for _, v := range replacer.ExtractVariablesFromFile(content, g.cfg.Formats) {
				variables[v] = true
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to map with empty values
	result := make(map[string]string)
	for v := range variables {
		result[v] = ""
	}

	return result, nil
}

// SetVariables updates the generator's variables
func (g *Generator) SetVariables(variables map[string]string) {
	g.cfg.Variables = variables
	g.replacer = replacer.NewReplacer(variables, g.cfg.Formats)
}

// TemplateDir returns the template directory path
func (g *Generator) TemplateDir() string {
	return g.cfg.TemplateDir
}

// OutputDir returns the output directory path
func (g *Generator) OutputDir() string {
	return g.cfg.OutputDir
}

// SkipConfirm returns whether to skip confirmation
func (g *Generator) SkipConfirm() bool {
	return g.cfg.SkipConfirm
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
