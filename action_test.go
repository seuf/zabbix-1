package zabbix_test

import (
	"fmt"
	"math/rand"
	"testing"

	. "."
)

func CreateAction(hostGroup *HostGroup, t *testing.T) *Action {
	name := fmt.Sprintf("%s-%d", "Action-test", rand.Int())
	actions := Actions{{

		Name:        name,
		EventSource: SourceActiveAgent,
		EscPeriod:   "60s",
		Operations: Operations{
			{OperationType: ADD_HOST},
			{OperationType: ADD_TO_HOST_GROUP, OpGroup: OpGroups{{GroupId: hostGroup.GroupId}}},
		},
	}}
	err := getAPI(t).ActionsCreate(actions)
	if err != nil {
		t.Fatal(fmt.Sprintf("Action Error : %s", err))
	}
	return &actions[0]
}

func DeleteAction(action *Action, t *testing.T) {
	err := getAPI(t).ActionsDelete(Actions{*action})
	if err != nil {
		t.Fatal(err)
	}
}

func TestActions(t *testing.T) {
	api := getAPI(t)

	//actions, err := api.ActionGet(Params{})
	//if err != nil {
	//	t.Fatal(err)
	//}

	hostGroup := CreateHostGroup(t)
	//defer DeleteHostGroup(hostGroup, t)

	action := CreateAction(hostGroup, t)
	if action.ActionId == "" || action.Name == "" {
		t.Errorf("Something is empty: %#v", action)
	}
	action2, err := api.ActionGetById(action.ActionId)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Action Created : %s", action2.ActionId)

	//DeleteAction(action2, t)
}
