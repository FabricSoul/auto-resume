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

type FloatDismissMsg struct{}

type ShowFloatInputMsg struct {
	Prompt       string
	InitialValue string
	Callback     func(string)
}

type GenerationCompleteMsg struct{}
