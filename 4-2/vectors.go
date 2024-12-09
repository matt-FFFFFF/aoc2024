package main

import "strings"

type vector struct {
	vert  int
	horiz int
}

var vectorN = vector{
	vert:  -1,
	horiz: 0,
}

var vectorNE = vector{
	vert:  -1,
	horiz: 1,
}

var vectorE = vector{
	vert:  0,
	horiz: 1,
}

var vectorSE = vector{
	vert:  1,
	horiz: 1,
}

var vectorS = vector{
	vert:  1,
	horiz: 0,
}

var vectorSW = vector{
	vert:  1,
	horiz: -1,
}

var vectorW = vector{
	vert:  0,
	horiz: -1,
}

var vectorNW = vector{
	vert:  -1,
	horiz: -1,
}

func (v vector) String() string {
	sb := strings.Builder{}
	switch v.vert {
	case 1:
		sb.WriteString("South")
	case -1:
		sb.WriteString("North")
	}
	switch v.horiz {
	case 1:
		sb.WriteString("East")
	case -1:
		sb.WriteString("West")
	}
	return sb.String()
}
