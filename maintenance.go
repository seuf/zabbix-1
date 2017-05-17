package zabbix

import (
	"github.com/AlekSi/reflector"
)

type (
	MaintType  int
	PeriodType int
)

const (
	WithData    MaintType = 0
	WithoutData MaintType = 1

	OneTime PeriodType = 0
	Daily   PeriodType = 2
	Weekly  PeriodType = 3
	Monthly PeriodType = 4
)

// Maintenance struct - https://www.zabbix.com/documentation/2.4/manual/api/reference/maintenance/object#maintenance
type Maintenance struct {
	MaintenanceID   string       `json:"maintenanceid"`
	Name            string       `json:"name"`
	ActiveSince     int64        `json:"active_since"`
	ActiveTill      int64        `json:"active_till"`
	Description     string       `json:"description,omitempty"`
	MaintenanceType MaintType    `json:"maintenance_type"`
	HostIDs         HostIDs      `json:"hostids"`
	HostGroupIDs    HostGroupIds `json:"groupids"`
	Hosts           Hosts        `json:"hosts,omitempty"`
	HostGroups      HostGroups   `json:"groups,omitempty"`
	TimePeriods     TimePeriods  `json:"timeperiods,omitempty"`
}

// TimePeriod struct - https://www.zabbix.com/documentation/2.4/manual/api/reference/maintenance/object#time_period
type TimePeriod struct {
	TimePeriodID   string     `json:"timeperiodid"`
	Day            string     `json:"day,omitempty"`
	DayOfWeek      int        `json:"dayofweek,omitempty"`
	Every          int        `json:"every,omitempty"`
	Month          int        `json:"month,omitempty"`
	Period         int64      `json:"period,omitempty"`
	StartDate      int64      `json:"start_date,omitempty"`
	StartTime      int64      `json:"start_time,omitempty"`
	TimePeriodType PeriodType `json:"timeperiod_type"`
}

// Maintenances slice struct for storing result returned from get method
type Maintenances []Maintenance

// TimePeriods slice struct for storing result returned from get method
type TimePeriods []TimePeriod

// MaintenancesGet returns all available maintenances according to given parameters -
// https://www.zabbix.com/documentation/2.4/manual/api/reference/maintenance/get
func (api *API) MaintenancesGet(params Params) (res Maintenances, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	if _, present := params["selectHosts"]; !present {
		params["selectHosts"] = "extend"
	}
	if _, present := params["selectGroups"]; !present {
		params["selectGroups"] = "extend"
	}
	if _, present := params["selectTimeperiods"]; !present {
		params["selectTimeperiods"] = "extend"
	}
	response, err := api.CallWithError("maintenance.get", params)
	if err != nil {
		return
	}
	res = make(Maintenances, len(response.Result.([]interface{})))
	for i, h := range response.Result.([]interface{}) {
		h2 := h.(map[string]interface{})
		reflector.MapToStruct(h2, &res[i], reflector.Strconv, "json")

		if hosts, ok := h2["hosts"]; ok {
			reflector.MapsToStructs2(hosts.([]interface{}), &res[i].Hosts, reflector.Strconv, "json")
		}
		if hostgroups, ok := h2["groups"]; ok {
			reflector.MapsToStructs2(hostgroups.([]interface{}), &res[i].HostGroups, reflector.Strconv, "json")
		}
		if timeperiods, ok := h2["timeperiods"]; ok {
			reflector.MapsToStructs2(timeperiods.([]interface{}), &res[i].TimePeriods, reflector.Strconv, "json")
		}

	}

	return
}

// MaintenanceGetByID returns maintenance by ID only if there is exactly 1 matching maintenance.
func (api *API) MaintenanceGetByID(id string) (res *Maintenance, err error) {
	maintenances, err := api.MaintenancesGet(Params{"maintenanceids": id})
	if err != nil {
		return
	}

	if len(maintenances) == 1 {
		res = &maintenances[0]
	} else {
		e := ExpectedOneResult(len(maintenances))
		err = &e
	}
	return
}

// MaintenanceGetByName returns maintenance by its name only if there is exactly 1 matching maintenance.
func (api *API) MaintenanceGetByName(name string) (res *Maintenance, err error) {
	maintenances, err := api.MaintenancesGet(Params{"filter": map[string]string{"name": name}})
	if err != nil {
		return
	}
	if len(maintenances) == 1 {
		res = &maintenances[0]
	} else {
		e := ExpectedOneResult(len(maintenances))
		err = &e
	}
	return
}

// MaintenancesCreate creates maintenances using maintenance.create - https://www.zabbix.com/documentation/2.4/manual/api/reference/maintenance/create
func (api *API) MaintenancesCreate(maintenances Maintenances) (err error) {
	response, err := api.CallWithError("maintenance.create", maintenances)
	if err != nil {
		return
	}
	result := response.Result.(map[string]interface{})
	maintenanceids := result["maintenanceids"].([]interface{})
	for i, id := range maintenanceids {
		maintenances[i].MaintenanceID = id.(string)
	}
	return
}

func (api *API) MaintenancesUpdate(maintenances Maintenances) (err error) {
	return
}

func (api *API) MaintenancesDelete(maintenances Maintenances) (err error) {
	return
}
