package expenses

type Storage interface {
	Insert(req ExpensesRequest) (int64, error)
}
