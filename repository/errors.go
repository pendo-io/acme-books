package repository
type BorrowedError struct{}

func(error *BorrowedError) Error() string{
	return "book already borrowed"
}

type ReturnedError struct{}

func(error *ReturnedError) Error() string{
	return "book not loaned out"
}