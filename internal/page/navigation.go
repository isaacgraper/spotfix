package page

func (p *Page) NavigateToInconsistencies() error {
	if err := p.Click(`a[href="/#/inconsistencies"]`, true); err != nil {
		return nil
	}

	if err := p.Click(`.btn.btn-default[data-toggle]`, true); err != nil {
		return nil
	}

	p.Page.MustEval(`() => document.querySelectorAll('a.ng-binding')[16].id = "hundred-lines"`)

	if err := p.Click(`#hundred-lines`, true); err != nil {
		return nil
	}
	return nil
}
