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
			"id": schema.StringAttribute{
				Description: "Pipeline Id",
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
	}
}

func (c *criblPipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.Pipeline
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var tags []string
	resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r, err := c.client.PostPipelines(ctx, cribl.Pipeline{
		Id: plan.ID.ValueString(),
		Conf: cribl.PipelineConf{
			AsyncFuncTimeout: lo.ToPtr(int(plan.TimeoutMS.ValueInt64())),
			Description:      plan.Description.ValueStringPointer(),
			//todo: fix streamtags - threw a 500 while calling into the api
			//Streamtags:       lo.ToPtr(tags),
			Output: plan.Output.ValueStringPointer(),
		},
	})
	if err != nil || r.StatusCode > 400 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating pipeline. Status %v", r.StatusCode),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *criblPipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.Pipeline
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *criblPipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan models.Pipeline
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := c.client.DeletePipelinesId(ctx, plan.ID.ValueString(), c.client.RequestEditors...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unabled to delete Cribl pipline",
			err.Error(),
		)
		return
	}
	
}

func (c *criblPipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.Pipeline
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

	pipelines := struct {
		Items []cribl.Pipeline `json:"items"`
	}{}
	if err := cribl.HandleResult(pipelineResp, err, &pipelines); err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch pipelines from Cribl",
			err.Error(),
		)
		return
	}
	for _, pipeline := range pipelines.Items {
		if pipeline.Id != state.ID.ValueString() {
			continue
		}
		state.Description = types.StringValue(lo.FromPtr(pipeline.Conf.Description))
		if pipeline.Conf.AsyncFuncTimeout != nil {
			state.TimeoutMS = types.Int64Value(int64(*pipeline.Conf.AsyncFuncTimeout))
		}
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
