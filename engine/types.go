package engine

type Renderable interface {
	Render(render *Renderer)
}

type Command interface {
	Invoke()
}
