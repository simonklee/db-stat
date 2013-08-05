package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func parseSinceFlag(v string) time.Time {
	d, err := time.Parse("2006-01-02", v)

	if err != nil {
		d = time.Unix(0, 0).UTC()
	}

	return d
}

func parseToFlag(v string) time.Time {
	d, err := time.Parse("2006-01-02", v)

	if err != nil {
		d = time.Now().UTC()
	}

	return d
}

func parseOutputFlag(v string) []outputType {
	words := parseWords(v)
	out := make([]outputType, 0, len(words))

	for _, w := range words {
		w = strings.ToUpper(w)

		switch w {
		case "TERM":
			out = append(out, termOutput)
		case "PNG":
			out = append(out, imageOutput)
		default:
			flag.Usage()
			os.Exit(1)
		}
	}

	return out
}

func parseGroupByFlag(v string) string {
	out := strings.ToUpper(v)

	if !stringInSlice(out, []string{"DAY", "WEEK", "MONTH", "YEAR"}) {
		flag.Usage()
		os.Exit(1)
	}
	return out
}

func parseWords(raw string) []string {
	var words []string

	for _, word := range strings.Split(raw, ",") {
		word = strings.TrimSpace(word)

		if word != "" {
			words = append(words, word)
		}
	}
	return words
}

func main() {
	var (
		flagHelp        = flag.Bool("h", false, "this help")
		flagGrowth      = flag.Bool("growth", false, "display table growth")
		flagDns         = flag.String("dns", "kogama:kogama@tcp(localhost:3306)/kogama", "Data Source Name")
		flagDatabase    = flag.String("database", "kogama", "database name")
		flagTables      = flag.String("tables", "", "comma separated list of tables")
		flagOutput      = flag.String("output", "term", "comma separated list of output types. Options; png, term")
		flagDateColumns = flag.String("dateColumns", "", "comma separated list of dateColumns")
		flagSinceDate   = flag.String("since", "", "limit queries from this date")
		flagToDate      = flag.String("to", "", "limit queries to this date")
		flagGroupBy     = flag.String("groupBy", "DAY", "DAY|WEEK|MONTH|YEAR")
		flagVersion     = flag.Bool("v", false, "show version and exit")
		flagCpuprofile  = flag.String("cpuprofile", "", "write cpu profile to file")
	)
	flag.Usage = usage
	flag.Parse()

	if *flagVersion {
		fmt.Fprintln(os.Stderr, "0.0.1")
		return
	}

	if *flagHelp {
		flag.Usage()
		os.Exit(1)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if *flagCpuprofile != "" {
		f, err := os.Create(*flagCpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	dbConnect(*flagDns)

	if *flagGrowth {
		tables := parseWords(*flagTables)
		dateColumns := parseWords(*flagDateColumns)
		groupBy := parseGroupByFlag(*flagGroupBy)
		since := parseSinceFlag(*flagSinceDate)
		to := parseToFlag(*flagToDate)
		charts := tableGrowthStat(*flagDatabase, tables, dateColumns, groupBy, since, to)
		outputTypes := parseOutputFlag(*flagOutput)

		for _, c := range charts {
			for _, t := range outputTypes {
				c.Write(t)
			}
		}
	} else {
		tableStat(*flagDatabase, *flagTables)
	}

	defer db.Close()
}
