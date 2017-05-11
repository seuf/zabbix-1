package zabbix

import "github.com/AlekSi/reflector"

type (
	ObjectType     int
	SourceType     int
	EventValueType int
)

const (
	ObjectTrigger            ObjectType = 0
	ObjectDiscoveredHost     ObjectType = 1
	ObjectDiscoveredService  ObjectType = 2
	ObjectAutoRegisteredHost ObjectType = 3
	ObjectItem               ObjectType = 4
	ObjectLLDRule            ObjectType = 5
)

const (
	SourceTrigger       SourceType = 0
	SourceDiscoveryRule SourceType = 1
	SourceActiveAgent   SourceType = 2
	SourceInternal      SourceType = 3
)

const (
	SourceTriggerOK               EventValueType = 0
	SourceTriggerProblem          EventValueType = 1
	SourceDiscoveryHostUp         EventValueType = 0
	SourceDiscoveryHostDown       EventValueType = 1
	SourceDiscoveryHostDiscovered EventValueType = 2
	SourceDiscoveryHostLost       EventValueType = 3
	SourceInternalNormal          EventValueType = 0
	SourceInternalUnknown         EventValueType = 1
)

type Acknowledge struct {
	AcknowledgeId string `json:"acknowledgeid"`
	Alias         string `json:"alias,omitempty"`
	Clock         uint   `json:"clock"`
	EventId       string `json:"eventid"`
	Message       string `json:"message"`
	Name          string `json:"name,omitempty"`
	SurName       string `json:"surname,omitempty"`
	UserId        string `json:"userid,omitempty"`
}

type Acknowledges []Acknowledge

type Event struct {
	Acknowledged int          `json:"acknowledged"`
	Acknowledges Acknowledges `json:"acknowledges,omitempty"`
	Clock        uint         `json:"clock"`
	EventId      string       `json:"eventid"`
	Ns           uint64       `json:"ns"`
	Object       ObjectType   `json:"object"`
	ObjectId     string       `json:"objectid"`
	Source       SourceType   `json:"source"`
	Value        ValueType    `json:"value"`
	Triggers     Triggers     `json:"triggers,omitempty"`
}

type Events []Event

func (api *API) EventsGet(params Params) (res Events, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("event.get", params)
	if err != nil {
		return
	}
	res = make(Events, len(response.Result.([]interface{})))
	for i, h := range response.Result.([]interface{}) {
		h2 := h.(map[string]interface{})
		reflector.MapToStruct(h2, &res[i], reflector.Strconv, "json")

		if triggers, ok := h2["triggers"]; ok {
			reflector.MapsToStructs2(triggers.([]interface{}), &res[i].Triggers, reflector.Strconv, "json")
		}
		if acknowledges, ok := h2["acknowledges"]; ok {
			reflector.MapsToStructs2(acknowledges.([]interface{}), &res[i].Acknowledges, reflector.Strconv, "json")
		}
	}

	return
}
