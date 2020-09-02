/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyGatewayPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyGatewayPolicyRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"domain":       getDataSourceDomainNameSchema(),
			"category": {
				Type:         schema.TypeString,
				Description:  "Category",
				ValidateFunc: validation.StringInSlice(gatewayPolicyCategoryValues, false),
				Optional:     true,
				Computed:     true,
			},
		},
	}
}

// Local Manager Only
func listGatewayPolicies(domain string, connector *client.RestConnector) ([]model.GatewayPolicy, error) {
	client := domains.NewDefaultGatewayPoliciesClient(connector)

	var results []model.GatewayPolicy
	boolFalse := false
	var cursor *string
	total := 0

	for {
		policies, err := client.List(domain, cursor, nil, &boolFalse, nil, nil, &boolFalse, nil)
		if err != nil {
			return results, err
		}
		results = append(results, policies.Results...)
		if total == 0 && policies.ResultCount != nil {
			// first response
			total = int(*policies.ResultCount)
		}

		cursor = policies.Cursor
		if len(results) >= total {
			return results, nil
		}
	}
}

func dataSourceNsxtPolicyGatewayPolicyRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	category := d.Get("category").(string)
	domain := d.Get("domain").(string)
	if isPolicyGlobalManager(m) {
		query := make(map[string]string)
		query["parent_path"] = "*/" + domain
		if category != "" {
			query["category"] = category
		}
		obj, err := policyDataSourceResourceReadWithValidation(d, connector, "GatewayPolicy", query, false)
		if err != nil {
			return err
		}

		converter := bindings.NewTypeConverter()
		converter.SetMode(bindings.REST)
		dataValue, errors := converter.ConvertToGolang(obj, gm_model.GatewayPolicyBindingType())
		if len(errors) > 0 {
			return errors[0]
		}

		policy := dataValue.(gm_model.GatewayPolicy)
		d.Set("category", policy.Category)
		return nil
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)

	var obj model.GatewayPolicy
	if objID != "" {
		// Get by id
		client := domains.NewDefaultGatewayPoliciesClient(connector)
		objGet, err := client.Get(domain, objID)
		if isNotFoundError(err) {
			return fmt.Errorf("Gateway Policy with ID %s was not found", objID)
		}

		if err != nil {
			return fmt.Errorf("Error while reading Gateway Policy %s: %v", objID, err)
		}
		obj = objGet
	} else if objName == "" && category == "" {
		return fmt.Errorf("Gateway Policy id, display name or category must be specified")
	} else {
		objList, err := listGatewayPolicies(domain, connector)
		if err != nil {
			return fmt.Errorf("Error while reading Gateway Policies: %v", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.GatewayPolicy
		var prefixMatch []model.GatewayPolicy
		for _, objInList := range objList {
			if category != "" && category != *objInList.Category {
				continue
			}
			if objName != "" {
				if strings.HasPrefix(*objInList.DisplayName, objName) {
					prefixMatch = append(prefixMatch, objInList)
				}
				if *objInList.DisplayName == objName {
					perfectMatch = append(perfectMatch, objInList)
				}
			} else {
				prefixMatch = append(prefixMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple Gateway Policies with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple Gateway Policies with name starting with '%s' and category '%s'", objName, category)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("Gateway Policy with name '%s' and category '%s' was not found", objName, category)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("category", obj.Category)

	return nil
}