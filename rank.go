package rank

import "log"

func RankRunner(taskAttributes []string, runnerAttributes []string) {
	// ranking := []string{}
	matchCount := 0
	for _, taskAttr := range taskAttributes {
		for _, runnerAttr := range runnerAttributes {
			if runnerAttr == taskAttr {
				log.Printf("there's a match %v", runnerAttr)
				// ranking = append(ranking, runnerAttr)
				matchCount++
			}
		}
	}
	ranked := matchCount >= len(taskAttributes)/2
	if ranked {
		log.Printf("there is a %d chance", matchCount*100/len(taskAttributes))
	}else {
		log.Printf("low matching %d", matchCount*100/len(taskAttributes))
	}
}