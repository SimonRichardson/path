package main

import (
	"fmt"
	"log"

	"github.com/spoke-d/path"
	"github.com/spoke-d/path/set"
)

func main() {
	root := set.MakeSet(map[string]interface{}{
		"company": map[string]interface{}{
			"person": map[string]interface{}{
				"name": "fred",
			},
		},
	})

	query, err := path.Parse(`company.person.(name == "fred")`)
	if err != nil {
		log.Fatal(err)
	}

	done, err := query.Run(root)
	fmt.Println(done, err)
}
