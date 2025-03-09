package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/noodahl-org/cribl/internal/clients/cribl"
	"github.com/noodahl-org/cribl/internal/clients/cribl/models"
)

type criblDataSourceModel struct {
	Build models.Build `tfsdk:"build"`
}

type criblDataSource struct {
	client *cribl.Client
}

func NewCriblDataSource() datasource.DataSource {
	return &criblDataSource{}
}

func (d *criblDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

func (d *criblDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Cribl information from system/info",
		Attributes: map[string]schema.Attribute{
			"build": schema.SingleNestedAttribute{
				Description: "Build information for the Cribl system",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"version": schema.StringAttribute{
						Description: "Cribl build version",
						Computed:    true,
					},
					"hostname": schema.StringAttribute{
						Description: "Cribl build hostname",
						Computed:    true,
					},
					"branch": schema.StringAttribute{
						Description: "Cribl branch",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *criblDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state criblDataSourceModel

	info := struct {
		Items []cribl.SystemInfo `json:"items"`
	}{}
	infoResp, err := d.client.GetSystemInfo(ctx, d.client.RequestEditors...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch Cribl System Info",
			err.Error(),
		)
		return
	}
	if err := cribl.HandleResult(infoResp, err, &info); err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Cribl System Info response",
			err.Error(),
		)
		return
	}
	state.Build = models.Build{
		Hostname: info.Items[0].Hostname,
		Version:  fmt.Sprintf("%s", info.Items[0].BUILD["VERSION"]),
		Branch:   fmt.Sprintf("%s", info.Items[0].BUILD["BRANCH"]),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *criblDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*cribl.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *cribl.Client, go: %T.", req.ProviderData),
		)
		return
	}
	d.client = client
}
