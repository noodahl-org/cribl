package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/noodahl-org/cribl/internal/clients/cribl"
)

type Sample struct {
	EventsPerSec types.Int64  `tfsdk:"events_per_sec"`
	Sample       types.String `tfsdk:"sample"`
}

type Connection struct {
	Output string
}

type InputDatagen struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Environment types.String `tfsdk:"environment"`
	Samples     []Sample     `tfsdk:"samples"`
	//StreamTags   types.List   `tfsdk:"stream_tags"`
	Type         types.String `tfsdk:"type"`
	Disabled     types.Bool   `tfsdk:"disabled"`
	PQEnabled    types.Bool   `tfsdk:"pq_enabled"`
	SendToRoutes types.Bool   `tfsdk:"send_to_routes"`
	//Connections  []Connection `tfsdk:"connections"`
	Pipeline types.String `tfsdk:"pipeline"`
}

func (i *InputDatagen) ToCriblInputDatagen() cribl.InputDatagen {

	out := cribl.InputDatagen{
		Id:           i.ID.ValueStringPointer(),
		Description:  i.Description.ValueStringPointer(),
		Type:         cribl.InputDatagenType(*i.Type.ValueStringPointer()),
		Environment:  i.Environment.ValueStringPointer(),
		Disabled:     i.Disabled.ValueBoolPointer(),
		PqEnabled:    (*cribl.InputDatagenPqEnabled)(i.PQEnabled.ValueBoolPointer()),
		SendToRoutes: (*cribl.InputDatagenSendToRoutes)(i.SendToRoutes.ValueBoolPointer()),
		Pipeline:     i.Pipeline.ValueStringPointer(),
	}
	for _, sample := range i.Samples {
		out.Samples = append(out.Samples, struct {
			EventsPerSec float32 "json:\"eventsPerSec\""
			Sample       string  "json:\"sample\""
		}{
			EventsPerSec: float32(sample.EventsPerSec.ValueInt64()),
			Sample:       sample.Sample.ValueString(),
		})
	}

	return out
}

func (i *InputDatagen) FromCriblDagen(model cribl.InputDatagen) {
	i.ID = types.StringPointerValue(model.Id)
	i.Description = types.StringPointerValue(model.Description)
	i.Environment = types.StringPointerValue(model.Environment)
	i.Type = types.StringValue(string(model.Type))
	i.Disabled = types.BoolPointerValue(model.Disabled)
	i.PQEnabled = types.BoolPointerValue((*bool)(model.PqEnabled))
	i.SendToRoutes = types.BoolPointerValue((*bool)(model.SendToRoutes))
	i.Pipeline = types.StringPointerValue(model.Pipeline)
	for _, sample := range model.Samples {
		i.Samples = append(i.Samples, Sample{
			EventsPerSec: types.Int64Value(int64(sample.EventsPerSec)),
			Sample:       types.StringValue(sample.Sample),
		})
	}
}
