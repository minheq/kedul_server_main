package app

// Operation defines which individual API is being executed on the application level. It maps to a permissions.
type Operation struct {
	Name string
}

// Permission defines set of allowed operations. It is user facing, i.e. users will be able
// to modify permissions settings
type Permission struct {
	ID         string
	Name       string
	Operations []Operation
}
