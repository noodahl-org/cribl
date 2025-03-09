package models

type Build struct {
	Hostname string `tfsdk:"hostname"`
	Version  string `tfsdk:"version"`
	Branch   string `tfsdk:"branch"`
}
