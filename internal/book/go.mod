module github.com/brightercommand/Rewind/internal/book

go 1.20


replace github.com/brightercommand/Rewind/internal/pages => ../pages
replace github.com/brightercommand/Rewind/internal/sources => ../sources

require (
	gopkg.in/yaml.v3 v3.0.1
	github.com/gomarkdown/markdown v0.0.0-20230716120725-531d2d74bc12
	github.com/google/uuid v1.3.0
	github.com/brightercommand/Rewind/internal/pages v1.0.0
	github.com/brightercommand/Rewind/internal/pages v1.0.0
)
