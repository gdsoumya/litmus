package gql

import (
	"encoding/json"
	"github.com/gdsoumya/workflow_manager/pkg/types"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Subscription(server, cid, accessKey string) {
	client := &http.Client{}
	clusterID := `{cluster_id: \"` + cid + `\", access_key: \"` + accessKey + `\"}`
	query := `{
			"query": "subscription { clusterConnect(clusterInfo: ` + clusterID + `) { project_id action }  }"
		}`

	req, err := http.NewRequest("POST", server, strings.NewReader(query))
	if err != nil {
		log.Fatal(err)
	}

	// Headers for the calling the server endpoint
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")

	for {
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			log.Fatal(err)
		}

		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		/*
			The format of response is data ---map--> clusterSubsription --map--> data(string)
			- Unmarshal the byte response and store it in a interface
		*/
		var responseInterface map[string]map[string]map[string]string
		err = json.Unmarshal([]byte(bodyText), &responseInterface)
		if err != nil || len(responseInterface) == 0 {
			log.Print(err)
		}
		podRequest := types.PodLogRequest{
			RequestID: responseInterface["data"]["clusterConnect"]["project_id"],
		}
		err = json.Unmarshal([]byte(responseInterface["data"]["clusterConnect"]["action"]), &podRequest)
		if err != nil {
			log.Print(err)
		}
		SendPodLogs(server, cid, accessKey, podRequest)
	}

}
