package types

type TransitionMsg struct {
	To     Appstate
	Params interface{}
}

type BackMsg struct {
}

type ErrorMsg struct {
	Error error
}
