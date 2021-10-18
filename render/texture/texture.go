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
Texture Package loads in Texture files and binds resources. TexLibrary maintains
list of texture named resource ids and associates with GLTextureID
*/

type TexLibrary struct {
	ID        map[string]uint32 //Maps Texture.Names to their GLInt Gen Texture IDS
	Textures  []Texture
	Base      string
	HasDevice bool
}

//Holds GoLang Image Interface Type
type Texture struct {
	Name           string
	Image          image.Image
	TexGLUnit      int32
	TexUVCoordAttr int32
	Sampler        TexSampler
}

//TexSampler holds GL / Driver sampler constants
type TexSampler struct {
	MagFilter int32
	MinFilter int32
	WrapS     int32
	WrapT     int32
}

func NewTexLibrary() TexLibrary {
	TexLib := TexLibrary{}
	TexLib.Base, _ = os.Getwd()
	TexLib.HasDevice = false
	fmt.Printf("Texture Library Working Dir:\n%s\n", TexLib.Base)
	return TexLib
}

/* Direcly loads texture into the texture library, does not interface with GL/GPU operations*/
func (t *TexLibrary) Load(path string, texIndex int32) error {

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
	tex := Texture{path, img, gl.TEXTURE0 + texIndex, 0, TexSampler{}}
	if t.Textures == nil {
		t.Textures = make([]Texture, 1)
		t.Textures[0] = tex
	} else {
		t.Textures = append(t.Textures, tex)
	}
	return nil
}

/*Commits Textures into OpenGL Texture Memory Bindings*/
func (t *TexLibrary) CommitTexLibGL() {

	if t.HasDevice {

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
		gl.BindTexture(gl.TEXTURE_2D, 0)
	} else {
		fmt.Printf("FUNC:(TexLibrary).CommitTexLibGL() NOP -- No Device Context Set\n")
	}
}

//RemoveTexLibGL removes entire Texture Library Binding from gl texture memory
func (t *TexLibrary) RemoveTexLibGL() {
	if t.HasDevice {
		//Construct GL ID Tag Read Buffer
		texIDS := make([]uint32, len(t.ID))
		index := 0
		for texName := range t.ID {
			texIDS[index] = t.ID[texName] //GLUINT32 TEXTURE HARDWARE ID
		}
		n := int32(len(texIDS))
		gl.DeleteTextures(n, &texIDS[0])
	} else {
		fmt.Printf("FUNC:(TexLibrary).RemoveTexLibGL() NOP -- No Device Context Set\n")
	}
}
