package bot

import (
	"github.com/isaacgraper/spotfix.git/internal/common/config"
	"github.com/isaacgraper/spotfix.git/internal/page"
)

type Process struct {
	config *config.Config
	page   *page.Page
}

func (p *Process) ProcessHandler(c *config.Config) error {
	return nil
}

func (p *Process) ProcessBatch(start, end int) (error, bool) {
	return nil, false
}

func (p *Process) ProcessFilter() (error, bool) {
	return nil, false
}

func (p *Process) EndProcess() (error, bool) {
	return nil, false
}
