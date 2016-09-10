package project

type Task struct {
	Title string
}

func (t *Task) Clear() {
	t.Title = ""
}

func (t *Task) IsEmpty() bool {
	return t.Title == ""
}

func (t *Task) Len() int {
	return len(t.Title)
}
