package db

type Runtime struct {
	Uid          int64  `db:"uid"`
	SessionToken string `db:"session_token"`
	IKSM         string `db:"iksm"`
	Language     string `db:"language"`
	Timezone     int    `db:"timezone"`
}

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
	runtimeStmtUpdateTimezone stmtName = iota
	runtimeStmtSelectFirst    stmtName = iota
	runtimeStmtUpdateIKSM     stmtName = iota
)

var runtimeNamedStmts = map[namedStmtName]Declaration{
	runtimeNamedStmtInsert:        {false, "INSERT INTO runtime (uid, session_token, iksm, language, timezone) VALUES (:uid, :session_token, :iksm, :language, :timezone);"},
	runtimeNamedStmtUpdate:        {false, "UPDATE runtime SET session_token=:session_token, iksm=:iksm, language=:language, timezone=:timezone WHERE uid=:uid;"},
	runtimeNamedStmtUpdateAccount: {false, "UPDATE runtime SET session_token=:session_token, iksm=:iksm WHERE uid=:uid;"},
}

var runtimeStmts = map[stmtName]Declaration{
	runtimeStmtSelectByUid:    {true, "SELECT * FROM runtime WHERE uid=?;"},
	runtimeStmtUpdateLanguage: {false, "UPDATE runtime SET language=? WHERE uid=?;"},
	runtimeStmtUpdateTimezone: {false, "UPDATE runtime SET timezone=? WHERE uid=?;"},
	runtimeStmtSelectFirst:    {false, `SELECT * FROM runtime WHERE session_token!="" LIMIT 1;`},
	runtimeStmtUpdateIKSM:     {true, "UPDATE runtime SET iksm=? WHERE uid=?;"},
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

func (impl *RuntimeTableImpl) UpdateRuntimeTimezone(userID int64, timezone int) error {
	return impl.exec(runtimeStmtUpdateTimezone, timezone, userID)
}

func (impl *RuntimeTableImpl) GetRuntime(uid int64) (*Runtime, error) {
	runtime := &Runtime{}
	err := impl.get(runtimeStmtSelectByUid, runtime, uid)
	return runtime, err
}

func (impl *RuntimeTableImpl) GetFirstRuntime() (*Runtime, error) {
	runtime := &Runtime{}
	err := impl.get(runtimeStmtSelectFirst, runtime)
	return runtime, err
}

func (impl *RuntimeTableImpl) UpdateRuntimeIKSM(userID int64, iksm string) error {
	return impl.exec(runtimeStmtUpdateIKSM, iksm, userID)
}
