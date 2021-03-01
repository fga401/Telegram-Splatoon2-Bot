package language

type Case struct {
	Cond string `json:"cond"`
	Text string `json:"text"`
}

type Plural struct {
	Arg    int    `json:"arg"`
	Format string `json:"format"`
	Cases  []Case `json:"cases"`
}

type Var struct {
	Key    string `json:"key"`
	Plural Plural `json:"plural"`
}

type Message struct {
	Key  string `json:"key"`
	Text string `json:"text"`
	Vars []Var  `json:"vars"`
}

type langPackage []Message

