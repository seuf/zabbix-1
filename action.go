package zabbix

import (
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type (
	EvalType      int
	OperationType int
	CommandType   int
	AuthType      int
	ExecuteOn     int
)

const (
	ANDOR EvalType = 0
	AND   EvalType = 1
	OR    EvalType = 2
)

const (
	SEND_MESSAGE            OperationType = 0
	REMOTE_COMMAND          OperationType = 1
	ADD_HOST                OperationType = 2
	REMOVE_HOST             OperationType = 3
	ADD_TO_HOST_GROUP       OperationType = 4
	REMOVE_FROM_HOST_GROUP  OperationType = 5
	LINK_TO_TEMPLATE        OperationType = 6
	UNLINK_FROM_TEMPLATE    OperationType = 7
	ENABLE_HOST             OperationType = 8
	DISABLE_HOST            OperationType = 9
	SET_HOST_INVENTORY_MODE OperationType = 10
)

const (
	CUSTOM_SCRIPT CommandType = 0
	IPMI_CMD      CommandType = 1
	SSH           CommandType = 2
	TELNET        CommandType = 3
	GLOBAL_SCRIPT CommandType = 4
)

const (
	PASSWORD   AuthType = 0
	PUBLIC_KEY AuthType = 1
)

const (
	ZABBIX_AGENT  ExecuteOn = 0
	ZABBIX_SERVER ExecuteOn = 1
	ZABBIX_PROXY  ExecuteOn = 2
)

// https://www.zabbix.com/documentation/4.2/manual/appendix/api/action/definitions
type Action struct {
	ActionId        string     `json:"actionid,omitempty"`
	EscPeriod       string     `json:"esc_period,omitempty"`
	EventSource     SourceType `json:"eventsource"`
	Name            string     `json:"name"`
	DefLongdata     string     `json:"def_longdata,omitempty"`
	DefShortdata    string     `json:"def_shortdata,omitempty"`
	RLongData       string     `json:"r_longdata,omitempty"`
	RShortData      string     `json:"r_shortdata,omitempty"`
	AckLongData     string     `json:"ack_longdata,omitempty"`
	AckShortData    string     `json:"ack_shortdata,omitempty"`
	Status          int        `json:"status,omitempty"`
	PauseSuppressed int        `json:"pause_suppressed,omitempty"`
	Operations      Operations `json:"operations"`
}

type Actions []Action

type ActionId struct {
	ActionId string `json:"actionid"`
}

type Operation struct {
	OperationId   string        `json:"operationid,omitempty"`
	OperationType OperationType `json:"operationtype"`
	ActionId      string        `json:"actionid,omitempty"`
	EscPeriod     string        `json:"esc_period,omitempty"`
	EscStepFrom   int           `json:"esc_step_from,omitempty"`
	EscStepTo     int           `json:"esc_step_to,omitempty"`
	EvalType      EvalType      `json:"evaltype,omitempty"`
	OpCommand     OpCommand     `json:"opcommand,omitempty"`
	OpCommandGrp  OpCommandGrps `json:"opcommand_grp,omitempty"`
	OpCommandHst  OpCommandHsts `json:"opcommand_hst,omitempty"`
	OpConditions  OpConditions  `json:"opconditions,omitempty"`
	OpGroup       OpGroups      `json:"opgroup,omitempty"`
	OpMessage     OpMessage     `json:"opmessage,omitempty"`
	OpMessageGrp  OpMessageGrps `json:"opmessage_grp,omitempty"`
	OpMessageUsr  OpMessageUsrs `json:"opmessage_usr,omitempty"`
	OpTemplate    OpTemplates   `json:"optemplate,omitempty"`
	OpInventory   OpInventory   `json:"opinventory,omitempty"`
}

type Operations []Operation

type OperationId struct {
	OperationId string `json:"operationid"`
}

type OpCommand struct {
	OperationId string      `json:"operationid,omitempty"`
	Command     string      `json:"command,omitempty"`
	Type        CommandType `json:"type"`
	AuthType    int         `json:"authtype,omitempty"`
	ExecuteOn   ExecuteOn   `json:"execute_on,omitempty"`
	Password    string      `json:"password,omitempty"`
	Port        string      `json:"port,omitempty"`
	PrivateKey  string      `json:"privatekey,omitempty"`
	PublicKey   string      `json:"publickey,omitempty"`
	ScriptId    string      `json:"scriptid,omitempty"`
	Username    string      `json:"username,omitempty"`
}

type OpCommandGrp struct {
	OpCommandGrpid string      `json:"opcommand_grpid,omitempty"`
	OperationId    OperationId `json:"operationid,omitempty"`
	GroupId        string      `json:"groupid,omitempty"`
}

type OpCommandGrps []OpCommandGrp

type OpCommandHst struct {
	OpCommandHstid string      `json:"opcommand_hstid,omitempty"`
	OperationId    OperationId `json:"operationid,omitempty"`
	HostId         HostId      `json:"hostid,omitempty"`
}
type OpCommandHsts []OpCommandHst

type Opcondition struct {
	OpConditionId   string `json:"opconditionid,omitempty"`
	OpConditionType int    `json:"conditiontype"`
	// Possible values:
	// 14 - event acknowledged.
	Value       string      `json:"value,omitempty"`
	OperationId OperationId `json:"operationid,omitempty"`
	Operator    int         `json:"operator,omitempty"`
}

type OpConditions []Opcondition

type OpGroup struct {
	//OperationId OperationId `json:"operationid"`
	GroupId string `json:"groupid,omitempty"`
}
type OpGroups []OpGroup

type OpMessage struct {
	OperationId string `json:"operationid,omitempty"`
	DefaultMsg  int    `json:"default_msg,omitempty"`
	// Possible values:
	// 0 - (default) use the data from the operation;
	// 1 - use the data from the action.
	MediaTypeId string `json:"mediatypeid,omitempty"`
	Message     string `json:"message,omitempty"`
	Subject     string `json:"subject,omitempty"`
}

type OpMessageGrp struct {
	OperationId string      `json:"operationid,omitempty"`
	UsrGrpId    HostGroupId `json:"groupid,omitempty"`
}
type OpMessageGrps []OpMessageGrp

type OpMessageUsr struct {
	OperationId string `json:"operationid,omitempty"`
	UserId      string `json:"userid,omitempty"`
}
type OpMessageUsrs []OpMessageUsr

type OpTemplate struct {
	OperationId string     `json:"operationid,omitempty"`
	TemplateId  TemplateId `json:"templateid,omitempty"`
}
type OpTemplates []OpTemplate

type OpInventory struct {
	OperationId   string `json:"operationid,omitempty"`
	InventoryMode string `json:"inventory_mode,omitempty"`
}

// Wrapper for action.get: https://www.zabbix.com/documentation/4.2/manual/appendix/api/action/get
func (api *API) ActionGet(params Params) (res Actions, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("action.get", params)
	if err != nil {
		return
	}

	mapstructure.Decode(response.Result.([]interface{}), &res)

	return
}

// Gets action by Id only if there is exactly 1 matching action.
func (api *API) ActionGetById(id string) (res *Action, err error) {
	actions, err := api.ActionGet(Params{"actionids": id})
	if err != nil {
		return
	}

	if len(actions) == 1 {
		res = &actions[0]
	} else {
		e := ExpectedOneResult(len(actions))
		err = &e
	}
	return
}

// Wrapper for action.create: https://www.zabbix.com/documentation/4.2/manual/appendix/api/action/create
func (api *API) ActionsCreate(actions Actions) (err error) {
	response, err := api.CallWithError("action.create", actions)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	actionsids := result["actionids"].([]interface{})
	for i, id := range actionsids {
		actionid := strconv.FormatFloat(id.(float64), 'f', 0, 64)
		actions[i].ActionId = actionid
	}
	return
}

// Wrapper for action.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/action/delete
// Cleans ActionId in all actions elements if call succeed.
func (api *API) ActionsDelete(actions Actions) (err error) {
	ids := make([]string, len(actions))
	for i, action := range actions {
		ids[i] = action.ActionId
	}

	err = api.ActionsDeleteByIds(ids)
	if err == nil {
		for i := range actions {
			actions[i].ActionId = ""
		}
	}
	return
}

// Wrapper for action.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/action/delete
func (api *API) ActionsDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("action.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	actionids := result["actionids"].([]interface{})
	if len(ids) != len(actionids) {
		err = &ExpectedMore{len(ids), len(actionids)}
	}
	return
}
