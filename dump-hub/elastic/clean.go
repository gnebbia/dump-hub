package elastic

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/olivere/elastic/v7"
)

/*
cleanTmp :: Clean tmp folder
*/
func cleanTmp() {
	log.Println("Cleaning tmp folder...")

	dir, err := ioutil.ReadDir("/tmp")
	if err != nil {
		log.Println(err)
	}
	for _, d := range dir {
		if d.Name() != ".gitkeep" {
			os.Remove(path.Join("/tmp", d.Name()))
		}
	}
}

/*
cleanHistory :: Clean unprocessed files and update history status
*/
func (eClient *Client) cleanHistory() {
	log.Println("Cleaning history of unprocessed files...")

	matchQ := elastic.NewMatchQuery(
		"status",
		0,
	)
	query := elastic.
		NewBoolQuery().
		Must(matchQ)

	scroll := eClient.client.Scroll().
		Index("dump-hub-history").
		Query(query).
		Size(1)

	for {
		result, err := scroll.Do(eClient.ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}

		for _, hit := range result.Hits.Hits {
			err = eClient.UpdateHistoryStatus(hit.Id, -1)
			if err != nil {
				log.Println(err)
			}
		}
	}

	matchQ = elastic.NewMatchQuery(
		"status",
		2,
	)
	query = elastic.
		NewBoolQuery().
		Must(matchQ)

	scroll = eClient.client.Scroll().
		Index("dump-hub-history").
		Query(query).
		Size(1)

	for {
		result, err := scroll.Do(eClient.ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}

		for _, hit := range result.Hits.Hits {
			err = eClient.UpdateHistoryStatus(hit.Id, -1)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
