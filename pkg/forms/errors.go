package forms

type errors map[string][]string

func (e errors) Add(field, msg string) {
	e[field] = append(e[field], msg)
}

func (e errors) Get(field string) string {
	errmsg := e[field]
	if len(errmsg) == 0 {
		return ""
	}
	return errmsg[0]
}
