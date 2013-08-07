package main

import (
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
	Write(*Chart)
}

func point2Chart(in []*Point) []chart.XYErrValue {
	out := make([]chart.XYErrValue, 0, len(in))

	for _, v := range in {
		out = append(out, v)
	}

	return out
}

type Chart struct {
	c    chart.Chart
	name string
}

type ImageWriter struct {}

func (im *ImageWriter) Write(c *Chart) {
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

type TermWriter struct {}

func (tm *TermWriter) Write(c *Chart) {
	tgr := txtg.New(100, 40)
	c.c.Plot(tgr)
	os.Stdout.Write([]byte(tgr.String() + "\n\n\n"))
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

	c.FmtVal = func(value, sum float64) (s string) {
		return humanize.Bytes(uint64(value))
	}
	c.Inner = 0.3

	return &Chart{
		c:    c,
		name: safeFilename(title + " pie chart"),
	}
}
