package opengl

import (
	"3d-tile-editor/window"
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"io"
	"os"
)

type Model struct {
	Normals, Vecs []mgl32.Vec3
	Uvs           []mgl32.Vec2
	VecIndices, NormalIndices, UvIndices []float32
	Vao uint32
}

// NewModel will read an OBJ model file and create a Model from its contents
func NewModel(file string) *Model {
	// Open the file for reading and check for errors.
	objFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	// Don't forget to close the file reader.
	defer objFile.Close()

	// Create a model to store stuff.
	model := Model{}

	// Read the file and get it's contents.
	for {
		var lineType string

		// Scan the type field.
		_, err := fmt.Fscanf(objFile, "%s", &lineType)

		// Check if it's the end of the file
		// and break out of the loop.
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		// Check the type.
		switch lineType {
		// VERTICES.
		case "v":
			// Create a vec to assign digits to.
			vec := mgl32.Vec3{}

			// Get the digits from the file.
			fmt.Fscanf(objFile, "%f %f %f\n", &vec[0], &vec[1], &vec[2])

			// Add the vector to the model.
			model.Vecs = append(model.Vecs, vec)

		// NORMALS.
		case "vn":
			// Create a vec to assign digits to.
			vec := mgl32.Vec3{}

			// Get the digits from the file.
			fmt.Fscanf(objFile, "%f %f %f\n", &vec[0], &vec[1], &vec[2])

			// Add the vector to the model.
			model.Normals = append(model.Normals, vec)

		// TEXTURE VERTICES.
		case "vt":
			// Create a Uv pair.
			vec := mgl32.Vec2{}

			// Get the digits from the file.
			fmt.Fscanf(objFile, "%f %f\n", &vec[0], &vec[1])

			// Add the uv to the model.
			model.Uvs = append(model.Uvs, vec)

		// INDICES.
		case "f":
			// Create a vec to assign digits to.
			norm := make([]float32, 3)
			vec := make([]float32, 3)
			uv := make([]float32, 3)

			// Get the digits from the file.
			matches, _ := fmt.Fscanf(objFile, "%f/%f/%f %f/%f/%f %f/%f/%f\n", &vec[0], &uv[0], &norm[0], &vec[1], &uv[1], &norm[1], &vec[2], &uv[2], &norm[2])

			if matches != 9 {
				panic("Cannot read your file")
			}

			// Add the numbers to the model.
			model.NormalIndices = append(model.NormalIndices, norm[0])
			model.NormalIndices = append(model.NormalIndices, norm[1])
			model.NormalIndices = append(model.NormalIndices, norm[2])

			model.VecIndices = append(model.VecIndices, vec[0])
			model.VecIndices = append(model.VecIndices, vec[1])
			model.VecIndices = append(model.VecIndices, vec[2])

			model.UvIndices = append(model.UvIndices, uv[0])
			model.UvIndices = append(model.UvIndices, uv[1])
			model.UvIndices = append(model.UvIndices, uv[2])
		}
	}

	model.Vao = model.createVAO()

	// Return the newly created Model.
	return &model
}

// GetRenderableVertices returns a slice of float32s
// formatted in X, Y, Z, U, V. That is, XYZ of the
// vertex and the texture position.
func (model *Model) GetRenderableVertices() []float32 {
	// Create a slice for the outward float32s.
	var out []float32

	// Loop over each vec3 in the indices property.
	for _, position := range model.VecIndices {
		index := int(position) - 1
		vec := model.Vecs[index]
		uv := model.Uvs[int(model.UvIndices[index])-1]

		out = append(out, vec.X(), vec.Y(), vec.Z(), uv.X(), uv.Y())
	}

	// Return the array.
	return out
}

func (model *Model) Draw(program *Program, camera *Camera, window *window.Window) {
	vao := model.Vao

	program.Use()

	// creates perspective camera
	fov := float32(60.0)

	projectionTransform := mgl32.Perspective(mgl32.DegToRad(fov),
		float32(window.Width())/float32(window.Height()),
		0.1,
		10.0)
	camTransform := camera.GetTransform()
	modelTransform := mgl32.Ident4()

	gl.UniformMatrix4fv(program.GetUniformLocation("projection"), 1, false,
		&projectionTransform[0])
	gl.UniformMatrix4fv(program.GetUniformLocation("camera"), 1, false, &camTransform[0])
	gl.UniformMatrix4fv(program.GetUniformLocation("model"), 1, false, &modelTransform[0])

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
	gl.BindVertexArray(0)
}

// Creates the Vertex Array Object
func (model *Model) createVAO() uint32 {
	vertices := model.GetRenderableVertices()

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32;
	gl.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// Copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Size of one whole vertex (sum of attrib sizes)
	var stride int32 = 3*4 + 2*4
	var offset int = 0

	// Position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(0)
	offset += 3*4

	// Texture position
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(1)
	offset += 2*4

	// Bind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return VAO
}