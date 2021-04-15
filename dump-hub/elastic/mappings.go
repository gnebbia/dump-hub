package elastic

const entryMapping = `
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
    "refresh_interval" : "30s"
  },
  "mappings": {
    "dynamic_templates": [
      {
        "all_text": {
          "match_mapping_type": "string",
          "mapping": {
            "copy_to": "_all",
            "type": "text"
          }
        }
      }
    ],
    "properties": {
      "_all": {
        "type": "text"
      }
    }
  }
}
`

const historyMapping = `
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "date": {"type": "keyword" }, 
      "filename": { "type": "keyword" }, 
      "status": { "type": "integer" }
    }
  }
}
`
