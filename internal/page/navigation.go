package page

import "log"

func (p *Page) NavigateToInconsistencies() error {
	if err := p.Click(`a[href="/#/inconsistencies"]`, false); err != nil {
		return nil
	}

	if err := p.Click(`.btn.btn-default[data-toggle]`, false); err != nil {
		return nil
	}

	p.Page.MustEval(`() => document.querySelectorAll('a.ng-binding')[16].id = "hundred-lines"`)

	if err := p.Click(`#hundred-lines`, false); err != nil {
		return nil
	}

	has, el, err := p.Page.Has(`document.querySelector('.beamerAnnouncementSnippet')`)
	if err != nil {
		log.Printf("error while trying to select %v", el)
	}
	if !has {
		p.Page.MustEval(`() => document.querySelector('.beamerAnnouncementSnippet').style.display="none"`)
	}

	return nil
}
