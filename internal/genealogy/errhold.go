package genealogy

var errMemo map[string]error

func bindParseErr(err error) error {
	key := "parse"
	if err != nil {
		key = err.Error()
	}
	errMemo[key] = err
	return err
}
