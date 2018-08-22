package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/sung1011/tk-log"
)

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "join roll game",
	Long:  `join roll game`,
	Example: `
		join -j {工号} -m {今天美餐号}
	`,
	Run: func(cmd *cobra.Command, args []string) {
		join(args)
	},
}

//员工号
var joinJbn int

//今日美餐号
var joinMcn int

func init() {
	rootCmd.AddCommand(joinCmd)

	joinCmd.Flags().IntP("jbn", "j", 0, "工号")
	viper.BindPFlag("join-jbn", joinCmd.Flags().Lookup("jbn"))

	joinCmd.Flags().IntP("mcn", "m", 0, "美餐号")
	viper.BindPFlag("join-mcn", joinCmd.Flags().Lookup("mcn"))
}

//Info 信息
type Info struct {
	gorm.Model
	JbNum int `gorm:"unique"`
	McNum int
	Date  int
}

func newInfo() *Info {
	p := new(Info)
	p.JbNum = joinJbn
	p.McNum = joinMcn
	p.Date = Today()
	return p
}

func Today() int {
	t := time.Now()
	s := strconv.Itoa(t.Year()) + fmt.Sprintf("%0.2d", int(t.Month())) + fmt.Sprintf("%0.2d", t.Day())
	date, _ := strconv.Atoi(s)
	return date
}

func join(args []string) {
	_joinChkFlag()
	_doJoin()
}

func _joinChkFlag() {
	joinJbn = viper.GetInt("join-jbn")
	if joinJbn == 0 {
		log.Erro("缺少工号")
	}
	joinMcn = viper.GetInt("join-mcn")
	if joinMcn == 0 {
		log.Erro("缺少今日美餐订餐号")
	}
}

func _doJoin() {
	db, err := gorm.Open("sqlite3", "/tmp/meicanroll.db")
	if err != nil {
		log.Erro(err, db)
	}
	defer db.Close()

	var bi BaseInfo
	if db.Where("ID = ?", joinJbn).First(&bi).RecordNotFound() {
		log.Erro(fmt.Sprintf("未注册得员工号【%v】", joinJbn))
	}

	db.AutoMigrate(&Info{})
	if !db.HasTable(&Info{}) {
		db.CreateTable(&Info{})
	}
	var in Info
	i := newInfo()
	if db.Where("ID = ?", joinJbn).First(&in).RecordNotFound() {
		err = db.Create(&i).Error
		if err != nil {
			log.Erro(err)
		}
	} else {
		db.Save(&i)
	}
}
