package main

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
	"github.com/yourusername/cadengine/pkg/cadengine"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"CadEngine",
		"1.0.0",
	)

	// Create CAD manager
	manager := cadengine.NewCadManager()

	// Register all tools
	cadengine.RegisterTools(s, manager)

	// Start stdio server
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
