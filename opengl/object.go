package opengl

import (
	"3d-tile-editor/window"
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"io"
	"os"
)

type Object struct {
	vertices, normals []mgl32.Vec3
	faces, uvs []mgl32.Vec2
	vao uint32
	vecIndices, normalIndices, uvIndices []float32
}

func NewObject(file string) *Object {
	objFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	defer objFile.Close()

	obj := Object{}

	for {
		var lineType string

		_, err := fmt.Fscanf(objFile, "%s", &lineType)

		// Check if it's the end of the file
		// and break out of the loop.
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		switch lineType {
		// VERTICES.
		case "v":
			vec := mgl32.Vec3{}
			fmt.Fscanf(objFile, "%f %f %f\n", &vec[0], &vec[1], &vec[2])
			obj.vertices = append(obj.vertices, vec)
		// NORMALS.
		case "vn":
			vec := mgl32.Vec3{}
			fmt.Fscanf(objFile, "%f %f %f\n", &vec[0], &vec[1], &vec[2])
			obj.normals = append(obj.normals, vec)
		// UVS.
		case "vt":
			vec := mgl32.Vec2{}
			fmt.Fscanf(objFile, "%f %f\n", &vec[0], &vec[1])
			obj.uvs = append(obj.uvs, vec)
		// INDICES.
		case "f":
			norm := make([]float32, 3)
			vec := make([]float32, 3)
			uv := make([]float32, 3)
			matches, _ := fmt.Fscanf(objFile, "%f/%f/%f %f/%f/%f %f/%f/%f\n", &vec[0], &uv[0], &norm[0], &vec[1], &uv[1], &norm[1], &vec[2], &uv[2], &norm[2])
			if matches != 9 {
				panic("Cannot read your file")
			}

			obj.normalIndices = append(obj.normalIndices, norm[0])
			obj.normalIndices = append(obj.normalIndices, norm[1])
			obj.normalIndices = append(obj.normalIndices, norm[2])

			obj.vecIndices = append(obj.vecIndices, vec[0])
			obj.vecIndices = append(obj.vecIndices, vec[1])
			obj.vecIndices = append(obj.vecIndices, vec[2])

			obj.uvIndices = append(obj.uvIndices, uv[0])
			obj.uvIndices = append(obj.uvIndices, uv[1])
			obj.uvIndices = append(obj.uvIndices, uv[2])
		}
	}

	obj.vao = CreateVAO(obj.getRenderableVertices())

	return &obj
}

func (obj *Object) getRenderableVertices() []float32 {
	var out []float32

	for _, position := range obj.vecIndices {
		index := int(position) - 1
		vec := obj.vertices[index]
		uv := obj.uvs[int(obj.uvIndices[index])-1]

		out = append(out, vec.X(), vec.Y(), vec.Z(), uv.X(), uv.Y())
	}

	return out
}

func (obj *Object) Draw(program *Program, camera *Camera, window *window.Window) {
	vao := obj.vao

	program.Use()

	modelTransform := mgl32.Ident4()

	gl.UniformMatrix4fv(program.GetUniformLocation("model"), 1, false, &modelTransform[0])

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(obj.vecIndices)))
	gl.BindVertexArray(0)
}