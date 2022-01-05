package errorinfo

import (
	"log"
	"testing"
)

func Test(t *testing.T) {
	t.Run("Rank", func(t *testing.T) {
		success := map[string]int{
			"a": 1,
			"b": 2,
		}
		log.Println(TurnRank("2067615", success))
	})
}
