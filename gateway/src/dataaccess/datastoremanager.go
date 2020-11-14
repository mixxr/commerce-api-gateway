package dataaccess

import (
	"dataaccess/impl/mydatastore"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func loadConfigurations() *mydatastore.DBConfig {
	runmode, ok := os.LookupEnv("DCGW_RUNMODE")
	if !ok {
		runmode = "dev"
	}

	// Set the fileid name of the configurations file
	viper.SetConfigName("config." + runmode + ".yaml")

	// Set the path to look for the configurations file
	configPath, ok := os.LookupEnv("DCGW_CONFIGPATH")
	if !ok {
		viper.AddConfigPath(".")
	} else {
		viper.AddConfigPath(configPath)
	}

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	dbcfg := mydatastore.DBConfig{
		Uid:    "root",
		Pwd:    "secr3tZ",
		IP:     "127.0.0.1",
		Port:   "3306",
		Dbname: "dcgw",
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Using default DATAACCESS configurations...Error reading config file, %s\n", err)
	} else {
		fmt.Printf("Using env DATAACCESS configurations...\n")
		viper.SetDefault("DATAACCESS.HOST", dbcfg.IP)
		viper.SetDefault("DATAACCESS.PORT", dbcfg.Port)
		viper.SetDefault("DATAACCESS.USERNAME", dbcfg.Uid)
		viper.SetDefault("DATAACCESS.PASSWORD", dbcfg.Pwd)
		viper.SetDefault("DATAACCESS.DBNAME", dbcfg.Dbname)
		dbcfg.IP = viper.Get("DATAACCESS.HOST").(string)
		dbcfg.Port = viper.Get("DATAACCESS.PORT").(string)
		dbcfg.Uid = viper.Get("DATAACCESS.USERNAME").(string)
		dbcfg.Pwd = viper.Get("DATAACCESS.PASSWORD").(string)
		dbcfg.Dbname = viper.Get("DATAACCESS.DBNAME").(string)
	}

	return &dbcfg
}

func NewDatastore() (IDatastore, error) {

	dbcfg := loadConfigurations()
	fmt.Println("Connecting to...", dbcfg.String())

	myDatastore, err := mydatastore.NewDatastore(dbcfg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return myDatastore, nil
}
