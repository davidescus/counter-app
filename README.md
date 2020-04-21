# Counter

### Overview
Distributed storage with high availability and persistent storage.
Implemented based on [Conflict-free replicated data type](https://en.wikipedia.org/wiki/Conflict-free_replicated_data_type) algorithm. 
* Data is stored on disk on **binary format**
* Keywords are kept as **uint64** using **crc64**  (will help us to implement buckets for increase performance see TODO`s). Collision probabilities can be founded [here](https://en.wikipedia.org/wiki/Birthday_problem#Probability_table)
* Each node has his **specific ID** starting from 0.
* Each node holds entire collection of data offering **High Availability**

### Components:
* **storage**: its purpose is to keep data in memory and satisfy outside layers
* **persistence**: dump data from storage on disk periodically at specified interval
* **publicapi**: serve GET/POST requests that get number of occurrences for a keyword or index a text by incrementing number of occurrences for each keyword 
* **syncsrv**: share data with other seeds at specific period of time

### Run it:
* clone repo
* add env vars on each node (see bellow)
* go run cmd/main.go (or build it or install it)
* access node:port/swagger, you can get from there a swagger definition as json

#### Node prepare
This is an example of 3 nodes configuration
```bash
# node1
export COUNTER_NODE_ID=0
export COUNTER_SEEDS=http://localhost:4001,http://localhost:5001
export COUNTER_PUBLIC_PORT=3000
export COUNTER_SYNC_PORT=3001

#node2
export COUNTER_NODE_ID=1
export COUNTER_SEEDS=http://localhost:3001,http://localhost:5001
export COUNTER_PUBLIC_PORT=4000
export COUNTER_SYNC_PORT=4001

#node 3
export COUNTER_NODE_ID=2
export COUNTER_SEEDS=http://localhost:3001,http://localhost:4001
export COUNTER_PUBLIC_PORT=5000
export COUNTER_SYNC_PORT=5001
```

### TODO`s
* **Implement buckets** on storage(memory) layer, that`s why we are using key as uint64 at storage(memory) layer 
* **Implement lastUpdateTime** by adding additional property on each keyword on storage(memory) layer that will help to know which key should be sync with other nodes.
* **Redesign Persistence layer** to use lastUpdateTime, store oly diffs, maybe add **compaction** and **SS tables**
* **Implement streaming** on persistence(very easy)

### Not Covered:
* Add nodes (easy to implement)
* Remove nodes