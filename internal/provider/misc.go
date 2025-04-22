package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Filter struct {
	Name   types.String   `tfsdk:"name"`
	Values []types.String `tfsdk:"values"`
}

func matchesFilter(value string, filterValues []types.String) bool {
	for _, v := range filterValues {
		if value == v.ValueString() {
			return true
		}
	}
	return false
}
