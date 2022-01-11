package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/danangkonang/oauth2-golang/db/migration"
	"github.com/danangkonang/oauth2-golang/db/seeder"
	"github.com/joho/godotenv"
)

var (
	help    bool
	version bool
)

type Tables struct {
	NameTable []string
}

var (
	MigrationFolder = "db/migration"
	SeederFolder    = "db/seeder"
	green           = "\033[32m"
	reset           = "\033[0m"
)

type ComandUsage struct {
	CmdName  string
	CmdAlias string
	CmdDesc  string
}

type FlagCmd struct {
	FlagName  string
	FlagAlias string
	FlagDesc  string
}

type Helper struct {
	Usage    string
	Version  string
	Error    string
	Option   []*ComandUsage
	Argument []*FlagCmd
}

var versionTmp = `Version: {{ .Version }}
`

var errorTmp = `unknow comand '{{ .Error }}'

see 'gomigator --help'
`

var helperTmp = `
Usage: {{ .Usage }}
{{ if .Option }}
Commands:
{{- range .Option}}
	{{ .CmdName }}      {{"	"}}{{ .CmdDesc }}{{ end -}}
{{ end }}

Options:
{{- range .Argument}}
	{{ .FlagName }}  {{"	"}}{{ .FlagDesc }}{{ end }}

`

func printTemplate(temp string, data interface{}) {
	tmpl, err := template.New("help").Parse(temp)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	os.Exit(0)
}

func globalHelp() {
	hlp := &Helper{
		Usage: "gomigator [COMAND] [OPTIONS]",
		Option: []*ComandUsage{
			{
				CmdName: "init",
				CmdDesc: "generate db directory for",
			},
			{
				CmdName: "create",
				CmdDesc: "create migration or seeder file",
			},
			{
				CmdName: "up",
				CmdDesc: "exect migration to database",
			},
			{
				CmdName: "down",
				CmdDesc: "drop migration on databse",
			},
			{
				CmdName: "migration",
				CmdDesc: "generate type migration",
			},
			{
				CmdName: "seeder",
				CmdDesc: "generate type seeder",
			},
		},
		Argument: []*FlagCmd{
			{
				FlagName: "--table",
				FlagDesc: "table name",
			},
			{
				FlagName: "--tables",
				FlagDesc: "list tables",
			},
			{
				FlagName: "--name",
				FlagDesc: "generate file name",
			},
		},
	}
	printTemplate(helperTmp, hlp)
}

func init() {
	godotenv.Load()
}

func main() {
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&help, "help", false, "help")
	flag.BoolVar(&version, "v", false, "version")
	flag.BoolVar(&version, "version", false, "version")
	flag.Parse()
	if help || len(os.Args[1:]) == 0 {
		globalHelp()
	}
	if version {
		hlp := &Helper{
			Version: "1.1.1",
		}
		printTemplate(versionTmp, hlp)
	}
	var t Tables
	up := flag.NewFlagSet("up", flag.ExitOnError)
	upMigration := flag.NewFlagSet("migration", flag.ContinueOnError)
	upMigration.Func("tables", "list file name", func(s string) error {
		t.NameTable = strings.Fields(s)
		return nil
	})
	upSeeder := flag.NewFlagSet("seeder", flag.ContinueOnError)
	upSeeder.Func("tables", "list file name", func(s string) error {
		t.NameTable = strings.Fields(s)
		return nil
	})

	down := flag.NewFlagSet("down", flag.ExitOnError)
	downMigration := flag.NewFlagSet("migration", flag.ContinueOnError)
	downMigration.Func("tables", "list file name", func(s string) error {
		t.NameTable = strings.Fields(s)
		return nil
	})
	downSeeder := flag.NewFlagSet("seeder", flag.ContinueOnError)
	downSeeder.Func("tables", "list file name", func(s string) error {
		t.NameTable = strings.Fields(s)
		return nil
	})
	if len(os.Args) < 2 {
		globalHelp()
	}
	switch os.Args[1] {
	case "up":
		up.Parse(os.Args[2:])
		upHandle(os.Args, upMigration, upSeeder, &t)
	case "down":
		down.Parse(os.Args[2:])
		downHandle(os.Args, upMigration, upSeeder, &t)
	default:
		hlp := &Helper{
			Error: os.Args[1],
		}
		printTemplate(errorTmp, hlp)
	}
}

func downHandle(argument []string, upMigration, upSeeder *flag.FlagSet, tb *Tables) {
	switch argument[2] {
	case "migration":
		upMigration.Parse(os.Args[3:])
		migrationDown(tb)
	case "seeder":
		upSeeder.Parse(os.Args[3:])
		seederDown(tb)
	default:
		hlp := &Helper{
			Error: os.Args[1],
		}
		printTemplate(errorTmp, hlp)
	}
}

func migrationDown(t *Tables) {
	files, err := ioutil.ReadDir(MigrationFolder)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	if len(t.NameTable) == 0 {
		newFile := []string{}
		for _, file := range files {
			filename := file.Name()
			if filename != "0.go" {
				newFile = append(newFile, filename)
			}
		}
		t.NameTable = newFile
	}
	m := migration.Migration{}
	for _, migrate := range t.NameTable {
		list := strings.Split(migrate, "_migration_")
		tb_name := strings.Split(list[1], ".go")
		meth := reflect.ValueOf(&m).MethodByName("Down" + strings.Title(tb_name[0]))
		meth.Call(nil)
	}
}

func seederDown(t *Tables) {
	files, err := ioutil.ReadDir(SeederFolder)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	if len(t.NameTable) == 0 {
		newFile := []string{}
		for _, file := range files {
			filename := file.Name()
			if filename != "0.go" {
				newFile = append(newFile, filename)
			}
		}
		t.NameTable = newFile
	}
	for _, migrate := range t.NameTable {
		list := strings.Split(migrate, "_seeder_")
		tb_name := strings.Split(list[1], ".go")
		var query string
		if os.Getenv("DB_DRIVER") == "mysql" {
			query = "TRUNCATE " + tb_name[0] + " ;"
		} else {
			query = "TRUNCATE " + tb_name[0] + " RESTART IDENTITY;"
		}

		_, err := migration.Connection().Db.Exec(query)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		fmt.Println(string(green), "success", string(reset), "down ", migrate)
	}
}

func upHandle(argument []string, upMigration, upSeeder *flag.FlagSet, tb *Tables) {
	switch argument[2] {
	case "migration":
		upMigration.Parse(os.Args[3:])
		migrationUp(tb)
	case "seeder":
		upSeeder.Parse(os.Args[3:])
		seederUp(tb)
	default:
		hlp := &Helper{
			Error: os.Args[1],
		}
		printTemplate(errorTmp, hlp)
	}
}

func migrationUp(t *Tables) {
	files, err := ioutil.ReadDir(MigrationFolder)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	if len(t.NameTable) == 0 {
		newFile := []string{}
		for _, file := range files {
			filename := file.Name()
			if filename != "0.go" {
				newFile = append(newFile, filename)
			}
		}
		t.NameTable = newFile
	}
	m := migration.Migration{}
	for _, migrate := range t.NameTable {
		list := strings.Split(migrate, "_migration_")
		tb_name := strings.Split(list[1], ".go")
		meth := reflect.ValueOf(&m).MethodByName("Up" + strings.Title(tb_name[0]))
		meth.Call(nil)
	}
}

func seederUp(t *Tables) {
	files, err := ioutil.ReadDir(SeederFolder)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	if len(t.NameTable) == 0 {
		newFile := []string{}
		for _, file := range files {
			filename := file.Name()
			if filename != "0.go" {
				newFile = append(newFile, filename)
			}
		}
		t.NameTable = newFile
	}
	s := seeder.Seeder{}
	for _, seed := range t.NameTable {
		list := strings.Split(seed, "_seeder_")
		func_name := strings.Split(list[1], ".go")
		meth := reflect.ValueOf(&s).MethodByName(strings.Title(func_name[0]))
		meth.Call(nil)
	}
}
