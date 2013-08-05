package main

import (
	"database/sql"
	"strconv"
	"fmt"
	"github.com/dustin/go-humanize"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
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

	stmtTableSize, err = db.Prepare(`SELECT (data_length+index_length) tablesize
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

func tableSize(database, table string) float64 {
	//var (
	//	dataSize  string
	//	indexSize string
	//)
	var size string

	err := stmtTableSize.QueryRow(database, table).Scan(&size)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("Database: %s, Table: %s\n", database, table)
	//fmt.Printf("DataSize: %s, IndexSize: %s, TableSize: %s\n", dataSize, indexSize, tableSize)
	//fmt.Printf("TableSize: %s\n", humanize.Bytes(tableSize))
	f64, _ := strconv.ParseFloat(size, 0)
	return f64
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

func tableStat(database string, tables []string) []*Chart {
	if len(tables) == 0 {
		tables = tablesAvailable(database)
	}

	var data []float64

	for _, table := range tables {
		data = append(data, tableSize(database, table))
	}

	return []*Chart{BarChart("Table Size Stats", tables, data)}
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
