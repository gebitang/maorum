package tomato

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// DailyItem such as: 1,2021-12-20 15:49:05,10,dGVzdCBsaW5lIGluZm8=
type (
	DailyItem struct {
		T       int
		Name    string
		Stamp   time.Time
		Min     int
		Comment string
	}
	XYZ struct {
		Hour   float64
		Minute float64
		Radius float64
	}
)

func (d *DailyItem) String() string {
	return fmt.Sprintf("%d %s %d %s", d.T, d.Stamp, d.Min, d.Comment)
}

func GetDailyItems(p string) map[int][]*DailyItem {
	group := make(map[int][]*DailyItem, 0)
	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		a := strings.Split(l, ",")
		di := &DailyItem{}
		// 1,2021-12-20 10:04:05,20,dGVzdCBsaW5lIGluZm8=
		for t := range a {
			switch t {
			case 0:
				di.T, _ = strconv.Atoi(a[t])
			case 1:
				di.Stamp, _ = strToTime(a[t])
			case 2:
				di.Min, err = strconv.Atoi(a[t])
			case 3:
				di.Comment = a[t]
			}
		}
		if g, found := group[di.T]; found {
			g = append(g, di)
			group[di.T] = g
		} else {
			group[di.T] = []*DailyItem{di}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return group
}

func ToXYZ(t time.Time, min int) XYZ {
	r := min / 2
	i := XYZ{}
	p := t.Add(time.Duration(-r) * time.Minute)
	i.Hour = float64(p.Hour() * 10)
	i.Minute = float64(p.Minute())
	i.Radius = minToRadius(min)
	return i
}

func minToRadius(min int) float64 {
	base := 10

	base = base + (min-base)/base

	return float64(base)
}

func temp() {
	decodeString, err := base64.StdEncoding.DecodeString("QjFDNzhCRTlGODY5MDUwREVDMjNDMzNGM0IyNTQwRDA=")
	if err != nil {
		return
	}
	fmt.Println(string(decodeString))
	str := base64.StdEncoding.EncodeToString(decodeString)
	fmt.Println(str)

	str = base64.StdEncoding.EncodeToString([]byte(`test line info`))
	fmt.Println(str)

}

func strToTime(str string) (time.Time, error) {
	layOut := "2006-01-02 15:04:05"
	return time.ParseInLocation(layOut, str, time.Local)
}
