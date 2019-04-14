package main

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
)

// redrawInterval is how often termdash redraws the screen.
const redrawInterval = 500 * time.Millisecond

// widgets holds the widgets used by this demo.
type widgets struct {
	doneText   *text.Text
	errorsText *text.Text
	heartLC    *linechart.LineChart
}

// newWidgets creates all widgets used by this demo.
func newWidgets(ctx context.Context, c *container.Container) (*widgets, error) {
	doneT, err := newRollText(ctx)
	checkErr(err)
	errorsT, err := newRollText(ctx)
	checkErr(err)

	heartLC, err := newHeartbeat(ctx)
	checkErr(err)

	return &widgets{
		doneText:   doneT,
		errorsText: errorsT,
		heartLC:    heartLC,
	}, nil
}

// layoutType represents the possible layouts the buttons switch between.
type layoutType int

const (
	// layoutAll displays all the widgets.
	layoutAll layoutType = iota
	// layoutText focuses onto the text widget.
	layoutText
	// layoutSparkLines focuses onto the sparklines.
	layoutSparkLines
	// layoutLineChart focuses onto the linechart.
	layoutLineChart
)

// gridLayout prepares container options that represent the desired screen layout.
// This function demonstrates the use of the grid builder.
// gridLayout() and contLayout() demonstrate the two available layout APIs and
// both produce equivalent layouts for layoutType layoutAll.
func gridLayout(w *widgets) ([]container.Option, error) {
	leftRows := []grid.Element{
		grid.RowHeightPerc(59,
			grid.ColWidthPerc(99,
				grid.Widget(w.doneText,
					container.Border(linestyle.Light),
					container.BorderTitle("Current Request"),
				),
			),
		),
		grid.RowHeightPerc(40,
			grid.Widget(w.heartLC,
				container.Border(linestyle.Light),
				container.BorderTitle("RPM (Request per minute)"),
			),
		),
	}

	rightRows := []grid.Element{
		grid.RowHeightPerc(99,
			grid.ColWidthPerc(99,
				grid.Widget(w.errorsText,
					container.Border(linestyle.Light),
					container.BorderTitle("Errors"),
				),
			),
		),
	}

	builder := grid.New()
	builder.Add(
		grid.ColWidthPerc(70, leftRows...),
		grid.ColWidthPerc(30, rightRows...),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}

// rootID is the ID assigned to the root container.
const rootID = "root"

// periodic executes the provided closure periodically every interval.
// Exits when the context expires.
func periodic(ctx context.Context, interval time.Duration, fn func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := fn(); err != nil {
				panic(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// newRollText creates a new Text widget that displays rolling text.
func newRollText(ctx context.Context) (*text.Text, error) {
	t, err := text.New(text.RollContent())
	if err != nil {
		return nil, err
	}
	return t, nil
}

// newHeartbeat returns a line chart that displays a heartbeat-like progression.
func newHeartbeat(ctx context.Context) (*linechart.LineChart, error) {
	var inputs []float64
	for i := 0; i < 100; i++ {
		v := math.Pow(math.Sin(float64(i)), 63) * math.Sin(float64(i)+1.5) * 8
		inputs = append(inputs, v)
	}

	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorGreen)),
	)
	if err != nil {
		return nil, err
	}
	step := 0
	go periodic(ctx, redrawInterval/3, func() error {
		step = (step + 1) % len(inputs)
		return lc.Series("heartbeat", rotateFloats(inputs, step),
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(87))),
			linechart.SeriesXLabels(map[int]string{
				0: "zero",
			}),
		)
	})
	return lc, nil
}

// distance is a thread-safe int value used by the newSince method.
// Buttons write it and the line chart reads it.
type distance struct {
	v  int
	mu sync.Mutex
}

// add adds the provided value to the one stored.
func (d *distance) add(v int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.v += v
}

// get returns the current value.
func (d *distance) get() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.v
}

// rotateFloats returns a new slice with inputs rotated by step.
// I.e. for a step of one:
//   inputs[0] -> inputs[len(inputs)-1]
//   inputs[1] -> inputs[0]
// And so on.
func rotateFloats(inputs []float64, step int) []float64 {
	return append(inputs[step:], inputs[:step]...)
}

// rotateRunes returns a new slice with inputs rotated by step.
// I.e. for a step of one:
//   inputs[0] -> inputs[len(inputs)-1]
//   inputs[1] -> inputs[0]
// And so on.
func rotateRunes(inputs []rune, step int) []rune {
	return append(inputs[step:], inputs[:step]...)
}

func addText(url string, t *text.Text) error {
	if err := t.Write(fmt.Sprintf(url+"\n"), text.WriteCellOpts(cell.FgColor(cell.ColorNumber(142)))); err != nil {
		return err
	}
	return nil
}
