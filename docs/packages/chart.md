# chart

ASCII chart rendering for terminal output.

## Description

The `chart` package provides ASCII-based visualizations including bar charts, line charts, and pie charts that can be displayed in terminal environments.

**Key Features:**
- Horizontal and vertical bar charts
- Line charts with trend visualization
- Pie charts (rendered as horizontal bars with percentages)
- Number formatting with thousand separators
- Configurable chart dimensions and labels

## Installation

```go
import "github.com/valksor/go-toolkit/chart"
```

## Usage

### Bar Chart

```go
bars := []chart.Bar{
    {Label: "January", Value: 150},
    {Label: "February", Value: 230},
    {Label: "March", Value: 180},
}

opts := chart.Options{
    Title:      "Monthly Sales",
    Width:      50,
    ShowValues: true,
}

fmt.Println(chart.BarChart(bars, opts))
```

Output:
```
Monthly Sales
January         │████████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│ 150
February        │██████████████████████████████████████████████████│ 230
March           │████████████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░│ 180
```

### Line Chart

```go
points := []chart.DataPoint{
    {Label: "Mon", Value: 10},
    {Label: "Tue", Value: 25},
    {Label: "Wed", Value: 18},
    {Label: "Thu", Value: 32},
    {Label: "Fri", Value: 28},
}

opts := chart.Options{
    Title: "Weekly Trend",
}

fmt.Println(chart.LineChart(points, opts))
```

### Pie Chart

```go
slices := []chart.Slice{
    {Label: "Planning", Value: 30},
    {Label: "Development", Value: 50},
    {Label: "Testing", Value: 20},
}

opts := chart.Options{
    Title: "Time Distribution",
}

fmt.Println(chart.PieChart(slices, opts))
```

Output:
```
Time Distribution
Planning             │███████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│ 30.0%
Development          │█████████████████████████░░░░░░░░░░░░░░░░░░░░░░░░░│ 50.0%
Testing              │██████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│ 20.0%
```

## API Reference

### Types

- `Bar` - Represents a single bar in a chart
- `DataPoint` - Represents a point in a line chart
- `Slice` - Represents a slice in a pie chart
- `Options` - Controls chart rendering options

### Functions

- `BarChart(bars []Bar, opts Options) string` - Generates an ASCII bar chart
- `LineChart(points []DataPoint, opts Options) string` - Generates an ASCII line chart
- `PieChart(slices []Slice, opts Options) string` - Generates an ASCII pie chart
- `FormatNumber(n int) string` - Formats a number with thousand separators

### Options Fields

| Field | Type | Description |
|-------|------|-------------|
| `Title` | `string` | Chart title (displayed above chart) |
| `Width` | `int` | Chart width in characters (default 60) |
| `Height` | `int` | Number of bars (0 = auto) |
| `ShowValues` | `bool` | Show numeric values |
| `Vertical` | `bool` | Render vertical bar chart |
| `ScaleLabel` | `string` | Label for the value axis |

## Common Patterns

### Vertical Bar Chart

```go
opts := chart.Options{
    Title:    "Quarterly Results",
    Vertical: true,
}
```

### Custom Max Value

```go
bars := []chart.Bar{
    {Label: "A", Value: 50, MaxValue: 100},
    {Label: "B", Value: 75, MaxValue: 100},
}
```

## See Also

- [display](packages/display.md) - Terminal colors and spinners
- [cli](packages/cli.md) - CLI helpers
