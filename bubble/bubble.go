package bubble

import (
	"bytes"
	"embed"
	"fmt"
	"gebitang.com/maorum/colors"
	"gebitang.com/maorum/rum"
	"gebitang.com/maorum/tomato"
	"golang.org/x/image/font/opentype"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
	"log"
)

//go:embed STSONG.TTF
var f embed.FS

func itemToBubbleData(items []*tomato.DailyItem) (XYZs, string) {
	data := make(XYZs, len(items))
	total := 0
	st := ""
	for i, item := range items {

		xyz := tomato.ToXYZ(item.Stamp, item.Min)
		data[i].X = xyz.Hour
		data[i].Y = xyz.Minute
		data[i].Z = xyz.Radius

		total += item.Min
		st = item.Name
	}
	return data, fmt.Sprintf("%s：%d分钟", st, total)
}

func BuildBubble(items map[string][]*tomato.DailyItem) (bool, string) {
	b, s, bb := BuildBubbleByte(items)
	if b {
		return rum.PostToGroup(bb, s)
	} else {
		return b, s
	}
}

func BuildBubbleByte(items map[string][]*tomato.DailyItem) (bool, string, []byte) {
	ttf, _ := f.ReadFile("STSONG.TTF")
	fontTTF, err := opentype.Parse(ttf)
	if err != nil {
		log.Println(err)
		return false, err.Error(), []byte{0}
	}
	stSong := font.Font{Typeface: "简体中文"}
	font.DefaultCache.Add([]font.Face{
		{
			Font: stSong,
			Face: fontTTF,
		},
	})
	if !font.DefaultCache.Has(stSong) {
		log.Fatalf("no font %q!", stSong.Typeface)
	}
	plot.DefaultFont = stSong

	p := plot.New()

	p.Title.Text = "人生五味"
	p.X.Label.Text = "Hour*10"
	p.Y.Label.Text = "Minute"

	var w bytes.Buffer
	w.WriteString("合计：\n")
	for t, g := range items {
		if len(g) == 0 {
			continue
		}
		bubbleData, s := itemToBubbleData(g)
		w.WriteString(s)
		w.WriteString("\n")
		bs := NewBubbles(bubbleData, vg.Points(1), vg.Points(20))
		bs.Color, _ = colors.ToStdColor(ItemMap[t].HexColor)
		p.Add(bs)
		p.Legend.Add(ItemMap[t].Name, bs)
	}

	p.Legend.Left = true
	p.Legend.Top = true

	p.Legend.YPosition = -1

	to, err := p.WriterTo(10*vg.Inch, 6*vg.Inch, "png")
	if err != nil {
		return false, err.Error(), []byte{0}
	}
	var buf bytes.Buffer
	_, err = to.WriteTo(&buf)
	if err != nil {
		fmt.Println(err)
		return false, err.Error(), []byte{0}
	}
	return true, w.String(), buf.Bytes()
}
