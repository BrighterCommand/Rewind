# Rewind
Collates documents from folders to build a multi-version GitBook site

## Versioned Documentation

For BrighterCommand to support multiple major versions, it needs to be able to support multiple versions of our documentation.
With each new major version it is likely that something changes in the documentation but not everything changes. We need GitBook to
offer at least two versions- current and last. But GitBook does not allow multiple links to the same file within the TOC. This
means that each version must have its own copy of any file in the TOC. If both 9.0.0 and 10.0.0 have the file Introduction.md, then
each version's TOC must point to a separate copy of the file. So to avoid "cut and paste" of all documents into the documentation 
for each version, we want to be able to copy the common files to each version and then the version specific documentation 
to use with the TOC for each version.

## Directory structure

We take a directory structure that gives us the shared and versioned documents, and how to fit them into a TOC that covers both

source
- root - the root level files for a GitBook project. Readme, Summary, etc
- shared - the documentation for Brighter or Darker that is common across all versions. .toc.yaml and .md files
- 9.0.0 - the documentation for Brighter or Darker that is specific to V9. .toc.yaml and .md files
  - _static\images - any images used by the docs
- 10.0.0 - the documentation for Brighter or Darker that is specific to v10. .toc.yaml and .md files
  - _static\images - any images used by the docs
 
We then merge shared with each version and build the Summary.md file which from the .toc.yaml files

docs - the root level files for a GitBook project. Readme, Summary, etc
- 9.0.0 - the docs for v1
  - _static\images - any images used by the docs
- 10.0.0 - the docs for v2
  - _static\images - any images used by the docs

## Building the documentation

When we build the documentation for a version, we copy shared, then the documentation for each version. This means that 
if you have shared and 9.0.0 and 10.0.0, any docs that appear in both 9.0.0 and 10.0.0 should appear in shared. We won't copy 
from 9.0.0 to 10.0.0. 

This works to support two versions => the tip of the spear and the last release. If we want to support documentation for three versions
this would get more complicated as we would need to flow documents from one version to the next i.e. 8.0.0 to 9.0.0 to 10.0.0 unless 
overwritten. This is out-of-scope as we have no plans to support three versions right now.

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








