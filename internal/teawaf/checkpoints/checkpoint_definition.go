package checkpoints

// check point definition
type CheckpointDefinition struct {
	Name        string
	Description string
	Prefix      string
	HasParams   bool // has sub params
	Instance    CheckpointInterface
}
