package opengl

import (
	"3d-tile-editor/window"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type SimpleShape struct {
	Vertices []float32
	VAO uint32
}

func NewSimpleShape(vertices []float32) *SimpleShape {
	VAO := CreateVAO(vertices)
	return &SimpleShape{Vertices: vertices, VAO: VAO}
}

func (shape *SimpleShape) Draw(program *Program, camera *Camera, window *window.Window) {
	vao := shape.VAO

	program.Use()

	modelTransform := mgl32.Ident4()

	gl.UniformMatrix4fv(program.GetUniformLocation("model"), 1, false, &modelTransform[0])

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
	gl.BindVertexArray(0)
}

