package bot

import (
	"log"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/isaacgraper/spotfix.git/internal/common/config"
	"github.com/isaacgraper/spotfix.git/internal/page"
)

type Process struct {
	config *config.Config
	page   *page.Page
}

func NewProcess() *Process {
	return &Process{
		config: &config.Config{},
		page:   &page.Page{},
	}
}

func (pr *Process) ProcessHandler(c *config.Config) error {
	browser := rod.New().ControlURL(launcher.New().Headless(false).MustLaunch()).MustConnect()
	defer browser.MustClose()

	// URL must not working as expected in my env file
	pageInstance := browser.MustPage("https://orbenk1.nexti.com/").MustWaitLoad()

	pr.page = &page.Page{
		Page: pageInstance,
	}

	if err := pr.page.Login(c.NewCredential()); err != nil {
		log.Printf("login failed: %v", err)
		return nil
	}
	return nil
}

func (pr *Process) ProcessBatch(start, end int) (error, bool) {
	return nil, false
}

func (pr *Process) ProcessFilter() (error, bool) {
	return nil, false
}

func (pr *Process) EndProcess() (error, bool) {
	return nil, false
}
