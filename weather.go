package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func init() {
	key := commandFunc{"!weather", "Gives a loot forcast", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if command == key.command {
			resp, err := http.Get("http://clashofclansforecaster.com/STATS.json")
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()
			dec := json.NewDecoder(resp.Body)
			var f forcast
			if err := dec.Decode(&f); err != nil {
				return "", err
			}
			direction := ""
			switch {
			case f.CurrentLoot.Trend > 0:
				direction = "upwards"
			case f.CurrentLoot.Trend < 0:
				direction = "downward"
			default:
				direction = "stable"
			}
			return fmt.Sprintf("**Loot index on a 10 scale is %s and trend is %s**\n\n", f.LootIndexString, direction) +
				strings.Replace(f.ForecastMessages.English, ".  ", ".\n", -1), nil
		}
		return "", nil
	}
}

type forcast struct {
	CurrentLoot struct {
		TotalPlayers            int `json:"totalPlayers"`
		Trend                   int `json:"trend"`
		LootMinutes             int `json:"lootMinutes"`
		LootMinuteChange        int `json:"lootMinuteChange"`
		PlayersOnline           int `json:"playersOnline"`
		PlayersOnlineChange     int `json:"playersOnlineChange"`
		ShieldedPlayers         int `json:"shieldedPlayers"`
		ShieldedPlayersChange   int `json:"shieldedPlayersChange"`
		AttackablePlayers       int `json:"attackablePlayers"`
		AttackablePlayersChange int `json:"attackablePlayersChange"`
	} `json:"currentLoot"`
	//MainColorShadeNow string `json:"mainColorShadeNow"`
	LootIndexString string `json:"lootIndexString"`
	/*BgColor           string `json:"bgColor"`
	FgColor           string `json:"fgColor"`*/
	ForecastWordNow  string `json:"forecastWordNow"`
	ForecastMessages struct {
		English string `json:"english"`
		/*Spanish     string `json:"spanish"`
		Portuguese  string `json:"portuguese"`
		French      string `json:"french"`
		German      string `json:"german"`
		Indonesian  string `json:"indonesian"`
		Dutch       string `json:"dutch"`
		Finnish     string `json:"finnish"`
		Italian     string `json:"italian"`
		Russian     string `json:"russian"`
		Norwegian   string `json:"norwegian"`
		ChineseSimp string `json:"chinese-simp"`
		ChineseTrad string `json:"chinese-trad"`
		Japanese    string `json:"japanese"`
		Arabic      string `json:"arabic"`
		Persian     string `json:"persian"`
		Korean      string `json:"korean"`
		Hindi       string `json:"hindi"`*/
	} `json:"forecastMessages"`
	/*RegionStats []struct {
		Num0  int    `json:"0"`
		Num1  string `json:"1"`
		Num2  string `json:"2"`
		Num3  string `json:"3"`
		Num4  string `json:"4"`
		Num5  int    `json:"5"`
		Num6  int    `json:"6"`
		Num7  int    `json:"7"`
		Num8  int    `json:"8"`
		Num9  int    `json:"9"`
		Num10 int    `json:"10"`
		Num11 int    `json:"11"`
		Num12 int    `json:"12"`
		Num13 int    `json:"13"`
	} `json:"regionStats"`*/
}
