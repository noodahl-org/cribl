package provider

import (
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/noodahl-org/cribl/internal/clients/cribl"
)

var (
	_ provider.Provider = &criblProvider{}
)

type criblProviderModel struct {
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	Token          types.String `tfsdk:"token"`
	WorkspaceID    types.String `tfsdk:"workspace_id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	BaseURL        types.String `tfsdk:"base_url"`
}

type criblProvider struct {
	version string
	client  *cribl.Client
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &criblProvider{
			version: version,
		}
	}
}

func (p *criblProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cribl"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *criblProvider) Schema(_ context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional: true,
			},
			"base_url": schema.StringAttribute{
				Optional: true,
			},
			"token": schema.StringAttribute{
				Optional: true,
			},
			"workspace_id": schema.StringAttribute{
				Optional: true,
			},
			"organization_id": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *criblProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config criblProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//todo: cribl cloud support. https://docs.cribl.io/edge/api-tutorials#criblcloud
	client, err := cribl.NewClient(config.BaseURL.ValueString() + "/api/v1")
	if err != nil {
		diags.AddError(
			"Unable to instantiate cribl client.",
			err.Error(),
		)
		return
	}

	if config.BaseURL.IsUnknown() || config.BaseURL.IsNull() {
		if val, ok := os.LookupEnv("CRIBL_URL"); ok {
			config.BaseURL = types.StringPointerValue(&val)
		} else {
			resp.Diagnostics.AddAttributeError(
				path.Root("base_url"),
				"Unknown URL",
				"The provider cannot reach cribl without an endpoint defined",
			)
		}
	}
	if !config.Username.IsUnknown() && !config.Password.IsUnknown() {
		tokenData := cribl.AuthToken{}
		resp, err := client.PostAuthLogin(ctx, cribl.LoginInfo{
			Username: config.Username.ValueString(),
			Password: config.Password.ValueString(),
		})
		if err := cribl.HandleResult(resp, err, &tokenData); err != nil {
			diags.AddError(
				"Unable to fetch auth token for cribl client.",
				err.Error(),
			)
			return
		}
		// include the token in each request from the client
		client.RequestEditors = append(client.RequestEditors, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+tokenData.Token)
			return nil
		})
	}

	p.client = client
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *criblProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCriblDataSource,
	}
}

func (p *criblProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCriblPipelineResource,
		NewCriblOutputS3Resource,
	}
}
