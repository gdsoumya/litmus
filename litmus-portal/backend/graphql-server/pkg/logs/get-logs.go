package logs

import (
	"encoding/json"
	"github.com/litmuschaos/litmus/litmus-portal/backend/graphql-server/graph/model"
	store "github.com/litmuschaos/litmus/litmus-portal/backend/graphql-server/pkg/data-store"
	"log"
)

func GetLogs(uid string, pod model.PodLogRequest){
	data, err := json.Marshal(pod)
	if err!=nil{
		log.Print("ERROR WHILE MARSHALLING POD DETAILS")
	}
	payload := model.ClusterAction{
		ProjectID: uid,
		Action: string(data),
	}
	if clusterChan,ok:=store.State.ConnectedCluster[pod.ClusterID];ok {
		clusterChan <- &payload
	}
}