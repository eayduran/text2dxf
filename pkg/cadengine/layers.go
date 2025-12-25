package cadengine

import "github.com/yofu/dxf/color"

// LayerDefinition represents a CAD layer configuration
type LayerDefinition struct {
	Name  string
	Color color.ColorNumber
}

// GetStandardLayers returns the standard architectural layer definitions
func GetStandardLayers() []LayerDefinition {
	return []LayerDefinition{
		// Main Walls - White (7 in AutoCAD Color Index)
		{Name: "ARCH-WALL", Color: color.White},

		// Inner Walls - Gray (8 in AutoCAD Color Index)
		{Name: "ARCH-PARTITION", Color: color.Grey128},

		// Doors - Yellow (2 in AutoCAD Color Index)
		{Name: "ARCH-DOOR", Color: color.Yellow},

		// Windows - Cyan (4 in AutoCAD Color Index)
		{Name: "ARCH-WINDOW", Color: color.Cyan},

		// Furniture - Red (1 in AutoCAD Color Index)
		{Name: "ARCH-FURNITURE", Color: color.Red},

		// Text - Green (3 in AutoCAD Color Index)
		{Name: "ARCH-TEXT", Color: color.Green},

		// Curves - Magenta (6 in AutoCAD Color Index)
		{Name: "ARCH-CURVE", Color: color.Magenta},
	}
}

// Layer name constants for convenience
const (
	LayerWall      = "ARCH-WALL"
	LayerPartition = "ARCH-PARTITION"
	LayerDoor      = "ARCH-DOOR"
	LayerWindow    = "ARCH-WINDOW"
	LayerFurniture = "ARCH-FURNITURE"
	LayerText      = "ARCH-TEXT"
	LayerCurve     = "ARCH-CURVE"
)
