package nintendo

// ErrIKSMExpired identifies the error that the IKSM is expired.
type ErrIKSMExpired struct {
	iksm string
}

func (err *ErrIKSMExpired) Error() string {
	return "expired iksm: " + err.iksm
}

// Is checks if an error is ErrIKSMExpired.
func (err *ErrIKSMExpired) Is(e error) bool {
	_, ok := e.(*ErrIKSMExpired)
	return ok
}
