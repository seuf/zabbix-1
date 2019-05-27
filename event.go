package zabbix

import "github.com/mitchellh/mapstructure"

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

// Acknowledge struct from https://www.zabbix.com/documentation/2.4/manual/api/reference/event/get?s[]=acknowledgeid
type Acknowledge struct {
	AcknowledgeId string `json:"acknowledgeid"`
	Alias         string `json:"alias,omitempty"`
	Clock         int64  `json:"clock"`
	EventId       string `json:"eventid"`
	Message       string `json:"message"`
	Name          string `json:"name,omitempty"`
	SurName       string `json:"surname,omitempty"`
	UserId        string `json:"userid,omitempty"`
}

type Acknowledges []Acknowledge

// Event struct from https://www.zabbix.com/documentation/2.4/manual/api/reference/event/object
type Event struct {
	Acknowledged int          `json:"acknowledged"`
	Acknowledges Acknowledges `json:"acknowledges,omitempty"`
	Clock        int64        `json:"clock"`
	EventId      string       `json:"eventid"`
	Ns           int64        `json:"ns"`
	Object       ObjectType   `json:"object"`
	ObjectId     string       `json:"objectid"`
	Source       SourceType   `json:"source"`
	Value        ValueType    `json:"value"`
	Triggers     Triggers     `json:"triggers,omitempty"`
}

type Events []Event

// EventsGet gets all events https://www.zabbix.com/documentation/2.4/manual/api/reference/event/get
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
		mapstructure.Decode(h2, &res[1])

		if triggers, ok := h2["triggers"]; ok {
			mapstructure.Decode(triggers.([]interface{}), &res[i].Triggers)

		}
		if acknowledges, ok := h2["acknowledges"]; ok {
			mapstructure.Decode(acknowledges.([]interface{}), &res[i].Acknowledges)
		}
	}

	return
}

// EventsGetByID gets an event by item ID
func (api *API) EventsGetByID(id string) (res Events, err error) {
	return api.EventsGet(Params{"eventids": id, "select_acknowledges": "extend"})
}

// EventsGetByTriggerID gets an event by item ID, default source = 0 (triggers)
func (api *API) EventsGetByTriggerID(id string) (res Events, err error) {
	return api.EventsGet(Params{"objectids": id})
}

// EventsAckByID acknowledges event using id and text message - https://www.zabbix.com/documentation/2.4/manual/api/reference/event/acknowledge
func (api *API) EventsAckByID(id string, message string) (err error) {
	_, err = api.CallWithError("event.acknowledge", Params{"eventids": id, "message": message})
	if err != nil {
		return
	}
	return
}
