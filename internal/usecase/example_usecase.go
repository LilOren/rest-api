package usecase

type (
	ExampleUsecase interface{}
	exampleUsecase struct{}
)

func NewExampleRepository() ExampleUsecase {
	return &exampleUsecase{}
}
