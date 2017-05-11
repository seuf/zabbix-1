package zabbix

import "github.com/AlekSi/reflector"

type (
	PriorityType int
)

const (
	NotClassified PriorityType = 0
	Information   PriorityType = 1
	Warning       PriorityType = 2
	Average       PriorityType = 3
	High          PriorityType = 4
	Disaster      PriorityType = 5
)

const (
	TriggerOk      ValueType = 0
	TriggerProblem ValueType = 1
)

type Trigger struct {
	TriggerId   string       `json:"triggerid"`
	Description string       `json:"description"`
	Expression  string       `json:"expression"`
	Error       string       `json:"error"`
	Hosts       HostIds      `json:"hosts,omitempty"`
	Priority    PriorityType `json:"priority"`
	Value       ValueType    `json:"value"`
}

type HostId struct {
	HostId string `json:"hostid"`
	Name   string `json:"name"`
}

type HostIds []HostId

type Triggers []Trigger

func (api *API) TriggersGet(params Params) (res Triggers, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("trigger.get", params)
	if err != nil {
		return
	}
	res = make(Triggers, len(response.Result.([]interface{})))
	for i, h := range response.Result.([]interface{}) {
		h2 := h.(map[string]interface{})
		reflector.MapToStruct(h2, &res[i], reflector.Strconv, "json")

		if hosts, ok := h2["hosts"]; ok {
			reflector.MapsToStructs2(hosts.([]interface{}), &res[i].Hosts, reflector.Strconv, "json")
		}
	}

	return
}