package zabbix

import "github.com/mitchellh/mapstructure"

type (
	AvailableType int
	StatusType    int
)

const (
	Available   AvailableType = 1
	Unavailable AvailableType = 2

	Monitored   StatusType = 0
	Unmonitored StatusType = 1

	NotInMaint StatusType = 0
	InMaint    StatusType = 1
)

// Host definition by https://www.zabbix.com/documentation/2.4/manual/api/reference/host/object
type Host struct {
	HostId      string        `json:"hostid,omitempty"`
	Host        string        `json:"host"`
	Available   AvailableType `json:"available"`
	Error       string        `json:"error"`
	Name        string        `json:"name"`
	Status      StatusType    `json:"status"`
	MaintStatus StatusType    `json:"maintenance_status"`
	// Fields below used only when creating hosts
	GroupIds    HostGroupIds   `json:"groups,omitempty"`
	Interfaces  HostInterfaces `json:"interfaces,omitempty"`
	TemplateIds TemplateIds    `json:"templates,omitempty"`
}

type HostId struct {
	HostId string `json:"hostid"`
}

//type HostIDs []HostID

type Hosts []Host

// HostsGet is a wrapper for host.get: https://www.zabbix.com/documentation/2.2/manual/appendix/api/host/get
func (api *API) HostsGet(params Params) (res Hosts, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("host.get", params)
	if err != nil {
		return
	}

	mapstructure.Decode(response.Result.([]interface{}), &res)

	return
}

// Gets hosts by host group Ids.
func (api *API) HostsGetByHostGroupIds(ids []string) (res Hosts, err error) {
	return api.HostsGet(Params{"groupids": ids})
}

// Gets hosts by host groups.
func (api *API) HostsGetByHostGroups(hostGroups HostGroups) (res Hosts, err error) {
	ids := make([]string, len(hostGroups))
	for i, id := range hostGroups {
		ids[i] = id.GroupId
	}
	return api.HostsGetByHostGroupIds(ids)
}

// Gets host by Id only if there is exactly 1 matching host.
func (api *API) HostGetById(id string) (res *Host, err error) {
	hosts, err := api.HostsGet(Params{"hostids": id})
	if err != nil {
		return
	}

	if len(hosts) == 1 {
		res = &hosts[0]
	} else {
		e := ExpectedOneResult(len(hosts))
		err = &e
	}
	return
}

// HostGetByHost gets host by Host only if there is exactly 1 matching host.
func (api *API) HostGetByHost(host string) (res *Host, err error) {
	hosts, err := api.HostsGet(Params{"filter": map[string]string{"host": host}})
	if err != nil {
		return
	}

	if len(hosts) == 1 {
		res = &hosts[0]
	} else {
		e := ExpectedOneResult(len(hosts))
		err = &e
	}
	return
}

// Wrapper for host.create: https://www.zabbix.com/documentation/2.2/manual/appendix/api/host/create
func (api *API) HostsCreate(hosts Hosts) (err error) {
	response, err := api.CallWithError("host.create", hosts)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	hostids := result["hostids"].([]interface{})
	for i, id := range hostids {
		hosts[i].HostId = id.(string)
	}
	return
}

// Wrapper for host.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/host/delete
// Cleans HostId in all hosts elements if call succeed.
func (api *API) HostsDelete(hosts Hosts) (err error) {
	ids := make([]string, len(hosts))
	for i, host := range hosts {
		ids[i] = host.HostId
	}

	err = api.HostsDeleteByIds(ids)
	if err == nil {
		for i := range hosts {
			hosts[i].HostId = ""
		}
	}
	return
}

// Wrapper for host.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/host/delete
func (api *API) HostsDeleteByIds(ids []string) (err error) {
	hostIds := make([]map[string]string, len(ids))
	for i, id := range ids {
		hostIds[i] = map[string]string{"hostid": id}
	}

	response, err := api.CallWithError("host.delete", hostIds)
	if err != nil {
		// Zabbix 2.4 uses new syntax only
		if e, ok := err.(*Error); ok && e.Code == -32500 {
			response, err = api.CallWithError("host.delete", ids)
		}
	}
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	hostids := result["hostids"].([]interface{})
	if len(ids) != len(hostids) {
		err = &ExpectedMore{len(ids), len(hostids)}
	}
	return
}
