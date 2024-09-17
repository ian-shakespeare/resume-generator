package generator

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"resumegenerator/pkg/resume"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func year(t *time.Time) string {
	if t == nil {
		return ""
	}
	return strconv.Itoa(t.Year())
}

func monthYear(t *time.Time) string {
	if t == nil {
		return ""
	}
	return fmt.Sprintf("%s %d", t.Month().String(), t.Year())
}

func GenerateHtml(r *resume.Resume, tmpl []byte) ([]byte, error) {
	t, err := template.New("tmpl").Funcs(template.FuncMap{
		"year":       year,
		"month_year": monthYear,
	}).Parse(string(tmpl))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if err = t.Execute(buf, r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func GeneratePdf(r *resume.Resume, tmpl []byte) ([]byte, error) {
	htmlContent, err := GenerateHtml(r, tmpl)
	if err != nil {
		return nil, err
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			page.SetDocumentContent(frameTree.Frame.ID, string(htmlContent)).Do(ctx)

			buf, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}); err != nil {
		return nil, err
	}

	return buf, nil
}
