package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	var (
		d         data
		dbStruct  dbStruct
		dbGenFile dbGenFile
		f         *os.File
		err       error
	)

	flag.StringVar(&d.Package, "package", "", "Package name")
	flag.StringVar(&dbGenFile.GenDir, "dir", "", "Directory")
	flag.StringVar(&dbGenFile.TypeFile, "typefile", "", "Type file")
	flag.StringVar(&dbGenFile.MockFile, "mockfile", "", "Mock file")
	flag.StringVar(&dbGenFile.FuncFile, "funcfile", "", "Func file")
	flag.StringVar(&dbGenFile.TestFile, "testfile", "", "Test file")
	flag.StringVar(&dbStruct.TableName, "table", "", "Table name")
	flag.StringVar(&dbStruct.Host, "host", "", "Host name")
	flag.StringVar(&dbStruct.Port, "port", "5432", "Port")
	flag.StringVar(&dbStruct.User, "user", "", "User")
	flag.StringVar(&dbStruct.Password, "password", "", "Password")
	flag.StringVar(&dbStruct.DBName, "dbname", "", "DB Name")
	flag.StringVar(&dbStruct.Schema, "schema", "public", "Schema")
	flag.Parse()

	if dbGenFile.TypeFile == "" {
		log.Fatalln("-typefile is empty")
	}

	var golangTableName string
	splits := strings.Split(dbStruct.TableName, "_")
	for _, v := range splits {
		golangTableName += strings.Title(strings.ToLower(v))
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbStruct.User,
		dbStruct.Password,
		dbStruct.DBName,
		dbStruct.Host,
		dbStruct.Port,
	))
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	var dbStructures dbStructures
	if err := dbStructures.Select(db, dbStruct); err != nil {
		log.Fatalln(err)
	}

	arrDBStructure := []dbStructure(dbStructures)
	if len(arrDBStructure) == 0 {
		log.Fatalln("Error empty table")
	}

	mapDataType := parseDatatype()
	for i := 0; i < len(arrDBStructure); i++ {
		arrDBStructure[i].GolangDataType = "interface{}"
		if mapDataType[arrDBStructure[i].DataType] != "" {
			arrDBStructure[i].GolangDataType = mapDataType[arrDBStructure[i].DataType]
		}

		splits := strings.Split(arrDBStructure[i].Column, "_")
		for _, v := range splits {
			arrDBStructure[i].GolangColumn += strings.Title(strings.ToLower(v))
		}
	}

	d.GolangTableName = golangTableName
	d.DBStructures = arrDBStructure

	typePath := dbGenFile.TypeFile
	if dbGenFile.GenDir != "" {
		typePath = fmt.Sprintf("%s/%s", dbGenFile.GenDir, dbGenFile.TypeFile)
	}

	templateString := queueTemplate
	if fileExists(typePath) {
		if f, err = os.OpenFile(typePath, os.O_RDWR|os.O_APPEND, 0660); err != nil {
			log.Fatalln("Error open file : ", err)
		}
	} else {
		templateString = newTemplate + queueTemplate

		dirName := filepath.Dir(typePath)
		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			log.Fatalln("Error create directory : ", err)
		}

		if f, err = os.Create(typePath); err != nil {
			log.Fatalln("Error create file : ", err)
		}
	}
	defer f.Close()

	buffData := new(bytes.Buffer)

	t := template.Must(template.New("genfile").Parse(templateString))
	w := tabwriter.NewWriter(buffData, 8, 8, 8, ' ', 0)
	t.Execute(w, d)
	w.Flush()

	if _, err := f.WriteString(buffData.String()); err != nil {
		log.Fatalln("Error write file : ", err)
	}
}
