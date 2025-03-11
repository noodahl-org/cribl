package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/noodahl-org/cribl/internal/clients/cribl"
)

type OutputS3 struct {
	ID                            types.String  `tfsdk:"id" validate:"required"`
	DefaultID                     types.String  `tfsdk:"default_id"`
	Type                          types.String  `tfsdk:"type"`
	Description                   types.String  `tfsdk:"description"`
	Environment                   types.String  `tfsdk:"environment"`
	StreamTags                    types.List    `tfsdk:"stream_tags"`
	SystemFields                  types.List    `tfsdk:"system_fields"`
	DeadletterEnabled             types.Bool    `tfsdk:"deadletter_enabled"`
	DeadletterPath                types.String  `tfsdk:"deadletter_path"`
	Bucket                        types.String  `tfsdk:"bucket"`
	Region                        types.String  `tfsdk:"region"`
	DestPath                      types.String  `tfsdk:"dest_path"`
	StagePath                     types.String  `tfsdk:"stage_path"`
	AddIdToStagePath              types.Bool    `tfsdk:"add_id_to_stage_path"`
	RemoveEmptyDirs               types.Bool    `tfsdk:"remove_empty_dirs"`
	EmptyDirCleanupSec            types.Float32 `tfsdk:"empty_dir_cleanup_sec"`
	BaseFileName                  types.String  `tfsdk:"base_file_name"`
	FileNameSuffix                types.String  `tfsdk:"file_name_suffix"`
	PartitionExpr                 types.String  `tfsdk:"partition_expr"`
	PartitioningFields            types.List    `tfsdk:"partitioning_fields"`
	Format                        types.String  `tfsdk:"format"`
	Compress                      types.String  `tfsdk:"compress"`
	CompressionLevel              types.String  `tfsdk:"compression_level"`
	MaxFileSizeMB                 types.Float32 `tfsdk:"max_file_size_mb"`
	MaxFileOpenTimeSec            types.Float32 `tfsdk:"max_file_open_time_sec"`
	MaxFileIdleTimeSec            types.Float32 `tfsdk:"max_file_idle_time_sec"`
	MaxOpenFiles                  types.Float32 `tfsdk:"max_open_files"`
	MaxConcurrentFileParts        types.Float32 `tfsdk:"max_concurrent_file_parts"`
	MaxClosingFilesToBackpressure types.Float32 `tfsdk:"max_closing_files_to_backpressure"`
	MaxRetryNum                   types.Float32 `tfsdk:"max_retry_num"`
	HeaderLine                    types.String  `tfsdk:"header_line"`
	OnBackpressure                types.String  `tfsdk:"on_backpressure"`
	OnDiskFullBackpressure        types.String  `tfsdk:"on_disk_full_backpressure"`
	WriteHighWaterMark            types.Float32 `tfsdk:"write_high_water_mark"`
	AwsAuthenticationMethod       types.String  `tfsdk:"aws_authentication_method"`
	AwsApiKey                     types.String  `tfsdk:"aws_api_key"`
	AwsSecretKey                  types.String  `tfsdk:"aws_secret_key"`
	AwsSecret                     types.String  `tfsdk:"aws_secret"`
	AssumeRoleArn                 types.String  `tfsdk:"assume_role_arn"`
	AssumeRoleExternalId          types.String  `tfsdk:"assume_role_external_id"`
	SignatureVersion              types.String  `tfsdk:"signature_version"`
	ReuseConnections              types.Bool    `tfsdk:"reuse_connections"`
	RejectUnauthorized            types.Bool    `tfsdk:"reject_unauthorized"`
	ObjectACL                     types.String  `tfsdk:"object_acl"`
	StorageClass                  types.String  `tfsdk:"storage_class"`
	ServerSideEncryption          types.String  `tfsdk:"server_side_encryption"`
	KmsKeyId                      types.String  `tfsdk:"kms_key_id"`
	Endpoint                      types.String  `tfsdk:"endpoint"`
	VerifyPermissions             types.Bool    `tfsdk:"verify_permissions"`
	DurationSeconds               types.Float32 `tfsdk:"duration_seconds"`
	EnableAssumeRole              types.Bool    `tfsdk:"enable_assume_role"`
	AutomaticSchema               types.Bool    `tfsdk:"automatic_schema"`
	EnablePageChecksum            types.Bool    `tfsdk:"enable_page_checksum"`
	EnableStatistics              types.Bool    `tfsdk:"enable_statistics"`
	EnableWritePageIndex          types.Bool    `tfsdk:"enable_write_page_index"`
	ParquetDataPageVersion        types.String  `tfsdk:"parquet_data_page_version"`
	ParquetPageSize               types.String  `tfsdk:"parquet_page_size"`
	ParquetRowGroupLength         types.Float32 `tfsdk:"parquet_row_group_length"`
	ParquetVersion                types.String  `tfsdk:"parquet_version"`
	Pipeline                      types.String  `tfsdk:"pipeline"`
	ShouldLogInvalidRows          types.Bool    `tfsdk:"should_log_invalid_rows"`
}

func (o *OutputS3) SetDefaults() {

}

func (o *OutputS3) FromCriblOutputS3(model cribl.OutputS3) {
	o.ID = types.StringPointerValue(model.Id)
	o.Type = types.StringPointerValue((*string)(model.Type))
	o.Description = types.StringPointerValue(model.Description)
	o.Environment = types.StringPointerValue(model.Environment)

	// if model.StreamTags != nil {
	// 	streamTags, _ := types.ListValueFrom(context.Background(), types.StringType, model.StreamTags)
	// 	o.StreamTags = streamTags
	// }

	if model.SystemFields != nil {
		systemFields, _ := types.ListValueFrom(context.Background(), types.StringType, model.SystemFields)
		o.SystemFields = systemFields
	}

	o.DeadletterEnabled = types.BoolPointerValue((*bool)(model.DeadletterEnabled))
	o.Bucket = types.StringValue(model.Bucket)
	o.DestPath = types.StringPointerValue(model.DestPath)
	o.StagePath = types.StringValue(model.StagePath)
	o.AddIdToStagePath = types.BoolPointerValue(model.AddIdToStagePath)
	o.RemoveEmptyDirs = types.BoolPointerValue((*bool)(model.RemoveEmptyDirs))
	o.EmptyDirCleanupSec = types.Float32PointerValue(model.EmptyDirCleanupSec)
	o.BaseFileName = types.StringPointerValue(model.BaseFileName)
	o.FileNameSuffix = types.StringPointerValue(model.FileNameSuffix)

	// if model.PartitioningFields != nil {
	// 	partitioningFields, _ := types.ListValueFrom(context.Background(), types.StringType, model.PartitioningFields)
	// 	o.PartitioningFields = partitioningFields
	// }

	o.Format = types.StringPointerValue((*string)(model.Format))
	o.Compress = types.StringPointerValue((*string)(model.Compress))
	o.CompressionLevel = types.StringPointerValue((*string)(model.CompressionLevel))
	o.MaxFileSizeMB = types.Float32PointerValue(model.MaxFileSizeMB)
	o.MaxFileOpenTimeSec = types.Float32PointerValue(model.MaxFileOpenTimeSec)
	o.MaxFileIdleTimeSec = types.Float32PointerValue(model.MaxFileIdleTimeSec)
	o.MaxOpenFiles = types.Float32PointerValue(model.MaxOpenFiles)
	o.MaxConcurrentFileParts = types.Float32PointerValue(model.MaxConcurrentFileParts)
	o.MaxClosingFilesToBackpressure = types.Float32PointerValue(model.MaxClosingFilesToBackpressure)
	o.HeaderLine = types.StringPointerValue(model.HeaderLine)
	o.OnBackpressure = types.StringPointerValue((*string)(model.OnBackpressure))
	o.OnDiskFullBackpressure = types.StringPointerValue((*string)(model.OnDiskFullBackpressure))
	o.WriteHighWaterMark = types.Float32PointerValue(model.WriteHighWaterMark)
	o.AwsAuthenticationMethod = types.StringPointerValue((*string)(model.AwsAuthenticationMethod))
	o.AwsSecret = types.StringPointerValue(model.AwsSecret)
	o.SignatureVersion = types.StringPointerValue((*string)(model.SignatureVersion))
	o.ReuseConnections = types.BoolPointerValue(model.ReuseConnections)
	o.RejectUnauthorized = types.BoolPointerValue(model.RejectUnauthorized)
	o.ObjectACL = types.StringPointerValue((*string)(model.ObjectACL))
	o.StorageClass = types.StringPointerValue((*string)(model.StorageClass))
	o.VerifyPermissions = types.BoolPointerValue(model.VerifyPermissions)
	o.DurationSeconds = types.Float32PointerValue(model.DurationSeconds)
	o.EnableAssumeRole = types.BoolPointerValue(model.EnableAssumeRole)
}

func (o *OutputS3) ToCriblOutputS3() cribl.OutputS3 {
	return cribl.OutputS3{
		AddIdToStagePath:              o.AddIdToStagePath.ValueBoolPointer(),
		AssumeRoleArn:                 o.AssumeRoleArn.ValueStringPointer(),
		AssumeRoleExternalId:          o.AssumeRoleExternalId.ValueStringPointer(),
		AutomaticSchema:               o.AutomaticSchema.ValueBoolPointer(),
		AwsApiKey:                     o.AwsApiKey.ValueStringPointer(),
		AwsAuthenticationMethod:       (*cribl.OutputS3AwsAuthenticationMethod)(o.AwsAuthenticationMethod.ValueStringPointer()),
		AwsSecret:                     o.AwsSecret.ValueStringPointer(),
		AwsSecretKey:                  o.AwsSecretKey.ValueStringPointer(),
		BaseFileName:                  o.BaseFileName.ValueStringPointer(),
		Bucket:                        o.Bucket.ValueString(),
		Compress:                      (*cribl.OutputS3Compress)(o.Compress.ValueStringPointer()),
		CompressionLevel:              (*cribl.OutputS3CompressionLevel)(o.CompressionLevel.ValueStringPointer()),
		DeadletterEnabled:             (*cribl.OutputS3DeadletterEnabled)(o.DeadletterEnabled.ValueBoolPointer()),
		DeadletterPath:                o.DeadletterPath.ValueStringPointer(),
		Description:                   o.Description.ValueStringPointer(),
		DestPath:                      o.DestPath.ValueStringPointer(),
		DurationSeconds:               o.DurationSeconds.ValueFloat32Pointer(),
		EmptyDirCleanupSec:            o.EmptyDirCleanupSec.ValueFloat32Pointer(),
		EnableAssumeRole:              o.EnableAssumeRole.ValueBoolPointer(),
		EnablePageChecksum:            o.EnablePageChecksum.ValueBoolPointer(),
		EnableStatistics:              o.EnableStatistics.ValueBoolPointer(),
		EnableWritePageIndex:          o.EnableWritePageIndex.ValueBoolPointer(),
		Endpoint:                      o.Endpoint.ValueStringPointer(),
		Environment:                   o.Environment.ValueStringPointer(),
		FileNameSuffix:                o.FileNameSuffix.ValueStringPointer(),
		Format:                        (*cribl.OutputS3Format)(o.Format.ValueStringPointer()),
		HeaderLine:                    o.HeaderLine.ValueStringPointer(),
		Id:                            o.ID.ValueStringPointer(),
		KmsKeyId:                      o.KmsKeyId.ValueStringPointer(),
		MaxClosingFilesToBackpressure: o.MaxClosingFilesToBackpressure.ValueFloat32Pointer(),
		MaxConcurrentFileParts:        o.MaxConcurrentFileParts.ValueFloat32Pointer(),
		MaxFileIdleTimeSec:            o.MaxFileIdleTimeSec.ValueFloat32Pointer(),
		MaxFileOpenTimeSec:            o.MaxFileOpenTimeSec.ValueFloat32Pointer(),
		MaxFileSizeMB:                 o.MaxFileSizeMB.ValueFloat32Pointer(),
		MaxOpenFiles:                  o.MaxOpenFiles.ValueFloat32Pointer(),
		MaxRetryNum:                   o.MaxRetryNum.ValueFloat32Pointer(),
		ObjectACL:                     (*cribl.OutputS3ObjectACL)(o.ObjectACL.ValueStringPointer()),
		OnBackpressure:                (*cribl.OutputS3OnBackpressure)(o.OnBackpressure.ValueStringPointer()),
		OnDiskFullBackpressure:        (*cribl.OutputS3OnDiskFullBackpressure)(o.OnDiskFullBackpressure.ValueStringPointer()),
		ParquetDataPageVersion:        (*cribl.OutputS3ParquetDataPageVersion)(o.ParquetDataPageVersion.ValueStringPointer()),
		ParquetPageSize:               o.ParquetPageSize.ValueStringPointer(),
		ParquetRowGroupLength:         o.ParquetRowGroupLength.ValueFloat32Pointer(),
		ParquetVersion:                (*cribl.OutputS3ParquetVersion)(o.ParquetVersion.ValueStringPointer()),
		PartitionExpr:                 o.PartitionExpr.ValueStringPointer(),
		Pipeline:                      o.Pipeline.ValueStringPointer(),
		Region:                        o.Region.ValueStringPointer(),
		RejectUnauthorized:            o.RejectUnauthorized.ValueBoolPointer(),
		RemoveEmptyDirs:               (*cribl.OutputS3RemoveEmptyDirs)(o.RemoveEmptyDirs.ValueBoolPointer()),
		ReuseConnections:              o.ReuseConnections.ValueBoolPointer(),
		ServerSideEncryption:          (*cribl.OutputS3ServerSideEncryption)(o.ServerSideEncryption.ValueStringPointer()),
		ShouldLogInvalidRows:          o.ShouldLogInvalidRows.ValueBoolPointer(),
		SignatureVersion:              (*cribl.OutputS3SignatureVersion)(o.SignatureVersion.ValueStringPointer()),
		StagePath:                     o.StagePath.ValueString(),
		StorageClass:                  (*cribl.OutputS3StorageClass)(o.StorageClass.ValueStringPointer()),
		Type:                          (*cribl.OutputS3Type)(o.Type.ValueStringPointer()),
		VerifyPermissions:             o.VerifyPermissions.ValueBoolPointer(),
		WriteHighWaterMark:            o.WriteHighWaterMark.ValueFloat32Pointer(),
	}
}
