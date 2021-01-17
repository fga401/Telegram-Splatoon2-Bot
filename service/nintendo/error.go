package nintendo

type ErrIKSMExpired struct {
	iksm string
}

func (err *ErrIKSMExpired) Error() string {
	return "expired iksm: " + err.iksm
}

func (err *ErrIKSMExpired) Is(e error) bool {
	_, ok := e.(*ErrIKSMExpired)
	return ok
}