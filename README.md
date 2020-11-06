## Command Line
```
The Highspot take-home coding exercise.

Usage: highspot [arguments]

The arguments are:

  -c string
        The changes file. (default "changes.json")
  -h    Print the help text.
  -o string
        The output file path. (default "output.json")
  -p string
        The input file path.
  -u string
        The input file URL. (default "https://gist.githubusercontent.com/jmodjeska/0679cf6cd670f76f07f1874ce00daaeb/raw/a4ac53fa86452ac26d706df2e851fb7d02697b4b/mixtape-data.json")
   
```

By default, the input file is downloaded using the URL provided in the take-home exercise. The -p argument can be used to specify a filesystem path to the input file. (The -p argument, when specifed, overrides the -u argument).

The -c argument specifies a filesystem path to the changes file, the default is changes.json.

The -o argument specifies a filesystem path for the output file, the default is output.json

### Examples

To run with the default arguments.

> ./highspot

To display the help text.

> ./highspot -h

To read the input from the filesystem instead of a url.

> ./highspot -p mixtape.json

### Running the Program

The executable 'highspot' in the root directory is for macOS.

1. Download and unzip
2. chmod 777 highspot
3. Allow highspot in System Preferences/Security&Privacy/General

## Changes File

The changes file format is based on JSON patch RFC 6902.

```
[
    {
        "op": "add",
        "path": "/playlists/-",
        "value": {
            "user_id" : "7",
            "song_ids" : [
                "32",
                "40"
            ]
        }
    },
    {
        "op": "remove",
        "path": "/playlists/1"
    },
    {
        "op": "add",
        "path": "/playlists/3/song_ids/-",
        "value":"8"
    }
]

```

The changes file implements the operations specifed in the exercise:

1. The add /playlists/- operation "adds a new playlist" to the end of the playlists collection.
2. The remove /playlists/{id} operation "removes a playlist" with the specifed id from the playlists collection. In the example, the playlist with id 1 is removed.
3. The add /playlists/{id}/song_ids/- operation "adds an existing song to an existing playlist". In the example song id 8 is added to the playlist with id 3.

## Implementation Nodes

The implementation is written in Go. 

Start with the main function in
> cmd/main.go

and the Execute function in
> data/ingester.go

which is called by main.

Ingester Execute implements the three functions specified in the execise: 

1. Ingest the input file.
2. Ingest the changes file.
3. Apply the changes and produce the output fille.

## Scaling Discussion

To handle arbitrarily large data sets of users, songs, and playlists. Imagine a web service that allows these entities to be downloaded in batches of variable size. Consider a different endpoint for each entity. For example the fetch users endpoint:

> http://api.highspot.com/users?streamId=2345321&limit=100&apikey=123abc

```
[
  {
    "stream": {
      id: "2345321",
      action: "put"
    },
    "id": "1",
    "name": "Albin Jaye"
  },
  {
    "stream": {
      id: "2345325",
      action: "put"
    },
    "id": "2",
    "name": "Dipika Crescentia"
  }
  ... more data to specified batch limit=100
]
    
```

The entity data includes a stream id (and action) that is a monotonically increasing number. The stream id is a cursor into the source data. To download the next batch, the largest stream id from the previous batch is used. A uint64 would have a sufficiently large range for the stream id.

The stream metadata could include an action to indicate new, updated, patched, and deleted entities in the stream. 

Clients could pull data from the feed, from an appropriate starting cursor (0 for example), until the cursor reaches the end of the data source, and then periodically check for new data as appropriate.
