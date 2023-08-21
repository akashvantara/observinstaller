package conf

type DownloadOptions struct {
	InstallationType string // -minimal, -full
}

type RunOptions struct {
	RunType string // same as installation type
}

type KillOptions struct {
	KillType string // same as installation type
}

type RemoveOptions struct {
	RemoveType string // -minimal, -full
}
