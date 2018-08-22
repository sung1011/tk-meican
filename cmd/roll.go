package cmd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	log "github.com/sung1011/tk-log"
)

// rollCmd represents the roll command
var rollCmd = &cobra.Command{
	Use:   "roll",
	Short: "roll",
	Long:  `roll`,
	Run: func(cmd *cobra.Command, args []string) {
		_roll()
	},
}

func init() {
	rootCmd.AddCommand(rollCmd)
}

type rollRow struct {
	Name  string  `json:"名字"`
	Jbnum string  `json:"工号"`
	Mcnum string  `json:"今日美餐号"`
	Roll  float64 `json:"roll点数"`
}

type rollRows []rollRow

func _roll() {
	db, err := gorm.Open("sqlite3", "/tmp/meicanroll.db")
	if err != nil {
		log.Erro(err, db)
	}
	defer db.Close()
	var rrs rollRows
	db.Table("infos").Select("base_infos.name as name, base_infos.id as jbnum, infos.mc_num as mcnum").Joins("left join base_infos on infos.jb_num = base_infos.id").Where("infos.date = ?", Today()).Scan(&rrs)
	if len(rrs) == 0 {
		log.Erro("无人需要代取")
	}
	if len(rrs) == 1 {
		log.Erro("只有1人参加,匿了")
	}

	for k := range rrs {
		rrs[k].Roll = rand.Float64() * 100
	}

	sort.Sort(rrs)

	b, _ := json.MarshalIndent(rrs, "", "\t")
	log.Succ(string(b))
	b, _ = json.Marshal(rrs[0])
	log.Warn(fmt.Sprintf("\r\n恭喜 %s (%s) 您今天取餐", rrs[0].Name, rrs[0].Jbnum))
}

func (r rollRows) Len() int {
	return len(r)
}

func (r rollRows) Less(i, j int) bool {
	return r[i].Roll > r[j].Roll
}

func (r rollRows) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
