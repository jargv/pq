package project

import (
	logpkg "log"
	"os"
)

var log *logpkg.Logger

func init() {
	log = logpkg.New(os.Stderr, "", 0)
}

type Renderer interface {
	RenderLabel(i, x, y, width int, msg string)
	RenderArm(x, y, len int)
}

type Project struct {
	tasks []*Task
}

func New(taskStrings []string) *Project {
	tasks := make([]*Task, len(taskStrings))
	for i, str := range taskStrings {
		tasks[i] = &Task{Title: str}
	}
	return &Project{
		tasks: tasks,
	}
}

func (p *Project) TaskStrings() []string {
	taskStrings := make([]string, len(p.tasks))
	for i, task := range p.tasks {
		taskStrings[i] = task.Title
	}
	return taskStrings
}

func (p *Project) Task(i int) *Task {
	return p.tasks[i]
}

func (p *Project) NTasks() int {
	return len(p.tasks)
}

func (p *Project) New() *Task {
	task := &Task{}
	p.tasks = append(p.tasks, task)
	return task
}

func (p *Project) PriorityUp(idx int) int {
	parent := (idx - 1) / 2
	tmp := p.tasks[idx]
	if tmp.IsEmpty() {
		return -1
	}
	p.tasks[idx] = p.tasks[parent]
	p.tasks[parent] = tmp
	if p.tasks[idx].IsEmpty() {
		return idx*2 + 1
	} else {
		return parent
	}
}

func (p *Project) Render(renderer Renderer) {
	offsets := p.computeOffsets()
	y := 0

	var visit func(int, int) int
	visit = func(idx, tier int) int {
		if idx >= len(p.tasks) {
			return -1
		}

		armStart := visit(idx*2+1, tier+1)

		x := offsets[tier] + tier
		width := len(p.tasks[idx].Title)
		if tier+1 < len(offsets) {
			width = offsets[tier+1] - offsets[tier]
		}

		var pad int
		if armStart >= 0 {
			pad = width
		} else {
			pad = len(p.tasks[idx].Title)
		}

		renderer.RenderLabel(idx, x, y, pad, p.tasks[idx].Title)
		renderY := y
		y += 1

		armEnd := visit(idx*2+2, tier+1)
		if armEnd == -1 {
			armEnd = renderY
		}

		armLen := armEnd - armStart + 1
		if armLen < 0 {
			armLen = renderY - armStart + 1
		}
		if tier+1 < len(offsets) && armStart != -1 {
			renderer.RenderArm(offsets[tier+1]+tier, armStart, armLen)
		}

		return renderY
	}

	visit(0, 0)
}

func (p *Project) RunCompaction() int {
	if len(p.tasks) == 0 {
		return -1
	}

	//just drop any spaces at the end
	for len(p.tasks) > 0 && p.tasks[len(p.tasks)-1].Title == "" {
		p.tasks = p.tasks[:len(p.tasks)-1]
	}

	//now find the first space
	firstEmpty := 0
	for ; firstEmpty < len(p.tasks); firstEmpty++ {
		if p.tasks[firstEmpty].IsEmpty() {
			break
		}
	}

	if firstEmpty == len(p.tasks) {
		return len(p.tasks) - 1
	}

	p.tasks[firstEmpty] = p.tasks[len(p.tasks)-1]
	p.tasks = p.tasks[:len(p.tasks)-1]
	cursor := firstEmpty

	//drop any new spaces
	for p.tasks[len(p.tasks)-1].IsEmpty() {
		p.tasks = p.tasks[:len(p.tasks)-1]
	}

	if cursor >= len(p.tasks) {
		cursor = 0
	}

	return cursor
}

func (p *Project) computeOffsets() []int {
	offsets := []int{}
	total := 0
	start := 0
	for i := 0; ; i++ {
		count := 1 << uint(i)
		last := false
		if count+start >= len(p.tasks) {
			count = len(p.tasks) - start
			last = true
		}

		end := start + count

		width := 0
		for _, task := range p.tasks[start:end] {
			if len(task.Title) > width {
				width = len(task.Title)
			}
		}
		offsets = append(offsets, total)
		total += width

		start += count
		if last {
			break
		}
	}

	return offsets
}
