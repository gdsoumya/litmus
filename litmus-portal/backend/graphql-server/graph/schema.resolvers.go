package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	store "github.com/litmuschaos/litmus/litmus-portal/backend/graphql-server/pkg/data-store"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/litmuschaos/litmus/litmus-portal/backend/graphql-server/graph/generated"
	"github.com/litmuschaos/litmus/litmus-portal/backend/graphql-server/graph/model"
	"github.com/litmuschaos/litmus/litmus-portal/backend/graphql-server/pkg/cluster"
	"github.com/litmuschaos/litmus/litmus-portal/backend/graphql-server/pkg/database"
	"github.com/litmuschaos/litmus/litmus-portal/backend/graphql-server/pkg/logs"
	"github.com/litmuschaos/litmus/litmus-portal/backend/graphql-server/pkg/workflow"
)

func (r *mutationResolver) UserClusterReg(ctx context.Context, clusterInput model.ClusterInput) (string, error) {
	return cluster.ClusterRegister(clusterInput)
}

func (r *mutationResolver) ClusterConfirm(ctx context.Context, identity model.ClusterIdentity) (string, error) {
	return cluster.ConfirmClusterRegistration(identity, store.State)
}

func (r *mutationResolver) NewClusterEvent(ctx context.Context, clusterEvent model.ClusterEventInput) (string, error) {
	return cluster.NewEvent(clusterEvent, store.State)
}

func (r *mutationResolver) CreateChaosWorkFlow(ctx context.Context, input *model.ChaosWorkFlowInput) (*model.ChaosWorkFlowResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChaosWorkflowRun(ctx context.Context, workflowData model.WorkflowRunInput) (string, error) {
	return workflow.UpsterWorkFlowRun(workflowData, store.State)
}

func (r *mutationResolver) PodLog(ctx context.Context, log model.PodLog) (string, error) {
	return logs.LogsHandler(log)
}

func (r *queryResolver) GetWorkFlowRuns(ctx context.Context, projectID string) ([]*model.WorkflowRun, error) {
	return workflow.QueryWorkflowRuns(projectID)
}

func (r *subscriptionResolver) ClusterEventListener(ctx context.Context, projectID string) (<-chan *model.ClusterEvent, error) {
	log.Print("NEW EVENT LISTENER ", projectID)
	clusterEvent := make(chan *model.ClusterEvent, 1)
	store.State.Mutex.Lock()
	store.State.ClusterEventPublish[projectID] = append(store.State.ClusterEventPublish[projectID], clusterEvent)
	store.State.Mutex.Unlock()
	go func() {
		<-ctx.Done()
	}()
	return clusterEvent, nil
}

func (r *subscriptionResolver) ClusterConnect(ctx context.Context, clusterInfo model.ClusterIdentity) (<-chan *model.ClusterAction, error) {
	log.Print("CLUSTER CONNECT REQUEST")
	clusterAction := make(chan *model.ClusterAction, 1)
	verifiedCluster, err := cluster.VerifyCluster(clusterInfo)
	if err != nil {
		return clusterAction, err
	}
	store.State.Mutex.Lock()
	if _, ok := store.State.ConnectedCluster[clusterInfo.ClusterID]; ok {
		store.State.Mutex.Unlock()
		return clusterAction, errors.New("CLUSTER ALREADY CONNECTED")
	}
	store.State.ConnectedCluster[clusterInfo.ClusterID] = clusterAction
	store.State.Mutex.Unlock()
	go func() {
		<-ctx.Done()
		verifiedCluster.IsActive = false
		cluster.SendClusterEvent("cluster-status", "Cluster Offline", "Cluster Disconnect", *verifiedCluster, store.State)
		store.State.Mutex.Lock()
		delete(store.State.ConnectedCluster, clusterInfo.ClusterID)
		store.State.Mutex.Unlock()
		err = database.UpdateClusterData(clusterInfo.ClusterID, "is_active", false, strconv.FormatInt(time.Now().Unix(), 10))
		if err != nil {
			log.Print("ERROR", err)
		}
	}()
	err = database.UpdateClusterData(clusterInfo.ClusterID, "is_active", true, strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		log.Print("ERROR", err)
		return clusterAction, err
	}
	verifiedCluster.IsActive = true
	cluster.SendClusterEvent("cluster-status", "Cluster Live", "Cluster is Live and Connected", *verifiedCluster, store.State)
	return clusterAction, nil
}

func (r *subscriptionResolver) WorkflowEventListener(ctx context.Context, projectID string) (<-chan *model.WorkflowRun, error) {
	log.Print("NEW WORKFLOW EVENT LISTENER", projectID)
	workflowEvent := make(chan *model.WorkflowRun, 1)
	store.State.Mutex.Lock()
	store.State.WorkflowEventPublish[projectID] = append(store.State.WorkflowEventPublish[projectID], workflowEvent)
	store.State.Mutex.Unlock()
	go func() {
		<-ctx.Done()
		log.Print("CLOSED WORKFLOW LISTENER", projectID)
	}()
	return workflowEvent, nil
}

func (r *subscriptionResolver) GetPodLog(ctx context.Context, podDetails model.PodLogRequest) (<-chan *model.PodLogResponse, error) {
	log.Print("NEW LOG REQUEST", podDetails.ClusterID, podDetails.PodName)
	workflowLog := make(chan *model.PodLogResponse, 1)
	cid := uuid.New()
	store.State.Mutex.Lock()
	store.State.WorkflowLog[cid.String()] = workflowLog
	store.State.Mutex.Unlock()
	go func() {
		<-ctx.Done()
		log.Print("CLOSED LOG LISTENER", podDetails.ClusterID, podDetails.PodName)
		delete(store.State.WorkflowLog, cid.String())
	}()
	go logs.GetLogs(cid.String(), podDetails)
	return workflowLog, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
