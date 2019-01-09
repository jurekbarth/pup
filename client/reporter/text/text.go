// Package text provides a reporter for humanized interactive events.
package text

import (
	"fmt"
	"time"

	"github.com/jurekbarth/pup/client/event"
	"github.com/jurekbarth/pup/client/internal/colors"
	spin "github.com/tj/go-spin"
	"github.com/tj/go/term"
)

// TODO: platform-specific reporting should live in the platform
// TODO: typed events would be nicer.. refactor event names
// TODO: refactor, this is a hot mess :D

// Report events.
func Report(events <-chan *event.Event) {
	r := reporter{
		events:       events,
		spinner:      spin.New(),
		pendingName:  "pendingName",
		pendingValue: "pendingValue",
	}

	r.Start()
}

// reporter struct.
type reporter struct {
	events       <-chan *event.Event
	spinner      *spin.Spinner
	prevTime     time.Time
	pendingName  string
	pendingValue string
}

// spin the spinner by moving to the start of the line and re-printing.
func (r *reporter) spin() {
	if r.pendingName != "" {
		r.pending(r.pendingName, r.pendingValue)
	}
}

// clear the liner.
func (r *reporter) clear() {
	r.pendingName = ""
	r.pendingValue = ""
	term.ClearLine()
}

// pending log with spinner.
func (r *reporter) pending(name, value string) {
	r.pendingName = name
	r.pendingValue = value
	term.ClearLine()
	fmt.Printf("\r   %s %s", colors.Purple(r.spinner.Next()+" "+name+":"), value)
}

// complete log with duration.
func (r *reporter) complete(name, value string, d time.Duration) {
	r.pendingName = ""
	r.pendingValue = ""
	term.ClearLine()
	duration := fmt.Sprintf("(%s)", d.Round(time.Millisecond))
	fmt.Printf("\r     %s %s %s\n", colors.Purple(name+":"), value, colors.Gray(duration))
}

// completeWithoutDuration log without duration.
func (r *reporter) completeWithoutDuration(name, value string) {
	r.pendingName = ""
	r.pendingValue = ""
	term.ClearLine()
	fmt.Printf("\r     %s %s\n", colors.Purple(name+":"), value)
}

// log line.
func (r *reporter) log(name, value string) {
	fmt.Printf("\r     %s %s\n", colors.Purple(name+":"), value)
}

// error line.
func (r *reporter) error(name, value string) {
	fmt.Printf("\r     %s %s\n", colors.Red(name+":"), value)
}

// Start handling events.
func (r *reporter) Start() {
	tick := time.NewTicker(150 * time.Millisecond)
	defer tick.Stop()
	defer fmt.Println("afterwards")
	for {
		select {
		case <-tick.C:
			r.spin()
		case e := <-r.events:
			switch e.Name {
			case "zip.start":
				term.HideCursor()
				r.pending("zip", "zipping files")
			case "zip.done":
				term.ShowCursor()
				r.completeWithoutDuration("zip", "complete zipping files")
			case "upload.start":
				term.HideCursor()
				r.pending("upload", "uploading zip")
			case "upload.done":
				term.ShowCursor()
				r.completeWithoutDuration("upload", "uploading done")
			case "logs.start":
				term.HideCursor()
				r.pending("logs", "wait for server logs")
			case "logs.done":
				term.ShowCursor()
				r.completeWithoutDuration("logs", "logs incoming...")
			case "s3download.start":
				term.HideCursor()
				r.pending("s3download", e.Value)
			case "s3download.done":
				term.ShowCursor()
				r.completeWithoutDuration("s3download", e.Value)
			case "unzip.start", "unzip.pending":
				term.HideCursor()
				r.pending("unzip", e.Value)
			case "unzip.done":
				term.ShowCursor()
				r.completeWithoutDuration("unzip", e.Value)
			case "usermanagement.start":
				term.HideCursor()
				r.pending("usermanagement", e.Value)
			case "usermanagement.done":
				term.ShowCursor()
				r.completeWithoutDuration("usermanagement", e.Value)
			case "ddb.start", "ddb.pending":
				term.HideCursor()
				r.pending("dynamodb", e.Value)
			case "ddb.done":
				term.ShowCursor()
				r.completeWithoutDuration("dynamodb", e.Value)
			case "lambda.start", "lambda.pending":
				term.HideCursor()
				r.pending("lambda", e.Value)
			case "lambda.done":
				term.ShowCursor()
				r.completeWithoutDuration("lambda", e.Value)
			case "cloudfrontupdate.start", "cloudfrontupdate.pending":
				term.HideCursor()
				r.pending("cloudfrontupdate", e.Value)
			case "cloudfrontupdate.done":
				term.ShowCursor()
				r.completeWithoutDuration("cloudfrontupdate", e.Value)
			case "s3.start", "s3.pending":
				term.HideCursor()
				r.pending("s3", e.Value)
			case "s3.done":
				term.ShowCursor()
				r.completeWithoutDuration("s3", e.Value)
			case "cloudfront.start":
				term.HideCursor()
				r.pending("cloudfront", e.Value)
			case "cloudfront.done":
				term.ShowCursor()
				r.completeWithoutDuration("cloudfront", e.Value)
			case "sqs.start":
				term.HideCursor()
				r.pending("sqs", e.Value)
			case "sqs.done":
				term.ShowCursor()
				r.completeWithoutDuration("sqs", e.Value)
			case "cf.start":
				term.HideCursor()
				r.pending("cloudfront", "Waiting (approx. 20-40min)")
			case "cf.done":
				term.ShowCursor()
				r.completeWithoutDuration("cloudfront", "all edges updated")
			case "cfinvalidation.start":
				term.HideCursor()
				r.pending("cloudfront", "Invalidating Cache Waiting (approx. 2-5min)")
			case "cfinvalidation.done":
				term.ShowCursor()
				r.completeWithoutDuration("cloudfront", "Invalidated Cache")
			}

			r.prevTime = time.Now()
		}
	}
}

// actionColor returns a color func by action.
func actionColor(s string) colors.Func {
	switch s {
	case "Add":
		return colors.Purple
	case "Remove":
		return colors.Red
	case "Modify":
		return colors.Blue
	default:
		return colors.Gray
	}
}
