# Counter

### Overview
Distributed storage with high availability and persistent storage.
Implemented based on [Conflict-free replicated data type](https://en.wikipedia.org/wiki/Conflict-free_replicated_data_type) algorithm. 
* Data is stored on disk on **binary format**
* Keywords are kept as **uint64** (will help us to implement buckets for increase performance see TODO`s)

**memory layout**    
64bit(key) + 64bit * nodes_number   

**disk layout**   
64bit(key) + 8bit(nodes_number) + 64bit * nodes_number


Each keyword is stored into memory on **64bits**, total occurrences for keyword is stored on **64bits * nodes number**. Other **64bits** are used for 

### Components:
* **storage**: its purpose is to keep data in memory and satisfy outside layers
* **persistence**: dump data from storage on disk periodically at specified interval
* **publicapi**: serve GET/POST requests that get number of occurrences for a keyword or index a text by incrementing number of occurrences for each keyword 
* **syncsrv**: share data with other seeds at specific period of time

### TODO`s
* **Implement buckets** on storage(memory) layer, that`s why we are using key as uint64 at storage(memory) layer 
* **Implement lastUpdateTime** by adding additional property on each keyword on storage(memory) layer that will help to know which key should be sync with other nodes.
* **Redesign Persistence layer** to use lastUpdateTime, store oly diffs, maybe add **compaction** and **SS tables**
* **Implement streaming** on persistence(very easy)

### Not Covered:
*