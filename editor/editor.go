package editor

import (
	"github.com/jargv/pq/project"
	logpkg "log"
	"os"
	"path"
	"strings"

	"github.com/nsf/termbox-go"
)

var log *logpkg.Logger

func init() {
	log = logpkg.New(os.Stderr, "", 0)
}

type Editor struct {
	loader
	current int
	cursor  int
}

func (e *Editor) moveCurrent(new int) {
	if new < e.Project.NTasks() && new >= 0 {
		e.current = new
	}
}

func (e *Editor) Edit() error {
	for {
		switch event := e.frame(e.Project); true {
		case event.Type == termbox.EventKey && event.Ch == 'q':
			return nil
		case event.Type == termbox.EventKey && event.Ch == 'h':
			if 0 < e.cursor {
				e.cursor -= 1
			} else {
				e.moveCurrent((e.current - 1) / 2)
				e.cursor = 0
			}
		case event.Type == termbox.EventKey && event.Ch == 'l':
			if 0 < e.cursor && e.cursor < len(e.Project.Get(e.current)) {
				e.cursor += 1
			} else {
				e.moveCurrent(e.current*2 + 1)
				e.cursor = 0
			}
		case event.Type == termbox.EventKey && event.Ch == 'j':
			e.moveCurrent(e.current + 1)
			e.cursor = 0
		case event.Type == termbox.EventKey && event.Ch == 'k':
			e.moveCurrent(e.current - 1)
			e.cursor = 0
		case event.Type == termbox.EventKey && event.Ch == 'A':
			e.cursor = len(e.Project.Get(e.current))
			e.edit()
		case event.Type == termbox.EventKey && event.Ch == 'i':
			e.edit()
		case event.Type == termbox.EventKey && event.Ch == 'D':
			item := e.Project.Get(e.current)
			item = item[:e.cursor]
			e.Project.Set(e.current, item)
			e.moveCurrent(e.current*2 + 1)
		case event.Type == termbox.EventKey && event.Ch == 'C':
			item := e.Project.Get(e.current)
			item = item[:e.cursor]
			e.Project.Set(e.current, item)
			e.moveCurrent(e.current*2 + 1)
			e.edit()
		case event.Type == termbox.EventKey && event.Ch == 'J':
			newCursor := e.Project.RunCompaction()
			e.moveCurrent(newCursor)
		case event.Type == termbox.EventKey && event.Ch == 'o':
			e.cursor = 0
			e.Project.New()
			e.moveCurrent(e.Project.NTasks() - 1)
			e.edit()
		case event.Type == termbox.EventKey && event.Ch == 'w':
			item := e.Project.Get(e.current)
			for ; e.cursor < len(item); e.cursor++ {
				if item[e.cursor] == ' ' {
					e.cursor++
					break
				}
			}
		case event.Type == termbox.EventKey && event.Ch == 'b':
			item := e.Project.Get(e.current)
			if e.cursor > 0 {
				e.cursor--
			}
			for e.cursor > 0 && item[e.cursor-1] != ' ' {
				e.cursor--
			}
		case event.Type == termbox.EventKey && event.Ch == '-':
			dir := path.Dir(e.filename)
			if e.isDir && path.Base(dir) != ".prio" {
				dir = path.Dir(dir)
			}
			e.current = 0
			e.Open(dir)
		case event.Type == termbox.EventKey && event.Key == termbox.KeyEsc:
			e.cursor = 0
		case event.Type == termbox.EventKey && event.Key == termbox.KeyEnter:
			if e.current > 0 {
				newCurrent := e.Project.PriorityUp(e.current)
				e.moveCurrent(newCurrent)
				continue
			}

			if !e.isDir {
				e.Project.Clear(0)
				e.current = 1
				continue
			}

			dir := strings.TrimSuffix(e.filename, ".index")
			path := path.Join(dir, e.Project.Get(e.current))
			e.current = 0
			e.cursor = 0
			err := e.Open(path)
			if err != nil {
				return err
			}
		default:
			log.Printf("unhandled event: %#v", event)
		}
	}
}

func (e *Editor) frame(project *project.Project) termbox.Event {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	project.Render(e)
	termbox.Sync()
	return termbox.PollEvent()
}

func (e *Editor) edit() {
	addChar := func(c string) {
		item := e.Project.Get(e.current)
		new := item[:e.cursor] + c + item[e.cursor:]
		e.Project.Set(e.current, new)
		e.cursor += 1
	}

	for {
		switch event := e.frame(e.Project); true {
		case event.Type == termbox.EventKey && event.Ch != 0:
			addChar(string(event.Ch))
		case event.Type == termbox.EventKey && event.Key == termbox.KeySpace:
			addChar(string(" "))
		case event.Type == termbox.EventKey && event.Key == termbox.KeyBackspace2:
			item := e.Project.Get(e.current)
			if len(item) != 0 {
				new := item[:e.cursor-1] + item[e.cursor:]
				e.Project.Set(e.current, new)
				e.cursor -= 1
			}
		case event.Type == termbox.EventKey && event.Key == termbox.KeyEsc:
			return
		default:
			return
		}
	}
}

func (e *Editor) RenderLabel(idx, x, y, width int, label string) {
	if e.current == idx {
		termbox.SetCursor(x+e.cursor, y)
	}
	pad := width - len(label)
	if pad < 0 {
		pad = 0
	}
	printAt(x, y, label+strings.Repeat("-", pad))
}

func (e *Editor) RenderArm(x, y, len int) {
	for i := 0; i < len; i++ {
		termbox.SetCell(x, y+i, '|', termbox.ColorWhite, termbox.ColorBlack)
	}
}

func printAt(x, y int, msg string) {
	for _, char := range msg {
		termbox.SetCell(x, y, char, termbox.ColorWhite, termbox.ColorBlack)
		x += 1
	}
}
