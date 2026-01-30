package chart

import (
	"strings"
	"testing"
)

func TestBarChart(t *testing.T) {
	tests := []struct {
		name         string
		bars         []Bar
		opts         Options
		wantContains []string
	}{
		{
			name: "simple bar chart",
			bars: []Bar{
				{Label: "Task A", Value: 50, MaxValue: 100},
				{Label: "Task B", Value: 75, MaxValue: 100},
			},
			opts: Options{
				Title: "Test Chart",
				Width: 40,
			},
			wantContains: []string{"Test Chart", "Task A", "Task B"},
		},
		{
			name: "with values",
			bars: []Bar{
				{Label: "Item", Value: 1000},
			},
			opts: Options{
				Width:      20,
				ShowValues: true,
			},
			wantContains: []string{"Item", "1,000"},
		},
		{
			name: "long label truncation",
			bars: []Bar{
				{Label: "This is a very long label that should be truncated", Value: 50},
			},
			opts: Options{
				Width: 30,
			},
			wantContains: []string{"..."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BarChart(tt.bars, tt.opts)
			if got == "" {
				t.Error("BarChart returned empty string")
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("BarChart() should contain %q", want)
				}
			}
		})
	}
}

func TestLineChart(t *testing.T) {
	tests := []struct {
		name           string
		dataPoints     []DataPoint
		opts           Options
		shouldHaveData bool
	}{
		{
			name: "simple line chart",
			dataPoints: []DataPoint{
				{"Mon", 10},
				{"Tue", 20},
				{"Wed", 15},
			},
			opts: Options{
				Title: "Weekly Trend",
			},
			shouldHaveData: true,
		},
		{
			name:           "empty data",
			dataPoints:     []DataPoint{},
			opts:           Options{},
			shouldHaveData: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LineChart(tt.dataPoints, tt.opts)
			if got == "" {
				t.Error("LineChart returned empty string")
			}
			if tt.shouldHaveData {
				// Just check that the chart has some content
				if len(got) < 10 {
					t.Error("LineChart output seems too short")
				}
			} else {
				if !strings.Contains(got, "(no data)") {
					t.Error("LineChart() should contain '(no data)' for empty input")
				}
			}
		})
	}
}

func TestPieChart(t *testing.T) {
	tests := []struct {
		name         string
		slices       []Slice
		opts         Options
		wantContains []string
	}{
		{
			name: "simple pie chart",
			slices: []Slice{
				{"Planning", 30, 30.0},
				{"Implementing", 50, 50.0},
				{"Review", 20, 20.0},
			},
			opts: Options{
				Title: "Time Distribution",
			},
			wantContains: []string{"Time Distribution", "Planning", "Implementing", "Review"},
		},
		{
			name: "with percentages",
			slices: []Slice{
				{"Step A", 25, 25.0},
				{"Step B", 75, 75.0},
			},
			opts:         Options{},
			wantContains: []string{"25.0%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PieChart(tt.slices, tt.opts)
			if got == "" {
				t.Error("PieChart returned empty string")
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("PieChart() should contain %q", want)
				}
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0"},
		{1, "1"},
		{10, "10"},
		{100, "100"},
		{1000, "1,000"},
		{10000, "10,000"},
		{100000, "100,000"},
		{1000000, "1,000,000"},
		{1234567, "1,234,567"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := FormatNumber(tt.input); got != tt.want {
				t.Errorf("FormatNumber(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{0, 0},
		{1, 1},
		{-1, 1},
		{100, 100},
		{-100, 100},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := abs(tt.input); got != tt.want {
				t.Errorf("abs(%d) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestVerticalBarChart(t *testing.T) {
	bars := []Bar{
		{Label: "A", Value: 10, MaxValue: 20},
		{Label: "B", Value: 15, MaxValue: 20},
	}
	opts := Options{
		Title:    "Vertical",
		Vertical: true,
	}

	got := BarChart(bars, opts)
	if !strings.Contains(got, "Vertical") {
		t.Error("Vertical chart should contain title")
	}
	if !strings.Contains(got, "A") {
		t.Error("Vertical chart should contain label A")
	}
	if !strings.Contains(got, "B") {
		t.Error("Vertical chart should contain label B")
	}
}
