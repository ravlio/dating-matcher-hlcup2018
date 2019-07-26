This project is a https://highloadcup.ru/ru/ contest participant. The main goal is to make search engine for matching people on ephemeral Dating Site, with many index types (roaring bitmaps, skiplists, hashes, rbtrees) with tight cpu and memory consumption optimizations.

The engine reads and loads json data in well optimized manner into memory, makes indexes and starts to recieve http queries. 

Engine resolves filtration, collaborative filtering, grouping with filtering and suggestions. 

Here I preferred performance over code beautifully, so sorry for quality :) I have forked and patched several libraries like gojay, skiplists for better peformance.
