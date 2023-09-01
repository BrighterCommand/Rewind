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

For GitBook we build a [SUMMARY.md](https://docs.gitbook.com/product-tour/git-sync/content-configuration) file.

To build the SUMMARY.md file we use .toc.yaml files that we locate in the root of any document collection. 
The .toc.yaml file is a list of files that we want to include in the table of contents. 

Files are grouped by heading (for example Brighter Configuration, Brighter Request Handlers and Middleware Pipelines). 
This appears in the subsequent markdown file as a level 2 heading.

Under each heading is a collection of links to the pages in that section of the table of contents.

To create this in our .toc.yaml file we create a map for each configuration section. The key is the heading and the value. 
For each link you need to supply: filename, title and indent. The indent is how deeply indented you want the 
file's entry to be in the table of contents. The indent is optional and defaults to 1. This allows you to nest pages.

Assuming the following .toc.yaml file was located in the root of a v9.0.0 folder

```yaml
Brighter Configuration:
- filename: GettingStarted.md
  title: Getting Started
  indent: 1
```

this will yield the following in the SUMMARY.md file

```markdown
## Brighter Configuration
* [Getting Started](/v9.0.0/GettingStarted.md)
```

### Ordering

NOTE: WE NEED TO THINK ABOUT HOW WE WOULD ORDER ENTRIES. THIS MAY ALSO SOLVE THE ERRATIC TEST PROBLEM!!

### Anchor Links

We don't support Anchor Links for page links within the TOC in this version. This is because whilst GitBook displays them 
it does not allow them to be navigated. GitBook does show an outline control that indicates the sub-headings. For this
reason, you should just document whole pages in the TOC.








