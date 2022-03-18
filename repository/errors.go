package repository
type BorrowedError struct{}

func(error *BorrowedError) Error() string{
	return "book already borrowed"
}