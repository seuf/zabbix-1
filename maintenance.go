package zabbix

import "github.com/AlekSi/reflector"

// Maintenance struct - https://www.zabbix.com/documentation/2.4/manual/api/reference/maintenance/object#maintenance
type Maintenance struct {
	MaintenanceID   string      `json:"maintenanceid"`
	Name            string      `json:"name"`
	ActiveSince     int64       `json:"active_since"`
	ActiveTill      int64       `json:"active_till"`
	Description     string      `json:"description"`
	MaintenanceType int         `json:"maintenance_type"`
	TimePeriods     TimePeriods `json:"timeperiods,omitempty"`
}

// TimePeriod struct - https://www.zabbix.com/documentation/2.4/manual/api/reference/maintenance/object#time_period
type TimePeriod struct {
	TimePeriodID   string `json:"timeperiodid"`
	Day            string `json:"day"`
	DayOfWeek      int    `json:"dayofweek"`
	Every          int    `json:"every"`
	Month          int    `json:"month"`
	Period         int64  `json:"period"`
	StartDate      int64  `json:"start_date"`
	StartTime      int64  `json:"start_time"`
	TimePeriodType int    `json:"timeperiod_type"`
}

// Maintenances slice struct for storing result returned from get method
type Maintenances []Maintenance

// TimePeriods slice struct for storing result returned from get method
type TimePeriods []TimePeriod

// MaintenanceGet returns all available maintenances according to given parameters -
// https://www.zabbix.com/documentation/2.4/manual/api/reference/maintenance/get
func (api *API) MaintenanceGet(params Params) (res Maintenances, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("maintenance.get", params)
	if err != nil {
		return
	}
	res = make(Maintenances, len(response.Result.([]interface{})))
	for i, h := range response.Result.([]interface{}) {
		h2 := h.(map[string]interface{})
		reflector.MapToStruct(h2, &res[i], reflector.Strconv, "json")

		if timeperiods, ok := h2["timeperiods"]; ok {
			reflector.MapsToStructs2(timeperiods.([]interface{}), &res[i].TimePeriods, reflector.Strconv, "json")
		}

	}

	return
}
