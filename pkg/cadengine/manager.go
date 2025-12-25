package cadengine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yofu/dxf/drawing"
)

// CadManager manages the DXF drawing state
type CadManager struct {
	drawing *drawing.Drawing
}

// NewCadManager creates a new CAD manager instance
func NewCadManager() *CadManager {
	m := &CadManager{}
	m.ResetDrawing()
	return m
}

// ResetDrawing clears the canvas and sets up standard architectural layers
func (m *CadManager) ResetDrawing() error {
	m.drawing = drawing.New()
	return m.setupLayers()
}

// setupLayers defines standard CAD layers with colors
func (m *CadManager) setupLayers() error {
	layers := GetStandardLayers()

	for _, layer := range layers {
		// Add layer with color, no specific line type (nil), not set as current (false)
		_, err := m.drawing.AddLayer(
			layer.Name,
			layer.Color,
			nil,   // Line type - use default continuous
			false, // Don't set as current layer
		)
		if err != nil {
			return fmt.Errorf("failed to add layer %s: %w", layer.Name, err)
		}
	}

	return nil
}

// Save saves the drawing to a DXF file
func (m *CadManager) Save(filename string) (string, error) {
	// Ensure .dxf extension
	if !strings.HasSuffix(strings.ToLower(filename), ".dxf") {
		filename += ".dxf"
	}

	// Get absolute path
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create parent directory if it doesn't exist
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Save the file
	err = m.drawing.SaveAs(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to save DXF file: %w", err)
	}

	return absPath, nil
}

// GetDrawing returns the underlying drawing instance
func (m *CadManager) GetDrawing() *drawing.Drawing {
	return m.drawing
}
