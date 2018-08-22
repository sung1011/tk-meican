package cmd

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/sung1011/tk-log"
)

// regiCmd represents the regi command
var regiCmd = &cobra.Command{
	Use:   "regi",
	Short: "register",
	Long:  `register`,
	Example: `
		regi -j {工号} -n {名字}
	`,
	Run: func(cmd *cobra.Command, args []string) {
		register(args)
	},
}

//姓名
var regiName string

//员工号
var regiJbn int

func init() {
	rootCmd.AddCommand(regiCmd)

	regiCmd.Flags().StringP("name", "n", "", "名字")
	viper.BindPFlag("regi-name", regiCmd.Flags().Lookup("name"))

	regiCmd.Flags().IntP("jbn", "j", 0, "工号")
	viper.BindPFlag("regi-jbn", regiCmd.Flags().Lookup("jbn"))
}

func register(args []string) {
	_regiChkFlag()
	_doRegi()
}

func _regiChkFlag() {
	regiName = viper.GetString("regi-name")
	if regiName == "" {
		log.Erro("缺少名字")
	}
	regiJbn = viper.GetInt("regi-jbn")
	if regiJbn == 0 {
		log.Erro("缺少工号")
	}
}

//BaseInfo 基础信息
//ID 为工号
//MeiCanNum 为美餐订单号
type BaseInfo struct {
	gorm.Model
	Name string
}

func newBaseInfo(jn int, name string) *BaseInfo {
	p := new(BaseInfo)
	p.ID = uint(jn)
	p.Name = name
	return p
}

func _doRegi() {
	db, err := gorm.Open("sqlite3", "/tmp/meicanroll.db")
	if err != nil {
		log.Erro(err, db)
	}
	defer db.Close()
	db.AutoMigrate(&BaseInfo{})
	if !db.HasTable(&BaseInfo{}) {
		db.CreateTable(&BaseInfo{})
	}

	bi := newBaseInfo(regiJbn, regiName)
	err = db.Create(&bi).Error
	if err != nil {
		log.Erro(fmt.Sprintf("该工号【%v】已被注册", regiJbn))
	}
}
