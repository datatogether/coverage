# Coverage
## Proof-of-concept to display coverage info, starting with epa.gov data

This project takes lists of urls and associated archiving information, and turns that into a tree of url paths with associated coverage information. The output of this program is cached in `cache.json`. This resulting cache file is very large, so we provide a server that delivers incremental pieces of the cached tree as a web server. If you're uninterested in dynamically calculating coverage completion, feel free to work with the `cache.json` file.

### Current Coverage Sources
Actual source datasets can be found in the `/repositoryies` directory. It currently includes the following:

* Archivers 2
* archivers.space
* EDGI Nomination Tool Uncrawlables
* The Internet Archive
* Project Svalbard json-ld crawl