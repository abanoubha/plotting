package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var vscodeBlueSeriesData = `
[2026-09-29 01:39:57] 0 MB
[2026-09-29 01:40:07] 2254.81 MB
[2026-09-29 01:40:17] 3739.42 MB
[2026-09-29 01:40:27] 3853.55 MB
[2026-09-29 01:40:38] 3653.62 MB
[2026-09-29 01:40:48] 4150.38 MB
[2026-09-29 01:40:58] 4140.62 MB
[2026-09-29 01:41:08] 3955.86 MB
[2026-09-29 01:41:18] 3938.73 MB
[2026-09-29 01:41:28] 3995.55 MB
[2026-09-29 01:41:38] 4079.58 MB
[2026-09-29 01:41:49] 4130.34 MB
[2026-09-29 01:41:59] 4148.77 MB
[2026-09-29 01:42:09] 4041.12 MB
[2026-09-29 01:42:19] 4005.03 MB
[2026-09-29 01:42:29] 914.078 MB
[2026-09-29 01:42:39] 822.188 MB
[2026-09-29 01:42:50] 822.156 MB
[2026-09-29 01:43:00] 818.609 MB
[2026-09-29 01:43:10] 0 MB
`

var zedbrownSeriesData = `
[2026-09-29 02:36:30] 0 MB
[2026-09-29 02:36:40] 474.109 MB
[2026-09-29 02:36:50] 880.203 MB
[2026-09-29 02:37:00] 880.516 MB
[2026-09-29 02:37:10] 890.062 MB
[2026-09-29 02:37:20] 889.766 MB
[2026-09-29 02:37:31] 890.047 MB
[2026-09-29 02:37:41] 896.312 MB
[2026-09-29 02:37:51] 898.375 MB
[2026-09-29 02:38:01] 558.094 MB
[2026-09-29 02:38:11] 560.422 MB
[2026-09-29 02:38:21] 532.328 MB
[2026-09-29 02:38:32] 522.75 MB
[2026-09-29 02:38:42] 537.656 MB
[2026-09-29 02:38:52] 537.719 MB
[2026-09-29 02:39:02] 793.391 MB
[2026-09-29 02:39:12] 866.875 MB
[2026-09-29 02:39:22] 911.312 MB
[2026-09-29 02:39:33] 911.359 MB
[2026-09-29 02:39:43] 868.156 MB
[2026-09-29 02:39:53] 868.25 MB
[2026-09-29 02:40:03] 868.328 MB
[2026-09-29 02:40:13] 868.359 MB
[2026-09-29 02:40:23] 908.156 MB
[2026-09-29 02:40:34] 1337.08 MB
[2026-09-29 02:40:44] 1351.16 MB
[2026-09-29 02:40:54] 1350.45 MB
[2026-09-29 02:41:04] 1350.48 MB
[2026-09-29 02:41:14] 1349.91 MB
[2026-09-29 02:41:24] 1349.97 MB
[2026-09-29 02:41:35] 886.438 MB
[2026-09-29 02:41:45] 886.438 MB
[2026-09-29 02:41:55] 0 MB
`

func parseValues(input string) []float64 {
	re := regexp.MustCompile(`\]\s*([\d\.]+)\s*MB`)
	lines := strings.Split(input, "\n")
	var values []float64
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			val, err := strconv.ParseFloat(matches[1], 64)
			if err == nil {
				values = append(values, val)
			}
		}
	}
	return values
}

func main() {
	blueSeries := parseValues(vscodeBlueSeriesData)
	brownSeries := parseValues(zedbrownSeriesData)

	maxVal := 0.0
	for _, v := range blueSeries {
		if v > maxVal {
			maxVal = v
		}
	}
	for _, v := range brownSeries {
		if v > maxVal {
			maxVal = v
		}
	}

	maxPoints := len(blueSeries)
	if len(brownSeries) > maxPoints {
		maxPoints = len(brownSeries)
	}

	// 1. Output ASCII Chart
	fmt.Print(renderASCII(blueSeries, brownSeries, maxVal, maxPoints))

	// 2. Generate and write SVG File
	svgContent := generateSVG(blueSeries, brownSeries, maxVal, maxPoints)
	err := os.WriteFile("memory_comparison.svg", []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("Error writing SVG file: %v\n", err)
	} else {
		fmt.Println("SVG chart generated and saved to memory_comparison.svg")
	}
}

func renderASCII(blue, brown []float64, maxVal float64, maxPoints int) string {
	var sb strings.Builder
	sb.WriteString("=== Memory Usage Comparison (ASCII) ===\n")
	sb.WriteString("Legend: B = Blue Series | A = Brown Series | X = Overlap\n\n")

	rows := 15
	cols := maxPoints

	grid := make([][]rune, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]rune, cols)
		for c := 0; c < cols; c++ {
			grid[r][c] = ' '
		}
	}

	for idx, val := range blue {
		norm := val / maxVal
		rIdx := (rows - 1) - int(math.Round(norm*float64(rows-1)))
		if rIdx >= 0 && rIdx < rows {
			grid[rIdx][idx] = 'B'
		}
	}

	for idx, val := range brown {
		norm := val / maxVal
		rIdx := (rows - 1) - int(math.Round(norm*float64(rows-1)))
		if rIdx >= 0 && rIdx < rows {
			if grid[rIdx][idx] == 'B' {
				grid[rIdx][idx] = 'X'
			} else {
				grid[rIdx][idx] = 'A'
			}
		}
	}

	for r := 0; r < rows; r++ {
		yVal := maxVal * float64(rows-1-r) / float64(rows-1)
		sb.WriteString(fmt.Sprintf("%7.1f MB | ", yVal))
		for c := 0; c < cols; c++ {
			sb.WriteRune(grid[r][c])
		}
		sb.WriteString("\n")
	}

	sb.WriteString("           +")
	for c := 0; c < cols; c++ {
		sb.WriteString("-")
	}
	sb.WriteString("\n           Index: 0")
	for c := 10; c < cols; c += 10 {
		sb.WriteString(fmt.Sprintf("%*d", 10, c))
	}
	sb.WriteString("\n\n")

	return sb.String()
}

func generateSVG(blue, brown []float64, maxVal float64, maxPoints int) string {
	width := 900
	height := 450
	padLeft := 80
	padRight := 40
	padTop := 50
	padBottom := 50

	plotW := float64(width - padLeft - padRight)
	plotH := float64(height - padTop - padBottom)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d" style="background-color: #1e1e1e; font-family: monospace;">`+"\n", width, height, width, height))

	// Chart Title
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="30" fill="#ffffff" font-size="16" font-weight="bold" text-anchor="middle">Memory Consumption Comparison</text>`+"\n", width/2))

	// Gridlines and Y-axis Labels
	ticks := 5
	for i := 0; i <= ticks; i++ {
		frac := float64(i) / float64(ticks)
		yPos := float64(height-padBottom) - frac*plotH
		val := frac * maxVal

		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%.2f" x2="%d" y2="%.2f" stroke="#333333" stroke-width="1" stroke-dasharray="4" />`+"\n", padLeft, yPos, width-padRight, yPos))
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%.2f" fill="#aaaaaa" font-size="12" text-anchor="end" dominant-baseline="middle">%.1f MB</text>`+"\n", padLeft-10, yPos, val))
	}

	// Gridlines and X-axis Labels
	xStep := 10
	for i := 0; i < maxPoints; i += xStep {
		xFrac := float64(i) / float64(maxPoints-1)
		xPos := float64(padLeft) + xFrac*plotW

		sb.WriteString(fmt.Sprintf(`<line x1="%.2f" y1="%d" x2="%.2f" y2="%d" stroke="#333333" stroke-width="1" stroke-dasharray="4" />`+"\n", xPos, padTop, xPos, height-padBottom))
		sb.WriteString(fmt.Sprintf(`<text x="%.2f" y="%d" fill="#aaaaaa" font-size="12" text-anchor="middle">Idx %d</text>`+"\n", xPos, height-padBottom+20, i))
	}

	// Axes Border Lines
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#aaaaaa" stroke-width="1.5" />`+"\n", padLeft, padTop, padLeft, height-padBottom))
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#aaaaaa" stroke-width="1.5" />`+"\n", padLeft, height-padBottom, width-padRight, height-padBottom))

	makePoints := func(series []float64) string {
		var pts []string
		for idx, val := range series {
			xFrac := float64(idx) / float64(maxPoints-1)
			xPos := float64(padLeft) + xFrac*plotW
			yPos := float64(height-padBottom) - (val/maxVal)*plotH
			pts = append(pts, fmt.Sprintf("%.2f,%.2f", xPos, yPos))
		}
		return strings.Join(pts, " ")
	}

	// Blue Series (#2196F3) Line
	sb.WriteString(fmt.Sprintf(`<polyline fill="none" stroke="#2196F3" stroke-width="2.5" points="%s" />`+"\n", makePoints(blue)))

	// Brown Series (#8D6E63) Line
	sb.WriteString(fmt.Sprintf(`<polyline fill="none" stroke="#8D6E63" stroke-width="2.5" points="%s" />`+"\n", makePoints(brown)))

	// Blue Points
	for idx, val := range blue {
		xFrac := float64(idx) / float64(maxPoints-1)
		xPos := float64(padLeft) + xFrac*plotW
		yPos := float64(height-padBottom) - (val/maxVal)*plotH
		sb.WriteString(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="3" fill="#2196F3" />`+"\n", xPos, yPos))
	}

	// Brown Points
	for idx, val := range brown {
		xFrac := float64(idx) / float64(maxPoints-1)
		xPos := float64(padLeft) + xFrac*plotW
		yPos := float64(height-padBottom) - (val/maxVal)*plotH
		sb.WriteString(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="3" fill="#8D6E63" />`+"\n", xPos, yPos))
	}

	// Legend Box
	legX := width - padRight - 150
	legY := padTop + 10
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="140" height="60" fill="#2d2d2d" stroke="#444444" rx="4" />`+"\n", legX, legY))
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#2196F3" stroke-width="3" />`+"\n", legX+10, legY+20, legX+30, legY+20))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="#ffffff" font-size="12" dominant-baseline="middle">VS Code</text>`+"\n", legX+40, legY+20))
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#8D6E63" stroke-width="3" />`+"\n", legX+10, legY+40, legX+30, legY+40))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="#ffffff" font-size="12" dominant-baseline="middle">Zed editor</text>`+"\n", legX+40, legY+40))

	sb.WriteString(`</svg>` + "\n")
	return sb.String()
}
