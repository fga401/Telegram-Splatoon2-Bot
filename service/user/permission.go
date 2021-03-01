package user

func (svc *serviceImpl) GetPermission(uid ID) (Permission, error) {
	return svc.db.GetPermission(uid)
}
