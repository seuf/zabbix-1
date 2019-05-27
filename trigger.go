package zabbix

import "github.com/mitchellh/mapstructure"

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
	Hosts       Hosts        `json:"hosts,omitempty"`
	Priority    PriorityType `json:"priority"`
	Value       ValueType    `json:"value"`
}

type Triggers []Trigger

// TriggersGet gets all triggers
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
		mapstructure.Decode(h2, &res[i])

		if hosts, ok := h2["hosts"]; ok {
			mapstructure.Decode(hosts.([]interface{}), &res[i].Hosts)

		}
	}

	return
}

// TriggerGetByID gets trigger by Id only if there is exactly 1 matching trigger.
func (api *API) TriggerGetByID(id string) (res *Trigger, err error) {
	triggers, err := api.TriggersGet(Params{"triggerids": id})
	if err != nil {
		return
	}

	if len(triggers) == 1 {
		res = &triggers[0]
	} else {
		e := ExpectedOneResult(len(triggers))
		err = &e
	}
	return
}
