package main

type data struct {
	Package         string
	GolangTableName string
	DBStructures    []dbStructure
}

type dbStruct struct {
	Host      string
	User      string
	Password  string
	DBName    string
	Schema    string
	Port      string
	TableName string
}

type dbGenFile struct {
	GenDir   string
	TypeFile string
	MockFile string
	FuncFile string
	TestFile string
}

type dbStructure struct {
	Column         string `db:"column_name"`
	DataType       string `db:"data_type"`
	GolangDataType string
	GolangColumn   string
}

type dbStructures []dbStructure
