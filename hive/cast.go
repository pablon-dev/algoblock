package hive

import (
	"errors"
	"fmt"
	"os"

	"github.com/pablon-dev/cubanblock/request"
)

var (
	wif   string
	voter string
)

const factor int = 100

func init() {
	voter = os.Getenv("hive_voter")
	wif = os.Getenv("wif_voter")
}

//CastVote cast a vote to Hive network
func CastVote(notes []string, power []uint64) {
	if len(notes) < 1 || len(power) != len(notes) {
		return
	}
	for i, x := range notes {
		author, permlink, err := GetAuthorAndPermLink(x)
		if err != nil {
			continue
		}
		peso := int(power[i]) * factor
		resp, err := request.Upvote(wif, voter, author, permlink, peso)
		if err != nil {
			continue
		}
		if len(resp.ID) > 0 {
			fmt.Println("Voto efectuado!")
		}

	}

}

//GetAuthorAndPermLink get the author post and permlink for given url
func GetAuthorAndPermLink(link string) (author string, permlink string, err error) {
	i := getindex(link, '@')
	if i+5 > len(link) || i < 5 {
		return "", "", errors.New("Not valid link")
	}
	link = link[i:]
	i = getindex(link, '/')
	if i+5 > len(link) || i < 2 {
		return "", "", errors.New("Not valid link")
	}
	author = link[1:i]
	permlink = link[i+1:]
	err = nil
	return

}
func getindex(text string, char rune) (result int) {
	for i, x := range text {
		if x == ' ' {
			return 0
		}
		if x == char {
			return i
		}

	}
	return 0
}
