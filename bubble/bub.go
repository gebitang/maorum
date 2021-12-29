package bubble

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"math"
)

// XYZer wraps the Len and XYZ methods.
type XYZer interface {
	// Len returns the number of x, y, z triples.
	Len() int

	// XYZ returns an x, y, z triple.
	XYZ(int) (float64, float64, float64)
}

// XYZs implements the XYZer interface using a slice.
type XYZs []struct{ X, Y, Z float64 }

// Len implements the Len method of the XYZer interface.
func (xyz XYZs) Len() int {
	return len(xyz)
}

// XYZ implements the XYZ method of the XYZer interface.
func (xyz XYZs) XYZ(i int) (float64, float64, float64) {
	return xyz[i].X, xyz[i].Y, xyz[i].Z
}

// CopyXYZs copies an XYZer.
func CopyXYZs(data XYZer) XYZs {
	cpy := make(XYZs, data.Len())
	for i := range cpy {
		cpy[i].X, cpy[i].Y, cpy[i].Z = data.XYZ(i)
	}
	return cpy
}

// Bubbles implements the Plotter interface, drawing
// a bubble plot of x, y, z triples where the z value
// determines the radius of the bubble.
type Bubbles struct {
	XYZs

	// Color is the color of the bubbles.
	color.Color

	// MinRadius and MaxRadius give the minimum and
	// maximum bubble radius respectively.  The radius
	// of each bubble is interpolated linearly between
	// these two values.
	MinRadius, MaxRadius vg.Length

	// MinZ and MaxZ are the minimum and maximum Z
	// values from the data.
	MinZ, MaxZ float64
}

// NewBubbles creates as new bubble plot plotter for
// the given data, with a minimum and maximum
// bubble radius.
func NewBubbles(xyz XYZer, min, max vg.Length) *Bubbles {
	cpy := CopyXYZs(xyz)
	minz := cpy[0].Z
	maxz := cpy[0].Z
	for _, d := range cpy {
		minz = math.Min(minz, d.Z)
		maxz = math.Max(maxz, d.Z)
	}
	return &Bubbles{
		XYZs:      cpy,
		MinRadius: min,
		MaxRadius: max,
		MinZ:      minz,
		MaxZ:      maxz,
	}
}

// Plot implements the Plot method of the plot.Plotter interface.
func (bs *Bubbles) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)

	c.SetColor(bs.Color)

	for _, d := range bs.XYZs {
		// Transform the data x, y coordinate of this bubble
		// to the corresponding drawing coordinate.
		x := trX(d.X)
		y := trY(d.Y)

		// Get the radius of this bubble.  The radius
		// is specified in drawing units (i.e., its size
		// is given as the final size at which it will
		// be drawn) so it does not need to be transformed.
		rad := vg.Length(d.Z)

		// Fill a circle centered at x,y on the draw area.
		var p vg.Path
		o := vg.Point{X: x + rad, Y: y}
		p.Move(o)
		o2 := vg.Point{X: x, Y: y}
		p.Arc(o2, rad, 0, 2*math.Pi)
		p.Close()
		c.Fill(p)
	}
}

// radius returns the radius of a bubble, in drawing
// units (vg.Lengths), by linear interpolation.
func (bs *Bubbles) radius(z float64) vg.Length {
	if bs.MinZ == bs.MaxZ {
		return (bs.MaxRadius-bs.MinRadius)/2 + bs.MinRadius
	}

	// Convert MinZ and MaxZ to vg.Lengths.  We just
	// want them to compute a slope so the units
	// don't matter, and the conversion is OK.
	minz := vg.Length(bs.MinZ)
	maxz := vg.Length(bs.MaxZ)

	slope := (bs.MaxRadius - bs.MinRadius) / (maxz - minz)
	intercept := bs.MaxRadius - maxz*slope
	return vg.Length(z)*slope + intercept
}

// XYValues implements the XYer interface, returning the
// x and y values from an XYZer.
type XYValues struct{ XYZer }

// XY implements the XY method of the XYer interface.
func (xy XYValues) XY(i int) (float64, float64) {
	x, y, _ := xy.XYZ(i)
	return x, y
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (bs *Bubbles) DataRange() (xmin, xmax, ymin, ymax float64) {
	// Note that by defining the XYValues type, which
	// implements the XYer interface, we can easily re-use
	// the XYRange function from the plotter package to
	// compute the minimum and maximum X and Y values.
	xm, xmax, ym, ymax := plotter.XYRange(XYValues{bs.XYZs})
	return xm - 20, xmax, ym, ymax + 30
}

func (bs *Bubbles) Thumbnail(c *draw.Canvas) {
	pts := []vg.Point{
		{X: c.Min.X, Y: c.Min.Y},
		{X: c.Min.X, Y: c.Max.Y},
		{X: c.Max.X, Y: c.Max.Y},
		{X: c.Max.X, Y: c.Min.Y},
	}
	poly := c.ClipPolygonY(pts)
	c.FillPolygon(bs.Color, poly)
	pts = append(pts, vg.Point{X: c.Min.X, Y: c.Min.Y})
	outline := c.ClipLinesY(pts)
	style := draw.LineStyle{Width: vg.Length(0)}
	c.StrokeLines(style, outline...)
}
