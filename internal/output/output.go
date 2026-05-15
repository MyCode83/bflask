package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"

	"bflask/pkg/bflask"
)

type Printer struct {
	json  bool
	quiet bool
	out   io.Writer
}

func New(jsonMode, quiet bool) *Printer {
	color.NoColor = quiet
	return &Printer{json: jsonMode, quiet: quiet, out: os.Stdout}
}

func (p *Printer) Banner() {
	if p.quiet || p.json {
		return
	}
	fmt.Fprintln(p.out, color.HiCyanString("bflask"), color.WhiteString("authorized Flask session cookie tester"))
}

func (p *Printer) Found(result bflask.Result) error {
	if p.json {
		enc := json.NewEncoder(p.out)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	if p.quiet {
		fmt.Fprintln(p.out, result.SecretKey)
		return nil
	}

	fmt.Fprintf(p.out, "%s SECRET_KEY found: %s\n", color.HiGreenString("[HIT]"), color.HiWhiteString(result.SecretKey))
	if result.Payload != "" {
		fmt.Fprintln(p.out, color.HiBlueString("[INF]"), "Decoded payload:")
		fmt.Fprintln(p.out, result.Payload)
	}
	return nil
}

func (p *Printer) NotFound(stats bflask.Stats) error {
	if p.json {
		return json.NewEncoder(p.out).Encode(map[string]any{
			"found":   false,
			"checked": stats.Checked,
			"elapsed": stats.Elapsed.String(),
		})
	}
	if !p.quiet {
		fmt.Fprintf(p.out, "%s No valid SECRET_KEY found after %d candidates\n", color.HiYellowString("[WRN]"), stats.Checked)
	}
	return nil
}
