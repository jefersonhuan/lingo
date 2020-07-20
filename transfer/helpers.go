package transfer

import (
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"strings"
)

func getStat(field *float64, value interface{}) {
	switch value.(type) {
	case float64:
		*field = value.(float64)
	case int32:
		*field = float64(value.(int32))
	}
}

func pushError(err error) {
	failures = append(failures, err)
}

func startBarForCollection(name string, total int64, p *mpb.Progress) *mpb.Bar {
	if len(name) < barTitleWidth {
		name += strings.Repeat(" ", barTitleWidth-len(name))
	}

	return p.AddBar(total,
		mpb.PrependDecorators(
			decor.Name(name),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60), "finished",
			),
		),
	)
}
