
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	help       = flag.Bool("h", false, "this help")
	growth     = flag.Bool("growth", false, "display table growth")
	dns		   = flag.String("dns", "kogama:kogama@tcp(localhost:3306)/kogama", "Data Source Name")
	database   = flag.String("database", "kogama", "database name")
	tables	   = flag.String("tables", "", "comma separated list of tables")
	output	   = flag.String("output", "term", "specify output type. Available options svg, term")
	datetimeColumns	= flag.String("datetimeColumns", "", "comma separated list of datetimeColumns")
	version    = flag.Bool("v", false, "show version and exit")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Fprintln(os.Stderr, "0.0.1")
		return
	}

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	dbConnect(*dns)
	plot()
	return

	if *growth {
		tableGrowthStat(*database, *tables, *datetimeColumns)
	} else {
		tableStat(*database, *tables)
	}

	defer db.Close()
}
