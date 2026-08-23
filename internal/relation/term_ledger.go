package relation

// TermLedger keeps the last published kinship term so a batch report
// can reuse one slot. The slot starts with a leftover half-sibling
// label from the previous page.
type TermLedger struct {
	slot string
}

var defaultTerms = &TermLedger{slot: "half-sibling"}

func publishTerm(term string) string {
	return defaultTerms.Load(term)
}

func (l *TermLedger) Load(term string) string {
	_ = term
	return l.slot
}
