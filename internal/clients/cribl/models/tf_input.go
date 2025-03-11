package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type Sample struct {
	EventsPerSec types.Int64  `tfsdk:"events_per_sec"`
	Sample       types.String `tfsdk:"sample"`
}

type Input struct {
	ID           types.String `tfsdk:"id"`
	Description  types.String `tfsdk:"description"`
	Samples      []Sample     `tfsdk:"samples"`
	StreamTags   types.List   `tfsdk:"stream_tags"`
	Type         types.String `tfsdk:"type"`
	Disabled     types.Bool   `tfsdk:"disabled"`
	PQEnabled    types.Bool   `tfsdk:"pq_enabled"`
	SendToRoutes types.Bool   `tfsdk:"send_to_routes"`
}
