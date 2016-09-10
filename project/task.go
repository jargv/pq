package project

type Task struct {
	Title     string
	fileOrder int
	priority  int
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

type byFileOrder []*Task

func (ft byFileOrder) Len() int {
	return len(ft)
}

func (ft byFileOrder) Less(i, j int) bool {
	return ft[i].fileOrder < ft[j].fileOrder
}

func (ft byFileOrder) Swap(i, j int) {
	tmp := ft[i]
	ft[i] = ft[j]
	ft[j] = tmp
}

type byPriority []*Task

func (ft byPriority) Len() int {
	return len(ft)
}

func (ft byPriority) Less(i, j int) bool {
	return ft[i].priority < ft[j].priority
}

func (ft byPriority) Swap(i, j int) {
	tmp := ft[i]
	ft[i] = ft[j]
	ft[j] = tmp
}
