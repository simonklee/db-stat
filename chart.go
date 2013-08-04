package main

import (
	"bytes"
	//"github.com/simonz05/go-gnuplot"
	"fmt"
	//"github.com/vdobler/chart"
	//"github.com/vdobler/chart/txtg"
	"os"
)

// PlotX will create a 2-d plot using `data` as input and `title` as the plot
// title.
// The index of the element in the `data` slice will be used as the x-coordinate
// and its correspinding value as the y-coordinate.
// Example:
//  err = p.PlotX([]float64{10, 20, 30}, "my title")
func PlotX(data []float64, title string) error {
	var buf bytes.Buffer
	for _, d := range data {
		buf.WriteString(fmt.Sprintf("%v\n", d))
	}

	fmt.Fprintf(os.Stdout, "%s", buf.String())
	//cmd = "replot"

	//var line string
	//if title == "" {
	//	line = fmt.Sprintf("%s \"%s\" with %s", cmd, fname, self.style)
	//} else {
	//	line = fmt.Sprintf("%s \"%s\" title \"%s\" with %s",
	//		cmd, fname, title, self.style)
	//}
	return nil
}
func plot() {
	fmt.Fprintf(os.Stdout, "set terminal 'dumb';")
	fmt.Fprintf(os.Stdout, "plot \"-\" using 1:2 notitle with lines;")
    PlotX([]float64{0,1,2,3,4,5,6,7,8,9,10}, "some data")

    //p.CheckedCmd("set terminal pdf")
    //p.CheckedCmd("set output 'plot002.pdf'")
    //p.CheckedCmd("replot")
    //p.CheckedCmd("q")
    return
}
