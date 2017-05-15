# Coverage
## Proof-of-concept to display coverage info, starting with epa.gov data

This project takes lists of urls and associated archiving information, and turns that into a tree of url paths with associated coverage information. The output of this program is cached in `cache.json`. This resulting cache file is very large, so we provide a server that delivers incremental pieces of the cached tree as a web server.

### Current Coverage Sources
Actual source datasets can be found in the `/coverage` directory. It currently includes the following:

* Archivers 2
* EDGI Nomination Tool Uncrawlables
* The Internet Archive
* Project Svalbard json-ld crawl