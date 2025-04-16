package main

import (
	"errors"
	"github.com/gookit/color"
	"github.com/obnahsgnaw/pbhttp/core/application"
	config2 "github.com/obnahsgnaw/pbhttp/core/config"
	gen2 "github.com/obnahsgnaw/socketgwservice/internal/dal/gen"
	"github.com/spf13/viper"
	"gorm.io/gen"
	"gorm.io/gorm"
	"os"
)

type Config struct {
	IniFile  string            `short:"c" long:"conf" description:"Ini file"`
	Database *config2.Database `ini:"database"`
}

func parse() *config2.Database {
	var opt Config

	if _, err := flags.Parse(&opt); err != nil {
		var flagErr *flags.Error
		ok := errors.As(err, &flagErr)
		if !ok || !errors.Is(flagErr.Type, flags.ErrHelp) {
			color.Error.Println("input error")
		}
		os.Exit(0)
	}

	if opt.IniFile != "" {
		cc := viper.New()
		cc.SetConfigFile(opt.IniFile)
		if err := cc.ReadInConfig(); err != nil {
			color.Error.Println(err.Error())
			os.Exit(2)
		}
		if err := cc.Unmarshal(&opt); err != nil {
			color.Error.Println(err.Error())
			os.Exit(2)
		}
	}

	return opt.Database
}

func main() {
	cnf := parse()
	var dbConn *gorm.DB
	var err error
	if dbConn, err = application.NewDb(cnf); err != nil {
		panic(err)
	}

	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/dal/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery,
		//if you want the nullable field generation property to be pointer type, set FieldNullable true
		FieldNullable: true,
		//if you want to assign field which has default value in `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		FieldCoverable: true,
		// if you want to generate field with unsigned integer type, set FieldSignable true
		FieldSignable: true,
		//if you want to generate index tags from database, set FieldWithIndexTag true
		FieldWithIndexTag: true,
		//if you want to generate type tags from database, set FieldWithTypeTag true
		FieldWithTypeTag: true,
		//if you need unit tests for query code, set WithUnitTest true
		/* WithUnitTest: true, */
	})

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary or it will panic
	// db, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(dbConn)

	// apply basic crud api on structs or table models which is specified by table name with function
	// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Excute.
	// 想对已有的model生成crud等基础方法可以直接指定model struct ，例如model.User{}
	// 如果是想直接生成表的model和crud方法，则可以指定表的名称，例如g.GenerateModel("company")
	// 想自定义某个表生成特性，比如struct的名称/字段类型/tag等，可以指定opt，例如g.GenerateModel("company",gen.FieldIgnore("address")), g.GenerateModelAs("people", "Person", gen.FieldIgnore("address"))
	// g.ApplyBasic(model.User{})
	// g.GenerateModel("company")
	// g.GenerateModelAs("people", "Person", gen.FieldIgnore("address")))

	// apply diy interfaces on structs or table models
	// 如果想给某些表或者model生成自定义方法，可以用ApplyInterface，第一个参数是方法接口，可以参考DIY部分文档定义
	//g.ApplyInterface(func(method model.Method) {}, model.User{}, g.GenerateModel("company"))

	gen2.InitModel(g)

	// execute the action of code generation
	g.Execute()
}
