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
	items []string
}

func New(items []string) *Project {
	return &Project{
		items: items,
	}
}

func (p *Project) Items() []string {
	items := p.items
	l := len(items)
	return items[:l:l]
}

func (p *Project) Set(i int, s string) {
	p.items[i] = s
}

func (p *Project) Get(i int) string {
	return p.items[i]
}

func (p *Project) NTasks() int {
	return len(p.items)
}

func (p *Project) Clear(idx int) {
	p.items[idx] = ""
}

func (p *Project) New() {
	p.items = append(p.items, "")
}

func (p *Project) PriorityUp(idx int) int {
	parent := (idx - 1) / 2
	tmp := p.items[idx]
	if tmp == "" {
		return -1
	}
	p.items[idx] = p.items[parent]
	p.items[parent] = tmp
	if p.items[idx] == "" {
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
		if idx >= len(p.items) {
			return -1
		}

		armStart := visit(idx*2+1, tier+1)

		x := offsets[tier] + tier
		width := len(p.items[idx])
		if tier+1 < len(offsets) {
			width = offsets[tier+1] - offsets[tier]
		}

		var pad int
		if armStart >= 0 {
			pad = width
		} else {
			pad = len(p.items[idx])
		}

		renderer.RenderLabel(idx, x, y, pad, p.items[idx])
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
	if len(p.items) == 0 {
		return -1
	}

	//just drop any spaces at the end
	for len(p.items) > 0 && p.items[len(p.items)-1] == "" {
		p.items = p.items[:len(p.items)-1]
	}

	//now find the first space
	firstEmpty := 0
	for ; firstEmpty < len(p.items); firstEmpty++ {
		if p.items[firstEmpty] == "" {
			break
		}
	}

	if firstEmpty == len(p.items) {
		return -1
	}

	p.items[firstEmpty] = p.items[len(p.items)-1]
	p.items = p.items[:len(p.items)-1]
	cursor := firstEmpty

	//drop any new spaces
	for p.items[len(p.items)-1] == "" {
		p.items = p.items[:len(p.items)-1]
	}

	if cursor >= len(p.items) {
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
		if count+start >= len(p.items) {
			count = len(p.items) - start
			last = true
		}

		end := start + count

		width := 0
		for _, word := range p.items[start:end] {
			if len(word) > width {
				width = len(word)
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
