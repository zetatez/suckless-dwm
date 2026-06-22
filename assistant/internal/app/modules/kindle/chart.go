package kindle

import (
	"bytes"
	"fmt"
	"strings"
)

type chartPoint struct {
	Label string
	Value float64
}

type chartSeries struct {
	Name   string
	Color  string
	Points []chartPoint
}

func svgMultiLineChart(series []chartSeries, width, height int, title string) string {
	if len(series) == 0 {
		return ""
	}

	const padL, padR, padT, padB = 55, 20, 8, 35
	cw := width - padL - padR
	ch := height - padT - padB

	type normSeries struct {
		Name   string
		Color  string
		Points []float64
		Last   float64
	}
	var norms []normSeries
	var minY, maxY float64
	first := true

	for _, s := range series {
		if len(s.Points) < 2 {
			continue
		}
		base := s.Points[0].Value
		if base == 0 {
			continue
		}
		last := s.Points[len(s.Points)-1].Value
		var pts []float64
		for _, p := range s.Points {
			v := (p.Value - base) / base * 100
			pts = append(pts, v)
			if first || v < minY {
				minY = v
			}
			if first || v > maxY {
				maxY = v
			}
			first = false
		}
		norms = append(norms, normSeries{Name: s.Name, Color: s.Color, Points: pts, Last: last})
	}
	if len(norms) == 0 {
		return ""
	}

	ran := maxY - minY
	if ran == 0 {
		ran = 1
	}
	pad := ran * 0.1
	minY -= pad
	maxY += pad
	ran = maxY - minY

	var b bytes.Buffer
	b.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" style="width:100%%;height:auto;display:block">`, width, height))

	// grid + y labels (%)
	steps := 4
	for i := 0; i <= steps; i++ {
		y := padT + ch*i/steps
		v := maxY - ran*float64(i)/float64(steps)
		b.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#eee" stroke-width="1"/>`, padL, y, padL+cw, y))
		label := fmt.Sprintf("%.1f%%", v)
		if v > 0 {
			label = "+" + label
		}
		b.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="9" fill="#888" text-anchor="end">%s</text>`, padL-4, y+3, label))
	}

	// zero line
	zeroY := padT + ch - int(float64(ch)*(0-minY)/ran)
	b.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#ccc" stroke-width="1" stroke-dasharray="4,2"/>`, padL, zeroY, padL+cw, zeroY))

	// lines
	dashPatterns := []string{"", "4,3", "1,3"} // solid, dashed, dotted
	for idx, n := range norms {
		var d []string
		for i, v := range n.Points {
			x := padL + cw*i/(len(n.Points)-1)
			y := padT + ch - int(float64(ch)*(v-minY)/ran)
			if i == 0 {
				d = append(d, fmt.Sprintf("M%d %d", x, y))
			} else {
				d = append(d, fmt.Sprintf("L%d %d", x, y))
			}
		}
		dash := ""
		if idx < len(dashPatterns) && dashPatterns[idx] != "" {
			dash = fmt.Sprintf(` stroke-dasharray="%s"`, dashPatterns[idx])
		}
		b.WriteString(fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="2"%s/>`, strings.Join(d, " "), n.Color, dash))
	}

	// x labels (first, middle, last)
	ref := series[0].Points
	for _, idx := range []int{0, len(ref) / 2, len(ref) - 1} {
		x := padL + cw*idx/(len(ref)-1)
		b.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="9" fill="#888" text-anchor="middle">%s</text>`, x, padT+ch+16, ref[idx].Label))
	}

	// legend with latest values
	legendY := height - 4
	for i, n := range norms {
		lx := padL + i*95
		dash := ""
		if i < len(dashPatterns) && dashPatterns[i] != "" {
			dash = fmt.Sprintf(` stroke-dasharray="%s"`, dashPatterns[i])
		}
		b.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2"%s/>`, lx, legendY-4, lx+12, legendY-4, n.Color, dash))
		lastStr := fmt.Sprintf("%.0f", n.Last)
		if n.Last >= 1000 {
			lastStr = fmt.Sprintf("%.1fK", n.Last/1000)
		}
		b.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="9" fill="#666">%s %s</text>`, lx+16, legendY, n.Name, lastStr))
	}

	b.WriteString(`</svg>`)
	return b.String()
}
