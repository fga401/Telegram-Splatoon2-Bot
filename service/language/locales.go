package language

// Case JSON structure
type Case struct {
	Cond string `json:"cond"`
	Text string `json:"text"`
}

// Plural JSON structure
type Plural struct {
	Arg    int    `json:"arg"`
	Format string `json:"format"`
	Cases  []Case `json:"cases"`
}

// Var JSON structure
type Var struct {
	Key    string `json:"key"`
	Plural Plural `json:"plural"`
}

// Message JSON structure
type Message struct {
	Key  string `json:"key"`
	Text string `json:"text"`
	Vars []Var  `json:"vars"`
}

type langPackage []Message
