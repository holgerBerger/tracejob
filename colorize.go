package main

import (
	"fmt"
	"sort"
	"strings"
)

type RGB struct {
	r, g, b int
}

var colormap map[string]RGB
var cls []string
var red, green, blue, purple, orange RGB

func initcolors() {
	if opts.Light {
		red = RGB{200, 0, 0}
		green = RGB{0, 200, 0}
		blue = RGB{0, 0, 200}
		purple = RGB{100, 100, 100}
		orange = RGB{255, 100, 0}
	}
	if opts.Dark {
		red = RGB{255, 0, 0}
		green = RGB{0, 255, 0}
		blue = RGB{0, 0, 255}
		purple = RGB{200, 200, 200}
		orange = RGB{255, 100, 0}
	}
	colormap = map[string]RGB{
		"SUSPEND":         red,
		"resume":          red,
		"killed":          red,
		"Failed":          red,
		" error ":         red,
		" exceed":         red,
		"FAIL":            red,
		"SYSTEM_FAILURE":  red,
		"SUSPEND_FAIL":    red,
		"ARRIVE_FAIL":     red,
		"PRERUN_FAIL":     red,
		"ERROR":           red,
		"completed":       green,
		"SUCCEED":         green,
		"Output_reqinfo":  blue,
		" SIG":            blue,
		"RST_ARRIVING":    blue,
		"RST_EXITED":      blue,
		"RST_EXITING":     blue,
		"RST_HELD":        blue,
		"RST_MOVED":       blue,
		"RST_POSTRUNNING": blue,
		"RST_PRERUNNING":  blue,
		"RST_QUEUED":      blue,
		"RST_RESUMING":    blue,
		"RST_REVIVED":     blue,
		"RST_RUNNING":     blue,
		"RST_STAGING":     blue,
		"RST_STALLED":     blue,
		"RST_SUSPENDED":   blue,
		"RST_SUSPENDING":  blue,
		"RST_WAITING":     blue,
		"JST_DELETED":     blue,
		"EXITING":         blue,
		"JM_ASSIGNED":     blue,
		"PRERUNNING":      blue,
		"STAGING":         blue,
		"JM_QUEUED":       blue,
		"EXITED":          blue,
		"RUNNING":         blue,
		"POSTRUNNING":     blue,
		"INFO":            purple,
		"(NOTUSE)":        orange,
		"(RUN)":           orange,
		"(SCHED)":         orange,
		"(STAGING)":       orange,
		"(STG_DELAY)":     orange,
		"(WAIT_RUN)":      orange,
		"(EXEC)":          orange,
		"(EXT_DELAY)":     orange,
		"rst=":            orange,
		"exst":            orange,
		"terminated":      orange,
	}
	for s := range colormap {
		cls = append(cls, s)
	}

	sort.SliceStable(cls, func(i, j int) bool {
		return len(cls[i]) > len(cls[j])
	})

}

func colorize(line *string) {

	for _, s := range cls {
		*line = strings.ReplaceAll(*line, s, fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", colormap[s].r, colormap[s].g, colormap[s].b, s))
	}
}
