package main

import (
	"strings"
	"database/sql"
	"fmt"
	"log"
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
	stmtTableSize, err = db.Prepare(`SELECT 
		CONCAT(FORMAT(DAT/POWER(1024,pw1),2),' ',SUBSTR(units,pw1*2+1,2)) DATSIZE,
		CONCAT(FORMAT(NDX/POWER(1024,pw2),2),' ',SUBSTR(units,pw2*2+1,2)) NDXSIZE,
		CONCAT(FORMAT(TBL/POWER(1024,pw3),2),' ',SUBSTR(units,pw3*2+1,2)) TBLSIZE
	FROM
	(
		SELECT DAT,NDX,TBL,IF(px>4,4,px) pw1,IF(py>4,4,py) pw2,IF(pz>4,4,pz) pw3
		FROM 
		(
			SELECT data_length DAT,index_length NDX,data_length+index_length TBL,
			FLOOR(LOG(IF(data_length=0,1,data_length))/LOG(1024)) px,
			FLOOR(LOG(IF(index_length=0,1,index_length))/LOG(1024)) py,
			FLOOR(LOG(IF(data_length+index_length=0,1,data_length+index_length))/LOG(1024)) pz
			FROM information_schema.tables
			WHERE table_schema=?
			AND table_name=?
		) AA
	) A,(SELECT 'B KBMBGBTB' units) B`)

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

func tableGrowth(table, datetimeColumn string) {
	stmt := fmt.Sprintf(`SELECT DATE(%s), COUNT(*) AS Count 
	FROM %s 
	GROUP BY DATE(%s)`, datetimeColumn, table, datetimeColumn)
	rows, err := db.Query(stmt)

	if err != nil {
		log.Fatal(err)
	}

	var (
		date string
		count int
	)

	for rows.Next() {
		err := rows.Scan(&date, &count)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(date, count)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func tableSize(database, table string) {
	var (
		dataSize  string
		indexSize string
		tableSize string
	)

	err := stmtTableSize.QueryRow(database, table).Scan(&dataSize, &indexSize, &tableSize)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("Database: %s, Table: %s\n", database, table)
	fmt.Printf("DataSize: %s, IndexSize: %s, TableSize: %s\n", dataSize, indexSize, tableSize)
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

func parseWords(raw string) []string {
	var words []string

	for _, word := range strings.Split(raw, ",") {
		word = strings.TrimSpace(word)

		if word != "" {
			words = append(words , word)
		} 
	}
	return words
}

func tableStat(database, rawTables string) {
	var tables []string

	if rawTables == "" {
		tables = tablesAvailable(database)
	} else {
		tables = parseWords(rawTables)
	}

	for _, table := range tables {
		tableSize(database, table)
	}
}

func tableGrowthStat(database, rawTables, rawColumns string) {
	var tables []string
	var datetimeColumns []string

	tables = parseWords(rawTables)
	datetimeColumns = parseWords(rawColumns)

	if len(tables) != len(datetimeColumns) {
		log.Fatal(fmt.Sprintf("tables count %d != datetime columns count %d", len(tables), len(datetimeColumns)))
	}

	for i := range tables {
		tableGrowth(tables[i], datetimeColumns[i])
	}
}
