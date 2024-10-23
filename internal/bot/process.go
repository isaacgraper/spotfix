package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/isaacgraper/spotfix.git/internal/common"
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

func (pr *Process) Execute(c *config.Config) error {
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

	if err := pr.page.NavigateToInconsistencies(); err != nil {
		log.Printf("navigate to inconsistencies failed: %v", err)
		return nil
	}

	if c.Filter {
		if err := pr.page.Filter(); err != nil {
			log.Printf("filtering failed: %v", err)
			return nil
		}
		pr.ProcessFilter(c)
	}

	if !c.Filter {
		pr.ProcessHandler(c)
	}

	return nil
}

func (pr *Process) ProcessHandler(c *config.Config) (error, bool) {
	for {
		pr.ProcessResult(c)

		if !pr.page.Pagination() {
			log.Println("no more pages to process")
			break
		}
	}
	return nil, true
}

func (pr *Process) ProcessResult(c *config.Config) {
	if c.Max < 1 {
		log.Println("no results to process")
		return
	}

	pr.page.Page.MustEval(`() => {
		const el = document.querySelectorAll('[data-id]');
		for (let i = 1; i < el.length; i++) {
			el[i].id = 'inconsistence-' + i;
		}
	}`)

	batchSize := c.BatchSize
	for i := 0; i < c.Max; i += batchSize {
		end := i + batchSize
		if end > c.Max {
			end = c.Max
		}
		pr.ProcessBatch(i+1, end, c)
	}
	// implements SendEmail here before ending the process
	pr.EndProcess()
}

func (pr *Process) ProcessBatch(start, end int, c *config.Config) error {
	results := pr.page.Page.MustEval(fmt.Sprintf(`() => {
		const results = [];
		for (let i = %d; i <= %d; i++) {
			const row = document.querySelector('#inconsistence-' + i);
			if (row) {
				results.push({
					index: i,
					name: row.querySelector('td.ng-binding:nth-child(2)').textContent,
					hour: row.querySelector('td.ng-binding:nth-child(6)').textContent,
					category: row.querySelector('td.ng-binding:nth-child(7)').textContent,
				});
			}
		}
		return results;
	}`, start, end))

	pr.page.Loading()

	for _, result := range results.Arr() {
		category := result.Get("category").String()
		hour := result.Get("hour").String()
		name := result.Get("name").String()

		hourSplit := strings.Split(hour, " ")
		hour = strings.TrimSpace(hourSplit[1])

		if category != "Não registrado" {
			log.Printf("%s - %s - %s", name, hour, category)

			index := result.Get("index").Int()

			resultJSON, err := result.MarshalJSON()
			if err != nil {
				log.Printf("Erro ao converter result para JSON: %v", err)
				continue
			}

			concatStr := string(resultJSON)

			go func() {
				fileName := fmt.Sprintf("report-%s.txt", time.Now().Format("01-02-2006"))
				common.NewFile(fileName, []byte(concatStr)).SaveFile()
			}()

			pr.page.Loading()
			time.Sleep(time.Millisecond * 250)

			if err := pr.page.ClickWithRetry(fmt.Sprintf(`#inconsistence-%d i`, index), 6); err != nil {
				log.Printf("failed to click on inconsistency %v", err)
			}
		}
	}
	return nil
}

func (pr *Process) ProcessFilter(c *config.Config) {
	for {
		if err := pr.page.Click(`#content > div.app-content-body.nicescroll-continer > div.content-body > div.app-content-body > div.tab-lis > div.content-table > table > thead > tr > th:nth-child(1) > label > i`, false); err != nil {
			log.Printf("Failed to click filter checkbox: %v", err)
			break
		}
		pr.page.Loading()

		if pr.EndProcess() {
			if pr.page.Pagination() {
				continue
			}
		} else {
			break
		}
		pr.page.Loading()
	}
}

func (pr *Process) EndProcess() bool {
	pr.page.Page.MustElement(`td.ng-binding`).ScrollIntoView()

	elements := []string{
		`#content > div.app-content-body.nicescroll-continer > div.content-body > div.content-body-header > div.content-body-header-filters > div.filters-right > button`,
		`[btn-radio="\'CANCELED\'"]`,
		`#app > modal > div > div > div > div > div.modal-body > div > div > div:nth-child(2) > div > multiselect > div > div > div:nth-child(1) > div > i`,
		`[alt="Erro operacional"]`,
	}

	for _, selector := range elements {
		time.Sleep(time.Millisecond * 250)

		if err := pr.page.Click(selector, false); err != nil {
			log.Printf("Failed to click on %s: %v", selector, err)
			return false
		}
		pr.page.Loading()
	}

	note := pr.page.Page.MustElement(`input#note`)
	note.MustInput("Cancelamento automático via bot")

	if err := pr.page.Click(`a.btn.button_link.btn-primary.ng-binding`, false); err != nil {
		log.Printf("Failed to click on submit button: %v", err)
		return false
	}
	pr.page.Loading()

	log.Println("inconsistencies processed!")
	return true
}
