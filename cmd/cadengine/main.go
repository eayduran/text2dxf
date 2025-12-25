package main

import (
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/eayduran/text2dxf/pkg/cadengine"
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ldate | log.Ltime)

	log.Println("=== CadEngine MCP Server v1.0.0 ===")

	s := server.NewMCPServer("CadEngine", "1.0.0")
	manager := cadengine.NewCadManager()
	cadengine.RegisterTools(s, manager)

	log.Println("Server ready - waiting for requests...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
