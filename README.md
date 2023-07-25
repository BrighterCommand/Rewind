# Rewind
Collates documents from folders to build a multi-version GitBook site

## Build Versioned Documentation

For BrighterCommand to support multiple major versions, it needs to be able to support multiple versions of our documentation.
With each new major version it is likely that something changes in the documentation but not everything changes. So to
avoid "cut and paste" of all documents into the documentation for each version, we want to be able to copy the common
documentation and then the version specific documentation to build the documentation for each version.

## Directory structure

source
- root - the root level files for a GitBook project. Readme, Summary, etc
- shared - the documnetation for Brighter or Darker that is common across all versions
- v1 - the documentation for Brigher or Darker that is specific to V1
- V2 - the documnetation for Brighter or Darket that is specific to v2

docs - the root level files for a GitBook project. Readme, Summary, etc
- v1 - the docs for v1
- v2 - the docs for v2


## Versions copy earlier versions

When we build the documentation for a version, we copy shared, then each earlier version in turn i.e. we build v2 by
copying shared, then v1, then v2. This means that we can change v1 and it will be reflected in all later versions.

## Building TOC in Summary.md

The summary.md file is built by the build script. It is built by reading the summary.md file in the shared folder and
then adding the summary.md file from each version in turn. This means that we can change the TOC for v1 and it will be
reflected in all later versions.



