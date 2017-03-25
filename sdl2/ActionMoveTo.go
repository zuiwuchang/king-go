package sdl2

//移動到 指定 坐標
type ActionMoveTo struct {
	//目的坐標
	x, y float64
	//每秒移動
	speed float64
}
