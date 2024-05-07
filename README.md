# Game of Life

Written with Ebitengine.

Check out [my blog post](https://blog.manuelpepe.com/posts/004-game-of-life/) about this.

## Shader Benchmark Results

A shader version of the procedure to generate the next grid is available as [`gol.NextGridShader`](./gol/shader.go)

Here are the results of the benchmark:

```
goos: windows
goarch: amd64
pkg: github.com/manuelpepe/gol/gol
cpu: AMD Ryzen 5 3600X 6-Core Processor
                                           │  bench.txt   │
                                           │    sec/op    │
NextGridShader/window-(H:10-W:10)-12         436.2µ ± 71%
NextGridShader/window-(H:50-W:50)-12         38.22m ±  6%
NextGridShader/window-(H:100-W:100)-12       36.33m ±  6%
NextGridShader/window-(H:1000-W:1000)-12     42.75m ±  6%
NextGridShader/window-(H:10000-W:10000)-12    1.051 ± 23%
NextGrid/window-(H:10-W:10)-12               3.242µ ±  0%
NextGrid/window-(H:50-W:50)-12               90.41µ ±  1%
NextGrid/window-(H:100-W:100)-12             374.2µ ±  2%
NextGrid/window-(H:1000-W:1000)-12           41.63m ±  0%
NextGrid/window-(H:10000-W:10000)-12          4.432 ±  1%
geomean                                      7.483m

                                           │   bench.txt   │
                                           │     B/op      │
NextGridShader/window-(H:10-W:10)-12         3.922Ki ± 26%
NextGridShader/window-(H:50-W:50)-12         21.84Mi ±  0%
NextGridShader/window-(H:100-W:100)-12       21.91Mi ±  0%
NextGridShader/window-(H:1000-W:1000)-12     30.98Mi ±  0%
NextGridShader/window-(H:10000-W:10000)-12   880.7Mi ±  0%
NextGrid/window-(H:10-W:10)-12                 112.0 ±  0%
NextGrid/window-(H:50-W:50)-12               2.625Ki ±  0%
NextGrid/window-(H:100-W:100)-12             10.00Ki ±  0%
NextGrid/window-(H:1000-W:1000)-12           984.0Ki ±  0%
NextGrid/window-(H:10000-W:10000)-12         95.38Mi ±  0%
geomean                                      659.4Ki

                                           │  bench.txt  │
                                           │  allocs/op  │
NextGridShader/window-(H:10-W:10)-12          38.00 ± 3%
NextGridShader/window-(H:50-W:50)-12         7.025k ± 0%
NextGridShader/window-(H:100-W:100)-12       7.027k ± 0%
NextGridShader/window-(H:1000-W:1000)-12     7.337k ± 0%
NextGridShader/window-(H:10000-W:10000)-12   7.427k ± 1%
NextGrid/window-(H:10-W:10)-12                1.000 ± 0%
NextGrid/window-(H:50-W:50)-12                1.000 ± 0%
NextGrid/window-(H:100-W:100)-12              1.000 ± 0%
NextGrid/window-(H:1000-W:1000)-12            1.000 ± 0%
NextGrid/window-(H:10000-W:10000)-12          1.000 ± 0%
geomean                                       50.23
```

At tradeoff between speed and memory usage can start to be seen on grids bigger than 10k x 10k.

## Resources

Compiling: https://ebitengine.org/en/documents/webassembly.html
