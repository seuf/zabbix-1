package zabbix

import "github.com/mitchellh/mapstructure"

// https://www.zabbix.com/documentation/2.4/manual/api/reference/history/object
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

	mapstructure.Decode(response.Result.([]interface{}), &res)

	return
}
