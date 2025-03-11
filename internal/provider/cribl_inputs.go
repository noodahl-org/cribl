package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/noodahl-org/cribl/internal/clients/cribl"
)

type criblInputResource struct {
	client *cribl.Client
}

func NewCriblInputResource() resource.Resource {
	return &criblInputResource{}
}

func (r *criblInputResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.

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

func (r *criblInputResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inputs"
}

func (r *criblInputResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Input Id",
				Required:    true,
			},
		},
	}
}

func (r *criblInputResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

}

func (r *criblInputResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *criblInputResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
}

func (r *criblInputResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {

}
