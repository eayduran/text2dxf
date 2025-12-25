# CadEngine MCP Server

A Model Context Protocol (MCP) server that provides CAD drawing capabilities for architectural floor plans. This server enables AI assistants to generate DXF files with architectural elements like walls, rooms, doors, and furniture.

## Features

- **Professional CAD Output**: Generates industry-standard DXF files compatible with AutoCAD, LibreCAD, and other CAD software
- **Architectural Layers**: Pre-configured layer system with proper color coding
- **Drawing Primitives**: Lines, polylines, arcs, circles, and text annotations
- **Room Planning**: Perfect for creating residential and commercial floor plans

## Installation

### Prerequisites

- Go 1.22 or later - Download from [golang.org/dl](https://go.dev/dl/)

### Building from Source

```bash
# Navigate to the project directory
cd text2dxf

# Download dependencies
go mod tidy

# Build the binary
go build -o cadengine.exe ./cmd/cadengine
```

The executable `cadengine.exe` will be created in the current directory.

## Usage

This is an MCP (Model Context Protocol) server that runs as a subprocess and communicates via stdio. It's designed to be used with MCP clients like Claude Desktop or Cline (VS Code extension).

### Setting up with Claude Desktop

1. **Locate your Claude Desktop config file:**

   - Windows: `%APPDATA%\Claude\claude_desktop_config.json`
   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Linux: `~/.config/Claude/claude_desktop_config.json`

2. **Add the MCP server configuration:**

   ```json
   {
     "mcpServers": {
       "cadengine": {
         "command": "C:\\path\\to\\text2dxf\\cadengine.exe"
       }
     }
   }
   ```

   Replace `C:\\path\\to\\text2dxf\\cadengine.exe` with the actual full path to your executable.

3. **Restart Claude Desktop** to load the server.

### Setting up with Cline (VS Code)

1. Open VS Code settings for the Cline extension
2. Add the MCP server to the configuration:
   ```json
   {
     "mcpServers": {
       "cadengine": {
         "command": "C:\\path\\to\\text2dxf\\cadengine.exe"
       }
     }
   }
   ```

### Using the Server

Once configured, you can ask Claude to create CAD drawings:

**Example prompts:**

- "Create a new floor plan for a 10x8 meter room"
- "Draw a rectangular room with dimensions 4m x 3m"
- "Add a circular table in the center of the room"
- "Label this room as 'Living Room'"
- "Save the drawing as 'my_house.dxf'"

Claude will automatically use the available CAD tools to create your drawings!

### Available Tools

#### `new_project()`

Initializes a new blank CAD project. Call this at the start of each drawing session.

#### `draw_line(start, end, layer="ARCH-WALL")`

Draws a single line segment.

- `start`: [x, y] coordinates
- `end`: [x, y] coordinates
- `layer`: Layer name (optional)

#### `draw_polyline(points, closed=True, layer="ARCH-WALL")`

Draws connected line segments. Ideal for room boundaries.

- `points`: List of [x, y] coordinates, e.g., `[[0,0], [5,0], [5,5], [0,5]]`
- `closed`: Whether to close the shape (default: True)
- `layer`: Layer name (optional)

#### `draw_arc(center, radius, start_angle, end_angle, layer="ARCH-CURVE")`

Draws curved elements like door swings or rounded walls.

- `center`: [x, y] center point
- `radius`: Arc radius
- `start_angle`: Starting angle in degrees (0° = East, 90° = North)
- `end_angle`: Ending angle in degrees
- `layer`: Layer name (optional)

#### `draw_circle(center, radius, layer="ARCH-FURNITURE")`

Draws circular elements like columns or round tables.

- `center`: [x, y] center point
- `radius`: Circle radius
- `layer`: Layer name (optional)

#### `add_text(text, position, height=0.2, layer="ARCH-TEXT")`

Adds text labels for room names and annotations.

- `text`: String to display
- `position`: [x, y] insertion point
- `height`: Text height (default: 0.2 units)
- `layer`: Layer name (optional)

#### `save_file(filename="output_plan.dxf")`

Saves the drawing to a DXF file.

- `filename`: Output filename (automatically adds .dxf extension if missing)

## Layer System

The server uses a standard architectural layer structure:

| Layer          | Color   | Purpose                      |
| -------------- | ------- | ---------------------------- |
| ARCH-WALL      | White   | Main exterior/interior walls |
| ARCH-PARTITION | Gray    | Interior partition walls     |
| ARCH-DOOR      | Yellow  | Doors and door swings        |
| ARCH-WINDOW    | Cyan    | Windows                      |
| ARCH-FURNITURE | Red     | Furniture and fixtures       |
| ARCH-TEXT      | Green   | Room labels and annotations  |
| ARCH-CURVE     | Magenta | Special curved features      |

## Example: Drawing a Simple Room

```python
# Initialize new project
new_project()

# Draw room boundary (4m x 3m)
draw_polyline([[0, 0], [4, 0], [4, 3], [0, 3]], closed=True)

# Add door swing
draw_arc([0, 1.5], 0.8, 0, 90)

# Add furniture (table)
draw_circle([2, 1.5], 0.5, layer="ARCH-FURNITURE")

# Label the room
add_text("DINING ROOM", [1, 2], height=0.25)

# Save the file
save_file("my_room.dxf")
```

## Use Cases

- Residential floor plans
- Office layouts
- Apartment designs
- Commercial space planning
- Architectural documentation

## Project Structure

```
text2dxf/
├── go.mod                          # Go module definition
├── go.sum                          # Dependency checksums
├── README.md                       # This file
├── BUILD_INSTRUCTIONS.md           # Detailed build instructions
├── cadengine.exe                   # Compiled MCP server (after build)
│
├── cmd/
│   └── cadengine/
│       └── main.go                 # Application entry point
│
└── pkg/
    └── cadengine/
        ├── manager.go              # CAD state management
        ├── tools.go                # MCP tool handlers
        └── layers.go               # Layer configuration
```

## Requirements

- Go 1.23 or later
- Dependencies (automatically installed via `go mod tidy`):
  - github.com/mark3labs/mcp-go v0.43.2 (MCP server framework)
  - github.com/yofu/dxf v0.0.0-20250806094206-f3988c7f0176 (DXF file generation)

## Output Format

All drawings are saved in DXF ACAD2000 (AC1015) format, ensuring broad compatibility with:

- AutoCAD
- LibreCAD
- DraftSight
- QCAD
- BricsCAD
- And other CAD applications

## License

This is an MCP server tool for AI-assisted architectural drawing.
