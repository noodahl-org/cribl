package inputs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/noodahl-org/cribl/internal/clients/cribl"
	"github.com/noodahl-org/cribl/internal/clients/cribl/models"
)

type criblInputDatagenResource struct {
	client *cribl.Client
}

func NewCriblInputDatagenResource() resource.Resource {
	return &criblInputDatagenResource{}
}

func (r *criblInputDatagenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cribl.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *cribl.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *criblInputDatagenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_input_datagen"
}

func (r *criblInputDatagenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Input Id",
				Required:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "Disabled",
				Required:    true,
			},
			"environment": schema.StringAttribute{
				Description: "Environment",
				Optional:    true,
			},
			"pipeline": schema.StringAttribute{
				Description: "Pipeline",
				Optional:    true,
			},
			"pq_enabled": schema.BoolAttribute{
				Description: "PQ Enabled",
				Optional:    true,
			},
			"samples": schema.ListNestedAttribute{
				Description: "List of sample data configurations",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"events_per_sec": schema.Int64Attribute{
							Description: "Events Per Second",
							Required:    true,
						},
						"sample": schema.StringAttribute{
							Description: "Sample Type",
							Required:    true,
						},
					},
				},
			},
			// "streamtags": schema.ListAttribute{
			// 	Description: "Stream Tags",
			// 	Optional:    true,
			// 	ElementType: types.StringType,
			// },
			"type": schema.StringAttribute{
				Description: "Input Type",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description",
				Optional:    true,
			},
			"send_to_routes": schema.BoolAttribute{
				Description: "Send To Routes",
				Optional:    true,
			},
			// "metadata": schema.ListAttribute{
			// 	Description: "Metadata",
			// 	Optional:    true,
			// 	ElementType: types.StringType,
			// },
		},
	}
}

func (r *criblInputDatagenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.InputDatagen
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	inputBytes, err := json.Marshal(data.ToCriblInputDatagen())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to marshal input request to Cribl destination obj",
			err.Error(),
		)
	}
	inputRes, err := r.client.PostSystemInputs(ctx, cribl.Input{
		Union: json.RawMessage(inputBytes),
	})
	if err != nil || inputRes.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to create input datagen in Cribl",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *criblInputDatagenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.InputDatagen
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	inputBytes, err := json.Marshal(data.ToCriblInputDatagen())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to marshal input request to Cribl destination obj",
			err.Error(),
		)
		return
	}

	inputRes, err := r.client.PatchSystemInputsId(ctx, data.ID.ValueString(), cribl.Input{
		Union: json.RawMessage(inputBytes),
	})
	if err != nil || inputRes.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to update input datagen in Cribl",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *criblInputDatagenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan models.InputDatagen
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteSystemInputsId(ctx, plan.ID.ValueString(), r.client.RequestEditors...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unabled to delete datagen input from Cribl",
			err.Error(),
		)
	}

}

func (r *criblInputDatagenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.InputDatagen
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	inputRes, err := r.client.GetSystemInputsId(ctx, state.ID.ValueString(), r.client.RequestEditors...)
	if err != nil || inputRes.StatusCode == http.StatusInternalServerError {
		resp.Diagnostics.AddError(
			"Unable to fetch input from Cribl",
			err.Error(),
		)
	}
	if inputRes.StatusCode == http.StatusNotFound {
		return
	}
	tmp := struct {
		Items []cribl.InputDatagen `json:"items"`
	}{}
	if err := cribl.HandleResult(inputRes, err, &tmp); err != nil {
		resp.Diagnostics.AddError(
			"Unable to deseralize input response from Cribl",
			err.Error(),
		)
	}
	state.FromCriblDagen(tmp.Items[0])
	resp.Diagnostics.Append(req.State.Set(ctx, &state)...)
}
