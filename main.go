package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	myClient = &http.Client{Timeout: 10 * time.Second}
)

type Spreads struct {
	Odds   []float64
	Points []string
}

type Odds struct {
	Spreads Spreads
}

type Site struct {
	Site_key string
	Odds     Odds
}

type Game struct {
	Id        string
	Teams     []string
	Home_team string
	Sites     []Site
}

type Data struct {
	Success bool
	Data    []Game
}

func addRow(fname string, column []string) {
	// read the file
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	w := csv.NewWriter(f)
	w.Write(column)
	w.Flush()
}

func getOdds(team string, data Data) string {
	team = strings.ToLower(team)

	for _, game := range data.Data {

		team1_name := strings.Split(game.Teams[0], " ")
		team1 := strings.ToLower(team1_name[len(team1_name)-1])
		team2_name := strings.Split(game.Teams[1], " ")
		team2 := strings.ToLower(team2_name[len(team2_name)-1])

		if team1 == team {
			return game.Sites[0].Odds.Spreads.Points[0]
		} else if team2 == team {
			return game.Sites[0].Odds.Spreads.Points[1]
		}
	}
	return ""
}

func getOpponent(team string, data Data) string {
	team = strings.ToLower(team)

	for _, game := range data.Data {

		team1_name := strings.Split(game.Teams[0], " ")
		team1 := strings.ToLower(team1_name[len(team1_name)-1])
		team2_name := strings.Split(game.Teams[1], " ")
		team2 := strings.ToLower(team2_name[len(team2_name)-1])

		if team1 == team {
			return team2
		} else if team2 == team {
			return team1
		}
	}
	return ""
}

func main() {
	apiKey := "37783b2b45645133431ca784d46519ab"
	requestURL := "https://api.the-odds-api.com/v3/odds/?apiKey=" + apiKey + "&sport=basketball_nba&region=us&mkt=spreads"

	for true {
		r, err := myClient.Get(requestURL)
		if err == nil {
			defer r.Body.Close()

			data := Data{}
			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				panic(err)
			}

			// parse and store data for teams
			for i := 1; i < len(os.Args); i++ {
				team := os.Args[i]
				line := getOdds(team, data)

				if line != "" {
					fmt.Println(team + ": " + line)
					addRow(team+"&"+getOpponent(team, data)+time.Now().Format("01-02-2006")+".csv", []string{line})
				} else {
					fmt.Println("data not acquired")
				}
			}
			time.Sleep(time.Minute)
		} else {
			fmt.Printf("request error: %v\n", err)
			time.Sleep(time.Second * 10)
		}
	}
}
