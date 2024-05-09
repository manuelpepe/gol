# Game of Life

Written with Ebitengine.

Check out [my blog post](https://blog.manuelpepe.com/posts/004-game-of-life/) about this.

## Shaders

Two shader version are available as `gol.NextGridShader` and `gil.NextGridShaderV2` at [./gol/shader.go](./gol/shader.go), with a Keyboard input they can be swapped at runtime using the `H` key.

The first version encodes the grid by setting the Alpha value on each pixel of an image, the shader calculates the live neighbours to each pixel.

The second version encodes the grid in a smaller number of pixels by using the 4 channels of each pixel (RGBA). 


## Resources

Compiling: https://ebitengine.org/en/documents/webassembly.html
