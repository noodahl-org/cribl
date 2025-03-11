package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/noodahl-org/cribl/internal/clients/cribl"
	"github.com/noodahl-org/cribl/internal/clients/cribl/models"
)

type criblOutputS3Resource struct {
	client *cribl.Client
}

func NewCriblOutputS3Resource() resource.Resource {
	return &criblOutputS3Resource{}
}

func (r *criblOutputS3Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *criblOutputS3Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_output_s3"
}

func (r *criblOutputS3Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique ID for this output",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Output type",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of this output",
				Optional:    true,
			},
			"environment": schema.StringAttribute{
				Description: "Optionally, enable this config only on a specified Git branch",
				Optional:    true,
			},
			"stream_tags": schema.ListAttribute{
				Description: "Tags for filtering and grouping in Cribl",
				Optional:    true,
				ElementType: types.StringType,
			},
			"system_fields": schema.ListAttribute{
				Description: "Fields to automatically add to events, such as cribl_pipe. Supports wildcards.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"deadletter_enabled": schema.BoolAttribute{
				Description: "If a file fails to move to its final destination after the maximum number of retries, dead-letter it to prevent further errors",
				Optional:    true,
			},
			"deadletter_path": schema.StringAttribute{
				Description: "Storage location for files that fail to reach their final destination after maximum retries are exceeded",
				Optional:    true,
			},
			"bucket": schema.StringAttribute{
				Description: "Name of the destination S3 bucket",
				Optional:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region where the S3 bucket is located",
				Optional:    true,
			},
			"dest_path": schema.StringAttribute{
				Description: "Prefix to append to files before uploading",
				Optional:    true,
			},
			"stage_path": schema.StringAttribute{
				Description: "Filesystem location in which to buffer files, before compressing and moving to final destination",
				Optional:    true,
			},
			"add_id_to_stage_path": schema.BoolAttribute{
				Description: "Append output's ID to staging location",
				Optional:    true,
			},
			"remove_empty_dirs": schema.BoolAttribute{
				Description: "Remove empty staging directories after moving files",
				Optional:    true,
			},
			"empty_dir_cleanup_sec": schema.Float32Attribute{
				Description: "How frequently, in seconds, to clean up empty directories when 'Remove empty staging dirs' is enabled",
				Optional:    true,
			},
			"base_file_name": schema.StringAttribute{
				Description: "JavaScript expression to define the output filename prefix",
				Optional:    true,
			},
			"file_name_suffix": schema.StringAttribute{
				Description: "JavaScript expression to define the output filename suffix",
				Optional:    true,
			},
			"partition_expr": schema.StringAttribute{
				Description: "JavaScript expression defining how files are partitioned and organized",
				Optional:    true,
			},
			"partitioning_fields": schema.ListAttribute{
				Description: "Fields to use for partitioning",
				Optional:    true,
				ElementType: types.StringType,
			},
			"format": schema.StringAttribute{
				Description: "Format of the output data",
				Optional:    true,
			},
			"compress": schema.StringAttribute{
				Description: "Choose data compression format to apply before moving files to final destination",
				Optional:    true,
			},
			"compression_level": schema.StringAttribute{
				Description: "Compression level to apply before moving files to final destination",
				Optional:    true,
			},
			"max_file_size_mb": schema.Float32Attribute{
				Description: "Maximum uncompressed output file size. Files of this size will be closed and moved to final output location",
				Optional:    true,
			},
			"max_file_open_time_sec": schema.Float32Attribute{
				Description: "Maximum amount of time to write to a file. Files open for longer than this will be closed and moved to final output location",
				Optional:    true,
			},
			"max_file_idle_time_sec": schema.Float32Attribute{
				Description: "Maximum amount of time to keep inactive files open. Files open for longer than this will be closed and moved to final output location",
				Optional:    true,
			},
			"max_open_files": schema.Float32Attribute{
				Description: "Maximum number of files to keep open concurrently. When exceeded, Cribl will close the oldest open files and move them to the final output location",
				Optional:    true,
			},
			"max_concurrent_file_parts": schema.Float32Attribute{
				Description: "Maximum number of parts to upload in parallel per file. Minimum part size is 5MB",
				Optional:    true,
			},
			"max_closing_files_to_backpressure": schema.Float32Attribute{
				Description: "Maximum number of files that can be waiting for upload before backpressure is applied",
				Optional:    true,
			},
			"max_retry_num": schema.Float32Attribute{
				Description: "The maximum number of times a file will attempt to move to its final destination before being dead-lettered",
				Optional:    true,
			},
			"header_line": schema.StringAttribute{
				Description: "If set, this line will be written to the beginning of each output file",
				Optional:    true,
			},
			"on_backpressure": schema.StringAttribute{
				Description: "Whether to block or drop events when all receivers are exerting backpressure",
				Optional:    true,
			},
			"on_disk_full_backpressure": schema.StringAttribute{
				Description: "Whether to block or drop events when disk space is below the global 'Min free disk space' limit",
				Optional:    true,
			},
			"write_high_water_mark": schema.Float32Attribute{
				Description: "Buffer size used to write to a file",
				Optional:    true,
			},
			"aws_authentication_method": schema.StringAttribute{
				Description: "AWS authentication method. Choose Auto to use IAM roles",
				Optional:    true,
			},
			"aws_api_key": schema.StringAttribute{
				Description: "Access key. This value can be a constant or a JavaScript expression",
				Optional:    true,
				Sensitive:   true,
			},
			"aws_secret_key": schema.StringAttribute{
				Description: "Secret key. This value can be a constant or a JavaScript expression",
				Optional:    true,
				Sensitive:   true,
			},
			"aws_secret": schema.StringAttribute{
				Description: "Select or create a stored secret that references your access key and secret key",
				Optional:    true,
				Sensitive:   true,
			},
			"assume_role_arn": schema.StringAttribute{
				Description: "Amazon Resource Name (ARN) of the role to assume",
				Optional:    true,
			},
			"assume_role_external_id": schema.StringAttribute{
				Description: "External ID to use when assuming role",
				Optional:    true,
			},
			"signature_version": schema.StringAttribute{
				Description: "Signature version to use for signing S3 requests",
				Optional:    true,
			},
			"reuse_connections": schema.BoolAttribute{
				Description: "Reuse connections between requests, which can improve performance",
				Optional:    true,
			},
			"reject_unauthorized": schema.BoolAttribute{
				Description: "Reject certificates that cannot be verified against a valid CA, such as self-signed certificates",
				Optional:    true,
			},
			"object_acl": schema.StringAttribute{
				Description: "Object ACL to assign to uploaded objects",
				Optional:    true,
			},
			"storage_class": schema.StringAttribute{
				Description: "Storage class to select for uploaded objects",
				Optional:    true,
			},
			"server_side_encryption": schema.StringAttribute{
				Description: "Server-side encryption for uploaded objects",
				Optional:    true,
			},
			"kms_key_id": schema.StringAttribute{
				Description: "ID or ARN of the KMS customer-managed key to use for encryption",
				Optional:    true,
			},
			"endpoint": schema.StringAttribute{
				Description: "S3 service endpoint. If empty, defaults to AWS' Region-specific endpoint. Otherwise, it must point to S3-compatible endpoint",
				Optional:    true,
			},
			"verify_permissions": schema.BoolAttribute{
				Description: "Disable if you can access files within the bucket but not the bucket itself",
				Optional:    true,
			},
			"duration_seconds": schema.Float32Attribute{
				Description: "Duration of the assumed role's session, in seconds. Minimum is 900 (15 minutes), default is 3600 (1 hour), and maximum is 43200 (12 hours)",
				Optional:    true,
			},
			"enable_assume_role": schema.BoolAttribute{
				Description: "Use Assume Role credentials to access S3",
				Optional:    true,
			},
			"automatic_schema": schema.BoolAttribute{
				Description: "Automatically calculate the schema based on the events of each Parquet file generated",
				Optional:    true,
			},
			"enable_page_checksum": schema.BoolAttribute{
				Description: "Parquet tools can use the checksum of a Parquet page to verify data integrity",
				Optional:    true,
			},
			"enable_statistics": schema.BoolAttribute{
				Description: "Statistics profile an entire file in terms of minimum/maximum values within data, numbers of nulls, etc. You can use Parquet tools to view statistics",
				Optional:    true,
			},
			"enable_write_page_index": schema.BoolAttribute{
				Description: "One page index contains statistics for one data page. Parquet readers use statistics to enable page skipping",
				Optional:    true,
			},
			"parquet_data_page_version": schema.StringAttribute{
				Description: "Serialization format of data pages. Note that some reader implementations use Data page V2's attributes to work more efficiently, while others ignore it",
				Optional:    true,
			},
			"parquet_page_size": schema.StringAttribute{
				Description: "Target memory size for page segments, such as 1MB or 128MB. Generally, lower values improve reading speed, while higher values improve compression",
				Optional:    true,
			},
			"parquet_row_group_length": schema.Float32Attribute{
				Description: "The number of rows that every group will contain. The final group can contain a smaller number of rows",
				Optional:    true,
			},
			"parquet_version": schema.StringAttribute{
				Description: "Determines which data types are supported and how they are represented",
				Optional:    true,
			},
			"pipeline": schema.StringAttribute{
				Description: "Pipeline to process data before sending out to this output",
				Optional:    true,
			},
			"should_log_invalid_rows": schema.BoolAttribute{
				Description: "Log up to 3 rows that Cribl skips due to data mismatch",
				Optional:    true,
			},
			"default_id": schema.StringAttribute{
				Description: "Default ID",
				Required:    true,
			},
			// "key_value_metadata": schema.ListNestedAttribute{
			// 	Description: "The metadata of files the Destination writes will include the properties you add here as key-value pairs",
			// 	Optional:    true,
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"key": schema.StringAttribute{
			// 				Description: "Metadata key",
			// 				Required:    true,
			// 			},
			// 			"value": schema.StringAttribute{
			// 				Description: "Metadata value",
			// 				Required:    true,
			// 			},
			// 		},
			// 	},
			// },
		},
	}
}
func (r *criblOutputS3Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.OutputS3
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	//convert data to output s3 request
	outputBytes, err := json.Marshal(data.ToCriblOutputS3())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to marshal request to Cribl",
			err.Error(),
		)
		return
	}
	outputRes, err := r.client.PostSystemOutputs(ctx, cribl.Output{
		Union: json.RawMessage(outputBytes),
	})

	if err != nil || outputRes.StatusCode != http.StatusResetContent {
		resp.Diagnostics.AddError(
			"Unable to create output in Cribl",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// todo finish
func (r *criblOutputS3Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.OutputS3

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
}

func (r *criblOutputS3Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.OutputS3
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.DeleteSystemOutputsId(ctx, data.ID.ValueString(), r.client.RequestEditors...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete output from Cribl",
			err.Error(),
		)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *criblOutputS3Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.OutputS3
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tmp := cribl.OutputS3{}
	outputRes, err := r.client.GetSystemOutputsId(ctx, data.ID.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch output from Cribl",
			err.Error(),
		)
	}
	if err := cribl.HandleResult(outputRes, err, &tmp); err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse outputs response from Cribl",
			err.Error(),
		)
		return
	}
	data.FromCriblOutputS3(tmp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *criblOutputS3Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
