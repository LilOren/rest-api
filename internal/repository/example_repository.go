package repository

type (
	ExampleRepository interface{}
	exampleRepository struct{}
)

func NewExampleRepository() ExampleRepository {
	return &exampleRepository{}
}
