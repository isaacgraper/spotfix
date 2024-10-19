package page

import (
	"fmt"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

type Page struct {
	Page *rod.Page
}

// func (p *Page) NewPage(page *rod.Page) *Page {
// 	return &Page{
// 		Page: page,
// 	}
// }

func (p *Page) Click(selector string, screenshot bool) error {
	err := rod.Try(func() {
		element, err := p.Page.Element(selector)
		if err != nil {
			log.Printf("Element not found: %s", selector)
			return
		}

		p.Loading()

		err = element.Click(proto.InputMouseButtonLeft, 1)
		if err != nil {
			log.Printf("Failed to click element: %s", selector)
			return
		}

		time.Sleep(time.Millisecond * 200)

		if screenshot {
			p.Page.MustScreenshot(fmt.Sprintf("screenshot_%d.png", time.Now().Unix()))
		}
	})

	return err
}

func (p *Page) ClickWithRetry(selector string, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = rod.Try(func() {
			element, err := p.Page.Timeout(250 * time.Millisecond).Element(selector)
			if err != nil {
				log.Printf("Element not found: %s", selector)
				return
			}
			err = element.Click(proto.InputMouseButtonLeft, 1)
			if err != nil {
				log.Printf("Failed to click element: %s", selector)
				return
			}
		})

		if err == nil {
			return nil
		}

		log.Printf("Attempt %d failed to click on %s: %v", i+1, selector, err)

		time.Sleep(time.Second)
	}

	return fmt.Errorf("failed to click on %s after %d attempts: %v", selector, maxRetries, err)
}

func (p *Page) AddElementId(selector, id string) error {
	p.Page.Eval(fmt.Sprintf(`() => {
		const el = %s;
		if (el) {
		el.id = "%s";
		} else {
			console.error("Element not found:", %s);
		}
	}`, selector, id, selector))

	return nil
}

func (p *Page) ScrollToElement(selector string) error {
	p.Page.Eval(fmt.Sprintf(`() => {
		const element = document.querySelector('%s');
		if (element) {
			element.scrollIntoView({ behavior: 'smooth', block: 'center' });
		} else {
			console.error("Element not found:", '%s');
		}
	}`, selector, selector))

	return nil
}

func (p *Page) Pagination() (error, bool) {

	// implement error handling here
	hasNextPage := p.Page.MustHas(`[ng-click="changePage('next')"]`)
	if !hasNextPage {
		log.Println("No next page found or error occurred:")
		return nil, false
	}

	if err := p.ClickWithRetry(`[ng-click="changePage('next')"]`, 3); err != nil {
		log.Println("Failed to click next page button:", err)
		return nil, false
	}

	p.Loading()

	log.Println("Moved to next page")

	return nil, true
}

func (p *Page) Filter() error {
	if err := p.ClickWithRetry(`#inconsistenciesFilter`, 3); err != nil {
		return fmt.Errorf("failed to click inconsistencies filter: %w", err)
	}

	element, err := p.Page.Element(`select#clockingTypes`)
	if err != nil {
		return fmt.Errorf("failed to find clocking types element: %w", err)
	}

	element.MustWaitStable()

	if err := p.ClickWithRetry(`#clockingTypes`, 3); err != nil {
		return fmt.Errorf("failed to click clocking types: %w", err)
	}

	p.Loading()

	err = element.Type(input.ArrowDown)
	if err != nil {
		return fmt.Errorf("failed to type arrow down: %w", err)
	}

	err = element.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("failed to click selected option: %w", err)
	}

	if err := p.ClickWithRetry(`#app > searchfilterinconsistencies > div > div.row.overscreen_child > div.filter_container > div.hbox.filter_button.ng-scope > a.btn.button_link.btn-dark.ng-binding`, 3); err != nil {
		return fmt.Errorf("failed to apply filter: %w", err)
	}

	p.Loading()

	log.Println("Page filtered")

	return nil
}

func (p *Page) Loading() {
	p.Page.MustWaitLoad().MustWaitStable().MustWaitDOMStable()
}
