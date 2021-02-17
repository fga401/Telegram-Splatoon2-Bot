package nintendo

type ErrIKSMExpired struct {
	iksm string
}

func (err *ErrIKSMExpired) Error() string {
	return "expired iksm: " + err.iksm
}

