package data

//go:generate go run gen.go

// TuberYaml is a generic representation of a default yaml for new tuber apps
type TuberYaml struct {
	Filename string
	Contents string
}
