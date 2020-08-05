package data_store

import (
	"sync"

	"github.com/litmuschaos/litmus/litmus-portal/backend/graphql-server/graph/model"
)

//Application state, contains channels and mutexes used for subscriptions
type StateData struct {
	ClusterEventPublish map[string][]chan *model.ClusterEvent
	ConnectedCluster    map[string]chan *model.ClusterAction
	Mutex               *sync.Mutex
	WorkflowEventPublish map[string][]chan *model.WorkflowRun
	WorkflowLog map[string]chan *model.PodLogResponse
}

var State StateData

//StoreInit initializes Application State
func StoreInit() {
	State.ClusterEventPublish = make(map[string][]chan *model.ClusterEvent)
	State.WorkflowEventPublish = make(map[string][]chan *model.WorkflowRun)
	State.ConnectedCluster = make(map[string]chan *model.ClusterAction)
	State.WorkflowLog = make(map[string]chan *model.PodLogResponse)
	State.Mutex = &sync.Mutex{}
}
