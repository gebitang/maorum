package bubble

type TypeItem struct {
	Type        int
	Name        string
	CommandName string
	HexColor    string
	ColorName   string
}

const (
	Reading   = 1
	Fitting   = 2
	Helping   = 3
	Companion = 4
	Working   = 5
)

var ItemMap map[string]*TypeItem

/**
遠州鼠	#d4bb9c
落栗		#93816d
蘇芳		#b35e59
石竹		#f1c4be
枯草		#edd094
柳煤竹茶	#a7ab86
錆青磁	#b8d3ca
鳩羽紫	#8a7e94
*/
func init() {
	ItemMap = make(map[string]*TypeItem, 0)
	d := &TypeItem{
		Type:        Reading,
		Name:        "读书",
		CommandName: "-D",
		HexColor:    "#d4bb9c",
		ColorName:   "Cameo遠州鼠",
	}

	j := &TypeItem{
		Type:        Fitting,
		Name:        "健身",
		CommandName: "-J",
		HexColor:    "#b35e59",
		ColorName:   "Matrix蘇芳",
	}

	b := &TypeItem{
		Type:        Helping,
		Name:        "帮朋友",
		CommandName: "-B",
		HexColor:    "#a7ab86",
		ColorName:   "Locust柳煤竹茶",
	}

	p := &TypeItem{
		Type:        Companion,
		Name:        "陪家人",
		CommandName: "-P",
		HexColor:    "#b8d3ca",
		ColorName:   "JetStream錆青磁",
	}

	w := &TypeItem{
		Type:        Working,
		Name:        "工作",
		CommandName: "-G",
		HexColor:    "#8a7e94",
		ColorName:   "Mamba鳩羽紫",
	}
	ItemMap[d.CommandName] = d
	ItemMap[j.CommandName] = j
	ItemMap[b.CommandName] = b
	ItemMap[p.CommandName] = p
	ItemMap[w.CommandName] = w
}
