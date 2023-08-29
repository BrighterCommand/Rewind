# Rewind
Collates documents from folders to build a multi-version GitBook site

## Versioned Documentation

For BrighterCommand to support multiple major versions, it needs to be able to support multiple versions of our documentation.
With each new major version it is likely that something changes in the documentation but not everything changes. So to
avoid "cut and paste" of all documents into the documentation for each version, we want to be able to copy the common
documentation and then the version specific documentation to build the documentation for each version.

## Directory structure

source
- root - the root level files for a GitBook project. Readme, Summary, etc
- shared - the documentation for Brighter or Darker that is common across all versions
- v1 - the documentation for Brighter or Darker that is specific to V1
- V2 - the documentation for Brighter or Darker that is specific to v2

docs - the root level files for a GitBook project. Readme, Summary, etc
- v1 - the docs for v1
- v2 - the docs for v2

## Building the documentation

When we build the documentation for a version, we copy shared, then the documentation for each version. This means that 
if you have shared and V1 and V2, you will need to include docs from V1 that should appear in V2, within V2. In a later
version we may change this so that we include V1 within V2, but for now, we are keeping it simple.

## Table of Contents

The table of contents lives in the SUMMARY.md file. To build it, we use the .toc.yaml files that we locate in the root
of any document collection. The .toc.yaml file is a list of files that we want to include in the table of contents. 

Files are grouped by heading. This appears in the subsequent markdown file as a level 2 heading.

Under each heading is a collection of files.

For each file you need to record the filename, the title and the indent. The indent is how deeply indented you want the 
file's entry to be in the table of contents. The indent is optional and defaults to 1.

Assuming the following .toc.yaml file was located in the root of a v9.0.0 folder

```yaml
Brighter Configuration:
- filename: GettingStarted.md
  title: Getting Started
  indent: 1
```

will yield

```markdown
## Brighter Configuration
* [Getting Started](/v9.0.0/GettingStarted.md)
```







