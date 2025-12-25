package main

import (
	"fmt"
	"github.com/yourusername/cadengine/pkg/cadengine"
)

func main() {
	// Create CAD manager
	manager := cadengine.NewCadManager()

	// Reset drawing to start fresh
	if err := manager.ResetDrawing(); err != nil {
		fmt.Printf("Error resetting drawing: %v\n", err)
		return
	}
	fmt.Println("New project started")

	// Change to wall layer
	if err := manager.GetDrawing().ChangeLayer("ARCH-WALL"); err != nil {
		fmt.Printf("Error changing layer: %v\n", err)
		return
	}

	// Draw a 3x5 meter rectangular room using polyline
	// Define the 4 corners of the room
	vertices := [][]float64{
		{0, 0, 0},    // Bottom-left corner
		{3, 0, 0},    // Bottom-right corner (3 meters wide)
		{3, 5, 0},    // Top-right corner (5 meters tall)
		{0, 5, 0},    // Top-left corner
	}

	// Draw closed polyline (rectangle)
	_, err := manager.GetDrawing().LwPolyline(true, vertices...)
	if err != nil {
		fmt.Printf("Error drawing room: %v\n", err)
		return
	}
	fmt.Println("Room boundary added (3x5 meters)")

	// Add a text label
	if err := manager.GetDrawing().ChangeLayer("ARCH-TEXT"); err != nil {
		fmt.Printf("Error changing to text layer: %v\n", err)
		return
	}

	_, err = manager.GetDrawing().Text("3x5m Room", 1.5, 2.5, 0, 0.3)
	if err != nil {
		fmt.Printf("Error adding text: %v\n", err)
		return
	}
	fmt.Println("Text label added")

	// Save the drawing
	path, err := manager.Save("room_3x5.dxf")
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		return
	}

	fmt.Printf("Success! File saved to: %s\n", path)
}
