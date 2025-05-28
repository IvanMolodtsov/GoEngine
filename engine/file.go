package engine

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ReadFile(path string) Mesh {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var mesh Mesh

	var vertices = make([]*Vector3d, 0)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		text := scanner.Text()

		tokens := strings.Split(text, ` `)

		switch tokens[0] {
		case "v":
			var x, y, z float64
			x, err = strconv.ParseFloat(tokens[1], 64)
			if err != nil {
				panic(err)
			}
			y, err = strconv.ParseFloat(tokens[2], 64)
			if err != nil {
				panic(err)
			}
			z, err = strconv.ParseFloat(tokens[3], 64)
			if err != nil {
				panic(err)
			}
			vertices = append(vertices, NewVector3d(x, y, z))
		case "f":
			var a, b, c int64
			a, err = strconv.ParseInt(tokens[1], 10, 64)
			if err != nil {
				panic(err)
			}
			b, err = strconv.ParseInt(tokens[2], 10, 64)
			if err != nil {
				panic(err)
			}
			c, err = strconv.ParseInt(tokens[3], 10, 64)
			if err != nil {
				panic(err)
			}

			mesh.Tris = append(mesh.Tris, &Triangle{
				Points: [3]*Vector3d{
					vertices[a-1], vertices[b-1], vertices[c-1],
				},
			})
		default:
			//junk
		}
	}

	return mesh
}
