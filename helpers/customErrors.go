package helpers

type BoorrowError struct {
	Msg string
}

func (be *BoorrowError) Error() string {
	return be.Msg
}

func (be *BoorrowError) IsBorrowed() bool {
	return true
}

type IBorrowError interface {
	error
	IsBorrowed() bool
}

type ReturnBookError struct {
	Msg string
}

func (be *ReturnBookError) Error() string {
	return be.Msg
}

func (be *ReturnBookError) IsBorrowed() bool {
	return false
}

type IReturneBookError interface {
	error
	IsBorrowed() bool
}
