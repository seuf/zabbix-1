package zabbix

import "github.com/AlekSi/reflector"

type History struct {
	Clock  uint   `json:"clock"`
	ItemId string `json:"itemid"`
	Ns     int    `json:"ns"`
	Value  string `json:"value"` // Currently always returns strings

	Id         string `json:"id,omitempty"`
	LogEventId int    `json:"logeventid,omitempty"`
	Severity   int    `json:"severity,omitempty"`
	Source     string `json:"source,omitempty"`
	Timestamp  uint   `json:"timestamp,omitempty"`
}

type Histories []History

// Wrapper for item.get https://www.zabbix.com/documentation/2.0/manual/appendix/api/item/get
func (api *API) HistoriesGet(params Params) (res Histories, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("history.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}
