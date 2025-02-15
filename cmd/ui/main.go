package main

import (
	"log"
	"os"

	"github.com/Prettyletto/post-dude/cmd/server"
	"github.com/Prettyletto/post-dude/cmd/ui/menu"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	s := server.NewServer(":8080")
	s.Start()

	m := menu.New()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	s.Stop()
}
