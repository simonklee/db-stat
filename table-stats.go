package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db                  *sql.DB
	stmtTableSize       *sql.Stmt
	stmtTablesAvailable *sql.Stmt
)

func dbConnect(dns string) {
	var err error
	db, err = sql.Open("mysql", dns)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	prepareStatements()
}

func prepareStatements() {
	var err error
	//stmtTableSize, err = db.Prepare(`SELECT
	//	CONCAT(FORMAT(DAT/POWER(1024,pw1),2),' ',SUBSTR(units,pw1*2+1,2)) DATSIZE,
	//	CONCAT(FORMAT(NDX/POWER(1024,pw2),2),' ',SUBSTR(units,pw2*2+1,2)) NDXSIZE,
	//	CONCAT(FORMAT(TBL/POWER(1024,pw3),2),' ',SUBSTR(units,pw3*2+1,2)) TBLSIZE
	//FROM
	//(
	//	SELECT DAT,NDX,TBL,IF(px>4,4,px) pw1,IF(py>4,4,py) pw2,IF(pz>4,4,pz) pw3
	//	FROM
	//	(
	//		SELECT data_length DAT,index_length NDX,data_length+index_length TBL,
	//		FLOOR(LOG(IF(data_length=0,1,data_length))/LOG(1024)) px,
	//		FLOOR(LOG(IF(index_length=0,1,index_length))/LOG(1024)) py,
	//		FLOOR(LOG(IF(data_length+index_length=0,1,data_length+index_length))/LOG(1024)) pz
	//		FROM information_schema.tables
	//		WHERE table_schema=?
	//		AND table_name=?
	//	) AA
	//) A,(SELECT 'B KBMBGBTB' units) B`)

	stmtTableSize, err = db.Prepare(`SELECT data_length, index_length
		FROM information_schema.tables
		WHERE table_schema=? and table_name=?
	`)

	if err != nil {
		log.Fatal(err)
	}

	stmtTablesAvailable, err = db.Prepare(`SELECT TABLE_NAME
	FROM information_schema.TABLES
	WHERE TABLE_SCHEMA = ?`)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func tableGrowth(table, dateColumn, groupBy string, since, to time.Time) []*Point {
	if groupBy == "DAY" {
		groupBy = "DATE"
	}

	stmt := fmt.Sprintf(`SELECT DATE(%s), COUNT(*) AS Count 
	FROM %s 
	WHERE %s >= ? AND %s <= ?
	GROUP BY %s(%s)`, dateColumn, table, dateColumn, dateColumn, groupBy, dateColumn)
	rows, err := db.Query(stmt, since, to)

	if err != nil {
		log.Fatal(err)
	}

	var (
		date  string
		count int
	)

	var data []*Point

	for rows.Next() {
		err := rows.Scan(&date, &count)
		if err != nil {
			log.Fatal(err)
		}

		d, err := time.Parse("2006-01-02", date)

		if err != nil {
			log.Fatal(err)
		}

		data = append(data, &Point{float64(d.Unix()), float64(count)})
	}
	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	return data
}

func tablesAvailable(database string) (tables []string) {
	var tableName string

	rows, err := stmtTablesAvailable.Query(database)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}

		tables = append(tables, tableName)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return tables
}

func tableSize(database, table string) *TableSize {
	var (
		dataSize  float64
		indexSize float64
	)

	err := stmtTableSize.QueryRow(database, table).Scan(&dataSize, &indexSize)

	if err != nil {
		log.Fatal(err)
	}

	return &TableSize{
		Name:  table,
		Data:  dataSize,
		Index: indexSize,
		Total: dataSize + indexSize,
	}
}

func tableStat(database string, tables, ignoreTables []string, cutoff int) []*Chart {
	if len(tables) == 0 {
		tables = tablesAvailable(database)

		if len(ignoreTables) > 0 {
			// remove the ignored tables from the slice of tables
			last := len(tables)-1

			for i := 0; i < last + 1; i++ {
				table := tables[i]

				if stringInSlice(table, ignoreTables) {
					if i == last {
						last -= 1
					} else {
						// overwrite the current index
						tables[i] = tables[last]
						last -= 1
						i -= 1
					}
				}
			}
			tables = tables[0:last+1]
		}
	}

	var data []*TableSize

	for _, table := range tables {
		data = append(data, tableSize(database, table))
	}

	sort.Sort(sort.Reverse(ByTotal{data}))

	// Merge the smallest values into an other branch if we have many results
	if cutoff > 0 && len(data) > cutoff {
		other := data[cutoff]
		other.Name = "Other"

		for i := cutoff + 1; i < len(data); i++ {
			other.Index += data[i].Index
			other.Data += data[i].Data
			other.Total += data[i].Total
		}

		data = data[0 : cutoff+1]
	}

	vals := make([]float64, 0, len(data))
	labels := make([]string, 0, len(data))

	for i := range data {
		vals = append(vals, data[i].Total)
		labels = append(labels, data[i].Name)
	}

	return []*Chart{PieChart("Table Size Overview", labels, vals)}
}

func tableGrowthStat(database string, tables []string, dateColumns []string, groupBy string, since, to time.Time) []*Chart {
	if len(tables) != len(dateColumns) {
		log.Fatal(fmt.Sprintf("tables count %d != datetime columns count %d", len(tables), len(dateColumns)))
	}

	charts := make([]*Chart, 0, len(tables))

	for i := range tables {
		table := tables[i]
		dateColumn := dateColumns[i]
		data := tableGrowth(table, dateColumn, groupBy, since, to)

		var total float64
		for _, p := range data {
			total += p.Y
		}

		ylabel := fmt.Sprintf("%s: Created Per %s", table, strings.Title(strings.ToLower(groupBy)))
		xlabel := fmt.Sprintf("total in period: %s", humanize.Comma(int64(total)))

		charts = append(charts, TimeChart(ylabel, xlabel, ylabel, data))
	}

	return charts
}
