
/* Go module which converts IAM policy to HCL style IAM policy */

package main

import (
	"encoding/json"
	"fmt"
)

type IAMHCL struct {
	version string
	policy_id string
	statment []IAMHCLStatement
}

type IAMHCLStatement struct {
	sid string
	effect string
	actions []string
	not_actions []string
	resources []string
	not_resources []string
	principal []map[string]string
	not_principals []map[string]string
	condition []IAMHCLCondition
}

type IAMHCLCondition struct {
	test string
	variable string
	values []string
}

// IAMJson document representation
type IAMJson struct {
	Version string
	Id string
	Statement interface{}
}

type IAMJsonStatement struct {
	Sid string
	Effect string
	Principal interface{}
	NotPrincipal interface{}
	Action interface{}
	NotAction interface{}
	Resource interface{}
	NotResource interface{}
	Condition interface{}
}



func IAMJsonToHCL(obj interface{}) (interface{}, error) {
	return nil, nil
}

func IAMJsonToHCLb(b []byte) (interface{}, error) {
	var unm IAMJson

	err := json.Unmarshal(b, &unm)
	
	if err != nil {
		fmt.Println("error %s", err)
	}

	fmt.Printf("%+v", unm)

	return nil, nil
}