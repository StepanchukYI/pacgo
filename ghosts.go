package main

import "sync"

type GhostStatus string

const (
	GhostStatusNormal GhostStatus = "Normal"
	GhostStatusBlue   GhostStatus = "Blue"
)

type ghost struct {
	position sprite
	status   GhostStatus
}

var ghosts []*ghost

var ghostsStatusMx sync.RWMutex

func updateGhosts(ghosts []*ghost, ghostStatus GhostStatus) {
	ghostsStatusMx.Lock()
	defer ghostsStatusMx.Unlock()
	for _, g := range ghosts {
		g.status = ghostStatus
	}
}

func moveGhosts() {
	for _, g := range ghosts {
		dir := drawDirection()
		g.position.row, g.position.col = makeMove(g.position.row, g.position.col, dir)
	}
}
