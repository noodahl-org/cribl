package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/noodahl-org/cribl/internal/clients/cribl"
	"github.com/noodahl-org/cribl/internal/clients/cribl/models"
	"github.com/samber/lo"
)

type criblPipelineResourceModel struct {
	Pipelines []models.Pipeline `tfsdk:"pipelines"`
}

type criblPipelineResource struct {
	client *cribl.Client
}

func NewCriblPipelineResource() resource.Resource {
	return &criblPipelineResource{}
}

func (r *criblPipelineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pipeline"
}

func (r *criblPipelineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"pipelines": schema.ListNestedAttribute{
				Description: "List of pipelines",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Pipeline ID",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "Pipeline description",
							Optional:    true,
						},
						"timeout_ms": schema.Int64Attribute{
							Description: "Pipeline timeout in ms",
							Required:    true,
						},
						"output": schema.StringAttribute{
							Description: "Pipeline output",
							Required:    true,
						},
						"tags": schema.ListAttribute{
							Description: "Pipeline tags",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (c *criblPipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan criblPipelineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, pipeline := range plan.Pipelines {
		// Extract string values from the tags list
		var tags []string
		resp.Diagnostics.Append(pipeline.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if _, err := c.client.PostPipelines(ctx, cribl.Pipeline{
			Id: pipeline.ID.ValueString(),
			Conf: cribl.PipelineConf{
				AsyncFuncTimeout: lo.ToPtr(int(pipeline.TimeoutMS.ValueInt64())),
				Description:      pipeline.Description.ValueStringPointer(),
				Streamtags:       lo.ToPtr(tags),
				Output:           pipeline.Output.ValueStringPointer(),
			},
		}); err != nil {
			resp.Diagnostics.AddError(
				"Error creating pipeline",
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *criblPipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan criblPipelineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *criblPipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan criblPipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *criblPipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state criblPipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	//fetch pipelines
	pipelineResp, err := c.client.GetPipelines(ctx, c.client.RequestEditors...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch pipelines from Cribl",
			err.Error(),
		)
		return
	}

	pipelines := []cribl.Pipeline{}
	if err := cribl.HandleResult(pipelineResp, err, &pipelines); err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch pipelines from Cribl",
			err.Error(),
		)
		return
	}
	state.Pipelines = []models.Pipeline{}
	for _, pipeline := range pipelines {
		tmp := models.Pipeline{
			ID:        types.StringValue(pipeline.Id),
			TimeoutMS: types.Int64Value(int64(*pipeline.Conf.AsyncFuncTimeout)),
		}
		state.Pipelines = append(state.Pipelines, tmp)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *criblPipelineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
