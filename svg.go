package main

// import (
// 	"net/http"

// 	svg "github.com/ajstarks/svgo"
// )

// const width = 1024
// const height = 768
// const x = 50
// const y = height / 2
// const intervalY = 200
// const intervalX = 100
// const style = "fill:none;stroke:black"

// func generateSVG(w http.ResponseWriter, r *http.Request, nginx map[string]map[string]nginxMeta) {

// 	w.Header().Set("Content-Type", "image/svg+xml")

// 	s := svg.New(w)
// 	s.Start(x, y)
// 	s.Title("nginx")

// 	s.Def()
// 	s.Marker("dot", 5, 5, 8, 8)
// 	s.Circle(5, 5, 3, "fill:black")
// 	// s.Text(x-15, y+6, "nginx")
// 	s.MarkerEnd()

// 	s.Marker("arrow", 2, 6, 13, 13)
// 	s.Path("M2,2 L2,11 L10,6 L2,2", "fill:blue")
// 	s.MarkerEnd()
// 	s.DefEnd()

// 	// idx := 1
// 	// for domain, location := range nginx {
// 	// 	if idx%2 == 0 {
// 	// 		// 中线分割
// 	// 		// s.Marker("arrow", x+idx*intervalX, y+idx*intervalY, 100, 50, style)
// 	// 		// s.Path("M2,2 L2,11 L10,6 L2,2", "fill:blue")
// 	// 		// s.Text(x+idx*intervalX, y+idx*intervalY, domain)

// 	// 	} else {
// 	// 		// s.Marker("arrow", x+idx*intervalX, y-idx*intervalY, 100, 50, style)
// 	// 		// s.Path("M2,2 L2,11 L10,6 L2,2", "fill:blue")
// 	// 		// s.Text(x+idx*intervalX, y-idx*intervalY, domain)
// 	// 	}

// 	// 	s.Line(x+(idx)*x, y+(idx)*y, x+(idx+1)*x, y+(idx+1)*y)
// 	// 	for loc, value := range location {
// 	// 		idx++
// 	// 		s.Text(x+(idx+1)*x, y+(idx+1)*y, loc)
// 	// 		s.Line(x+(idx+1)*x, y+(idx+1)*y, x+(idx+2)*x, y+(idx+2)*y, "fill:none;stroke:black")
// 	// 		s.Text(x+(idx+2)*x, y+(idx+2)*y, value.Dest)
// 	// 	}
// 	// 	idx++
// 	// }

// 	var _x []int
// 	var _y []int

// 	idx := 1

// 	for range nginx {
// 		_x = append(_x, x+idx*intervalX)
// 		if idx%2 == 0 {
// 			_y = append(_y, y+idx*intervalY)
// 		} else {
// 			_y = append(_y, y-idx*intervalY)
// 		}
// 	}

// 	s.Polyline(
// 		_x,
// 		_y,
// 		`fill="none"`,
// 		`stroke="red"`,
// 		`marker-start="url(#dot)"`,
// 		`marker-mid="url(#arrow)"`)
// 	s.End()

// }
