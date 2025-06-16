package scanner

type AiResultUnit map[string][]Unit

type Unit struct {
	Result string
	Reason string
}
