package static

import "embed"

//go:embed assets/* main.js .vite/*
var EmbedDirStatic embed.FS