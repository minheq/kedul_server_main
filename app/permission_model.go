package app

// Operation defines which individual API is being executed on the application level. It maps to a permissions.
// The list of operations are stored in the database
type Operation struct {
	Name string
}

// Permission defines set of allowed operations. It is user facing, i.e. users will be able
// to modify permissions settings
// The list of permissions are stored in the database
type Permission struct {
	Name       string
	Operations []Operation
}
