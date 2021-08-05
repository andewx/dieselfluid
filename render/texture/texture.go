package texture

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

const (
	RGBA32 = 32
)

/*
Texture Package loads in Texture files and binds resources
*/

type TexLibrary struct {
	ID       map[string]uint32
	Textures []Texture
	Base     string
}

//Holds GoLang Image Interface Type
type Texture struct {
	Name  string
	Image image.Image
}

func NewTexLibrary() TexLibrary {
	TexLib := TexLibrary{}
	TexLib.Base, _ = os.Getwd()
	fmt.Printf("Texture Library Working Dir: %s\n", TexLib.Base)
	TexLib.Textures = make([]Texture, 1)
	return TexLib
}

/* Load Library should load in any image and decode images into a byte stream
to be fed into the OpenGL Texture Memory*/
func (t *TexLibrary) Load(path string) error {

	f, err := os.Open(t.Base + path) //File is the reader

	if err != nil {
		fmt.Printf("Error opening texture file %s\n", t.Base+path)
		return err
	}
	defer f.Close()

	img, _, err2 := image.Decode(f)
	if err2 != nil {
		fmt.Printf("Error Decoding Image\n")
		return err2
	}

	//Create new "Texture" and append to library
	tex := Texture{path, img}
	t.Textures = append(t.Textures, tex)
	return nil
}

/*Commits Textures into OpenGL Texture Memory Bindings*/
func (t *TexLibrary) Upload() {

	n := len(t.Textures)
	texID := make([]uint32, n)
	gl.GenTextures(int32(n), &texID[0])

	for i := 0; i < n; i++ {

		//Setup bounds and Texture Units
		tex := t.Textures[i]
		t.ID[tex.Name] = texID[i]
		min := tex.Image.Bounds().Min
		max := tex.Image.Bounds().Max
		width := max.X - min.X
		height := max.Y - min.Y
		//Stores images as 32 byte RGBA standard format
		rgbaBuffer := make([]byte, width*height*RGBA32)
		bIndex := 0
		//Construct RAW RGBA Buffer
		for y := min.Y; y < max.Y; y++ {
			for x := min.X; x < max.X; x++ {
				r, g, b, a := tex.Image.At(x, y).RGBA()
				rgbaBuffer[bIndex] = byte(r)
				rgbaBuffer[bIndex+1] = byte(g)
				rgbaBuffer[bIndex+2] = byte(b)
				rgbaBuffer[bIndex+3] = byte(a)
				bIndex += 4
			}
		}
		/*GL COMMANDS*/
		gl.BindTexture(gl.TEXTURE_2D, texID[i])
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&rgbaBuffer[0]))

	}
}
