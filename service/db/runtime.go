package db

//type RuntimeTableInterface interface {
//	TableInterface
//	InsertRuntime(runtime Runtime) error
//	UpdateRuntime(runtime Runtime) error
//	GetRuntime(uid int64) (Runtime, error)
//}

type RuntimeTableImpl struct {
	TableImpl
}

const (
	runtimeNamedStmtInsert        namedStmtName = iota
	runtimeNamedStmtUpdate        namedStmtName = iota
	runtimeNamedStmtUpdateAccount namedStmtName = iota
)

const (
	runtimeStmtSelectByUid    stmtName = iota
	runtimeStmtUpdateLanguage stmtName = iota
)

var runtimeNamedStmts = map[namedStmtName]Declaration{
	runtimeNamedStmtInsert: {false, "INSERT INTO runtime (uid, session_token, iksm, language) VALUES (:uid, :session_token, :iksm, :language);"},
	runtimeNamedStmtUpdate: {false, "UPDATE runtime SET session_token=:session_token, iksm=:iksm, language=:language WHERE uid=:uid;"},
	runtimeNamedStmtUpdateAccount: {false, "UPDATE runtime SET session_token=:session_token, iksm=:iksm WHERE uid=:uid;"},
}

var runtimeStmts = map[stmtName]Declaration{
	runtimeStmtSelectByUid:    {true, "SELECT * FROM runtime WHERE uid=?;"},
	runtimeStmtUpdateLanguage: {false, "UPDATE runtime SET language=? WHERE uid=?;"},
}

func (impl *RuntimeTableImpl) InsertRuntime(runtime *Runtime) error {
	return impl.namedExec(runtimeNamedStmtInsert, runtime)
}

func (impl *RuntimeTableImpl) UpdateRuntime(runtime *Runtime) error {
	return impl.namedExec(runtimeNamedStmtUpdate, runtime)
}

func (impl *RuntimeTableImpl) UpdateRuntimeAccount(runtime *Runtime) error {
	return impl.namedExec(runtimeNamedStmtUpdateAccount, runtime)
}

func (impl *RuntimeTableImpl) UpdateRuntimeLanguage(userID int64, language string) error {
	return impl.exec(runtimeStmtUpdateLanguage, language, userID)
}

func (impl *RuntimeTableImpl) GetRuntime(uid int64) (*Runtime, error) {
	runtime := &Runtime{}
	err := impl.get(runtimeStmtSelectByUid, runtime, uid)
	return runtime, err
}
