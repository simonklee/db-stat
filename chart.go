package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/vdobler/chart"
	"github.com/vdobler/chart/imgg"
	"github.com/vdobler/chart/txtg"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path"
)

type ChartWriter interface {
	Write(chart.Chart)
}

func point2Chart(in []*Point) []chart.XYErrValue {
	out := make([]chart.XYErrValue, 0, len(in))

	for _, v := range in {
		out = append(out, v)
	}

	return out
}

type outputType int

const (
	termOutput outputType = iota
	imageOutput
)

type Chart struct {
	c    chart.Chart
	name string
}

func (c *Chart) ImageWrite() {
	os.MkdirAll("data", os.ModePerm)

	fp, err := os.Create(path.Join("data", c.name+".png"))
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	img := image.NewRGBA(image.Rect(0, 0, 1024, 768))
	bg := image.NewUniform(color.RGBA{0xff, 0xff, 0xff, 0xff})
	draw.Draw(img, img.Bounds(), bg, image.ZP, draw.Src)

	//row, col := d.Cnt/d.N, d.Cnt%d.N
	igr := imgg.AddTo(img, 0, 0, 1024, 768, color.RGBA{0xff, 0xff, 0xff, 0xff}, nil, nil)
	c.c.Plot(igr)

	png.Encode(fp, img)
}

func (c *Chart) TermWrite() {
	tgr := txtg.New(100, 40)
	c.c.Plot(tgr)
	os.Stdout.Write([]byte(tgr.String() + "\n\n\n"))
}

func (c *Chart) Write(t outputType) {
	switch t {
	case termOutput:
		c.TermWrite()
	case imageOutput:
		c.ImageWrite()
	default:
	}
}

func TimeChart(title, xlabel, ylabel string, data []*Point) *Chart {
	c := &chart.ScatterChart{Title: title}
	c.XRange.Label = xlabel
	c.YRange.Label = ylabel
	c.XRange.Time = true
	c.XRange.TicSetting.Mirror = 1

	style := chart.AutoStyle(4, true)
	c.AddDataGeneric(ylabel, point2Chart(data), chart.PlotStyleLinesPoints, style)

	return &Chart{
		c:    c,
		name: safeFilename(title + " time chart"),
	}
}

func PieChart(title string, labels []string, data []float64) *Chart {
	c := &chart.PieChart{Title: title}
	c.AddDataPair("Tables", labels, data)

	c.FmtVal = func (value, sum float64) (s string) {
		return humanize.Bytes(uint64(value))
	}
	c.Inner = 0.3

	return &Chart{
		c:    c,
		name: safeFilename(title + " pie chart"),
	}
}

func BarChart(title string, labels []string, data []float64) *Chart {
	c := &chart.BarChart{Title: title}
	c.XRange.Category = labels
	blue := chart.Style{Symbol: '+', LineColor: color.NRGBA{0x00, 0x00, 0xff, 0xff}, LineWidth: 3 }
	//green := chart.Style{Symbol: 'x', LineColor: color.NRGBA{0x00, 0xaa, 0x00, 0xff}, LineWidth: 3, FillColor: color.NRGBA{0x40, 0xff, 0x40, 0xff}}

	x := make([]float64, 0, len(data))

	for i := 0; i < len(data); i++ {
		x = append(x, float64(i))
	}

	c.YRange.TicSetting.Format = func(f float64) string {
		return humanize.Bytes(uint64(f))
	}

	c.XRange.TicSetting.Format = func(f float64) string {
		return humanize.Bytes(uint64(f))
	}

	c.Stacked = true
	for i := range data {
		c.AddData("", []chart.Point{
			chart.Point{X:x[i], Y: data[i]},
		}, blue)
	}
	c.ShowVal = 2

	fmt.Println(title, labels, data)
	// Bar Chart
	//x := []float64{0, 1, 2}
	//europe := []float64{1, 0}
	//africa := []float64{20, 5, 5, 5}
	//blue := chart.Style{Symbol: '#', LineColor: color.NRGBA{0x00, 0x00, 0xff, 0xff}, LineWidth: 3, FillColor: color.NRGBA{0x40, 0x40, 0xff, 0xff}}
	//green := chart.Style{Symbol: 'x', LineColor: color.NRGBA{0x00, 0xaa, 0x00, 0xff}, LineWidth: 3, FillColor: color.NRGBA{0x40, 0xff, 0x40, 0xff}}

	//bar := chart.BarChart{Title: "Income Distribution"}
	//bar.XRange.Category = []string{"low", "average", "high"}
	//bar.XRange.Label = "Income category"
	//bar.YRange.Label = "Adult population"
	//bar.YRange.TicSetting.Format = func(f float64) string {
	//	return fmt.Sprintf("%d%%", int(f+0.5))
	//}
	//bar.Stacked = true
	//bar.Key.Pos, bar.Key.Cols = "obc", 1
	//bar.AddDataPair("Europe", x, europe, blue)
	//bar.AddDataPair("Africa", x, africa, green)
	//bar.ShowVal = 1

	return &Chart{
		c:    c,
		name: safeFilename(title + " bar chart"),
	}
}
