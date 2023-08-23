package conf

type DownloadOptions struct {
	InstallationType string // -minimal, -full
}

type RunOptions struct {
	RunType string // same as installation type
}

type KillOptions struct {
	KillType string // same as installation type
	Restart  bool   // restart apps instead of killing
}

type RemoveOptions struct {
	RemoveType string // -d downloads, -i installs, -a all
}

type OtelOptions struct {
	Build bool
	List  bool
}
