package cadengine

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers all CAD tools with the MCP server
func RegisterTools(s *server.MCPServer, manager *CadManager) {
	// Tool 1: new_project
	s.AddTool(
		mcp.NewTool("new_project",
			mcp.WithDescription("Clears the current session and starts a blank CAD project. Always call this at the beginning of a new request."),
		),
		newProjectHandler(manager),
	)

	// Tool 2: draw_line
	s.AddTool(
		mcp.NewTool("draw_line",
			mcp.WithDescription("Draws a simple line between two points."),
			mcp.WithArray("start",
				mcp.Description("[x, y] coordinates"),
				mcp.Required(),
			),
			mcp.WithArray("end",
				mcp.Description("[x, y] coordinates"),
				mcp.Required(),
			),
			mcp.WithString("layer",
				mcp.Description("Layer name (e.g., 'ARCH-WALL', 'ARCH-FURNITURE')"),
				mcp.DefaultString("ARCH-WALL"),
			),
		),
		drawLineHandler(manager),
	)

	// Tool 3: draw_polyline
	s.AddTool(
		mcp.NewTool("draw_polyline",
			mcp.WithDescription("Draws a sequence of connected lines (Polyline). Best for room boundaries."),
			mcp.WithArray("points",
				mcp.Description("List of [x, y] coordinates. e.g., [[0,0], [5,0], [5,5], [0,5]]"),
				mcp.Required(),
			),
			mcp.WithBoolean("closed",
				mcp.Description("If True, connects the last point back to the first"),
				mcp.DefaultBool(true),
			),
			mcp.WithString("layer",
				mcp.Description("Layer name"),
				mcp.DefaultString("ARCH-WALL"),
			),
		),
		drawPolylineHandler(manager),
	)

	// Tool 4: draw_arc
	s.AddTool(
		mcp.NewTool("draw_arc",
			mcp.WithDescription("Draws an Arc (curve). Used for curved walls, door swings, or circular designs."),
			mcp.WithArray("center",
				mcp.Description("[x, y] center point of the arc"),
				mcp.Required(),
			),
			mcp.WithNumber("radius",
				mcp.Description("Distance from center to edge"),
				mcp.Required(),
			),
			mcp.WithNumber("start_angle",
				mcp.Description("Starting angle in degrees (0 is East, 90 is North)"),
				mcp.Required(),
			),
			mcp.WithNumber("end_angle",
				mcp.Description("Ending angle in degrees"),
				mcp.Required(),
			),
			mcp.WithString("layer",
				mcp.Description("Layer name"),
				mcp.DefaultString("ARCH-CURVE"),
			),
		),
		drawArcHandler(manager),
	)

	// Tool 5: draw_circle
	s.AddTool(
		mcp.NewTool("draw_circle",
			mcp.WithDescription("Draws a full Circle. Useful for columns or round tables."),
			mcp.WithArray("center",
				mcp.Description("[x, y] center point"),
				mcp.Required(),
			),
			mcp.WithNumber("radius",
				mcp.Description("Circle radius"),
				mcp.Required(),
			),
			mcp.WithString("layer",
				mcp.Description("Layer name"),
				mcp.DefaultString("ARCH-FURNITURE"),
			),
		),
		drawCircleHandler(manager),
	)

	// Tool 6: add_text
	s.AddTool(
		mcp.NewTool("add_text",
			mcp.WithDescription("Adds text to the drawing. Use this for Room Names and Area labels."),
			mcp.WithString("text",
				mcp.Description("The string to display"),
				mcp.Required(),
			),
			mcp.WithArray("position",
				mcp.Description("[x, y] location"),
				mcp.Required(),
			),
			mcp.WithNumber("height",
				mcp.Description("Text size (default 0.2 units)"),
				mcp.DefaultNumber(0.2),
			),
			mcp.WithString("layer",
				mcp.Description("Layer name"),
				mcp.DefaultString("ARCH-TEXT"),
			),
		),
		addTextHandler(manager),
	)

	// Tool 7: save_file
	s.AddTool(
		mcp.NewTool("save_file",
			mcp.WithDescription("Saves the final drawing to a DXF file on the computer."),
			mcp.WithString("filename",
				mcp.Description("Output filename"),
				mcp.DefaultString("output_plan.dxf"),
			),
		),
		saveFileHandler(manager),
	)
}

// Handler implementations

func newProjectHandler(manager *CadManager) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		err := manager.ResetDrawing()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to reset drawing: %v", err)), nil
		}
		return mcp.NewToolResultText("New project started. Canvas is empty."), nil
	}
}

func drawLineHandler(manager *CadManager) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments format"), nil
		}

		// Extract parameters
		start, err := getFloatArray(args, "start", 2)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		end, err := getFloatArray(args, "end", 2)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		layer := getStringWithDefault(args, "layer", "ARCH-WALL")

		// Change to the specified layer
		if err := manager.drawing.ChangeLayer(layer); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to change layer: %v", err)), nil
		}

		// Draw the line (z=0 for 2D)
		_, err = manager.drawing.Line(start[0], start[1], 0, end[0], end[1], 0)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to draw line: %v", err)), nil
		}

		return mcp.NewToolResultText("Line added."), nil
	}
}

func drawPolylineHandler(manager *CadManager) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments format"), nil
		}

		// Extract points array
		pointsRaw, ok := args["points"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: points"), nil
		}

		points, err := parsePointsArray(pointsRaw)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid points format: %v", err)), nil
		}

		closed := getBoolWithDefault(args, "closed", true)
		layer := getStringWithDefault(args, "layer", "ARCH-WALL")

		// Change to the specified layer
		if err := manager.drawing.ChangeLayer(layer); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to change layer: %v", err)), nil
		}

		// Convert to []float64 format (x, y, z) for yofu/dxf
		vertices := make([][]float64, len(points))
		for i, p := range points {
			vertices[i] = []float64{p[0], p[1], 0}
		}

		// Draw polyline
		_, err = manager.drawing.LwPolyline(closed, vertices...)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to draw polyline: %v", err)), nil
		}

		return mcp.NewToolResultText("Polyline added."), nil
	}
}

func drawArcHandler(manager *CadManager) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments format"), nil
		}

		center, err := getFloatArray(args, "center", 2)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		radius := getNumberWithDefault(args, "radius", 0)
		startAngle := getNumberWithDefault(args, "start_angle", 0)
		endAngle := getNumberWithDefault(args, "end_angle", 0)
		layer := getStringWithDefault(args, "layer", "ARCH-CURVE")

		if radius <= 0 {
			return mcp.NewToolResultError("radius must be greater than 0"), nil
		}

		// Change to the specified layer
		if err := manager.drawing.ChangeLayer(layer); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to change layer: %v", err)), nil
		}

		// Draw arc (z=0 for 2D)
		_, err = manager.drawing.Arc(center[0], center[1], 0, radius, startAngle, endAngle)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to draw arc: %v", err)), nil
		}

		return mcp.NewToolResultText("Arc added."), nil
	}
}

func drawCircleHandler(manager *CadManager) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments format"), nil
		}

		center, err := getFloatArray(args, "center", 2)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		radius := getNumberWithDefault(args, "radius", 0)
		layer := getStringWithDefault(args, "layer", "ARCH-FURNITURE")

		if radius <= 0 {
			return mcp.NewToolResultError("radius must be greater than 0"), nil
		}

		// Change to the specified layer
		if err := manager.drawing.ChangeLayer(layer); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to change layer: %v", err)), nil
		}

		// Draw circle (z=0 for 2D)
		_, err = manager.drawing.Circle(center[0], center[1], 0, radius)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to draw circle: %v", err)), nil
		}

		return mcp.NewToolResultText("Circle added."), nil
	}
}

func addTextHandler(manager *CadManager) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments format"), nil
		}

		textVal, ok := args["text"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: text"), nil
		}

		text, ok := textVal.(string)
		if !ok {
			return mcp.NewToolResultError("Parameter text must be a string"), nil
		}

		position, err := getFloatArray(args, "position", 2)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		height := getNumberWithDefault(args, "height", 0.2)
		layer := getStringWithDefault(args, "layer", "ARCH-TEXT")

		// Change to the specified layer
		if err := manager.drawing.ChangeLayer(layer); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to change layer: %v", err)), nil
		}

		// Add text (z=0 for 2D)
		_, err = manager.drawing.Text(text, position[0], position[1], 0, height)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to add text: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Text '%s' added.", text)), nil
	}
}

func saveFileHandler(manager *CadManager) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments format"), nil
		}

		filename := getStringWithDefault(args, "filename", "output_plan.dxf")

		path, err := manager.Save(filename)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to save file: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("File saved successfully at: %s", path)), nil
	}
}

// Helper functions

func getFloatArray(args map[string]interface{}, key string, expectedLen int) ([]float64, error) {
	val, ok := args[key]
	if !ok {
		return nil, fmt.Errorf("missing required parameter: %s", key)
	}

	arr, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("parameter %s must be an array", key)
	}

	if len(arr) != expectedLen {
		return nil, fmt.Errorf("parameter %s must have exactly %d elements", key, expectedLen)
	}

	result := make([]float64, expectedLen)
	for i, v := range arr {
		switch num := v.(type) {
		case float64:
			result[i] = num
		case int:
			result[i] = float64(num)
		case int64:
			result[i] = float64(num)
		default:
			return nil, fmt.Errorf("parameter %s[%d] must be a number", key, i)
		}
	}

	return result, nil
}

func parsePointsArray(pointsRaw interface{}) ([][2]float64, error) {
	arr, ok := pointsRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("points must be an array")
	}

	points := make([][2]float64, len(arr))
	for i, pointRaw := range arr {
		pointArr, ok := pointRaw.([]interface{})
		if !ok {
			return nil, fmt.Errorf("point %d must be an array", i)
		}

		if len(pointArr) != 2 {
			return nil, fmt.Errorf("point %d must have exactly 2 elements", i)
		}

		for j := 0; j < 2; j++ {
			switch num := pointArr[j].(type) {
			case float64:
				points[i][j] = num
			case int:
				points[i][j] = float64(num)
			case int64:
				points[i][j] = float64(num)
			default:
				return nil, fmt.Errorf("point %d[%d] must be a number", i, j)
			}
		}
	}

	return points, nil
}

func getStringWithDefault(args map[string]interface{}, key string, defaultVal string) string {
	val, ok := args[key]
	if !ok {
		return defaultVal
	}

	str, ok := val.(string)
	if !ok {
		return defaultVal
	}

	return str
}

func getBoolWithDefault(args map[string]interface{}, key string, defaultVal bool) bool {
	val, ok := args[key]
	if !ok {
		return defaultVal
	}

	b, ok := val.(bool)
	if !ok {
		return defaultVal
	}

	return b
}

func getNumberWithDefault(args map[string]interface{}, key string, defaultVal float64) float64 {
	val, ok := args[key]
	if !ok {
		return defaultVal
	}

	switch num := val.(type) {
	case float64:
		return num
	case int:
		return float64(num)
	case int64:
		return float64(num)
	default:
		return defaultVal
	}
}
