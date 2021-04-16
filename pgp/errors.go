package pgp

type keyIncorrectError int

func (ki keyIncorrectError) Error() string {
	return "vaultz: incorrect key"
}

var ErrKeyIncorrect error = keyIncorrectError(0)
