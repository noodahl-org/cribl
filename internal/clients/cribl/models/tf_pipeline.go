package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type Pipeline struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	TimeoutMS   types.Int64  `tfsdk:"timeout_ms"`
	Tags        types.List   `tfsdk:"tags"`
	Output      types.String `tfsdk:"output"`
}
