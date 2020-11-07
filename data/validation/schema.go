package validation

var PatchSchema = `{
    "type": "array",
    "items": {
        "type": "object",
        "properties": {
            "op": {
                "type": "string",
                "enum": [
                    "add",
                    "remove"
                ]
            },
            "path": {
                "type": "string",
                "maxLength": 32,
                "pattern": "^(/playlists/-|/playlists/[0-9]+|/playlists/[0-9]+/song_ids/-)$"
            },
            "value": {}
        },
        "additionalProperties": false,
        "required": [
            "op",
            "path"
        ]
    }
}`

var PatchPlaylistSchema = `{
    "type": "object",
    "properties": {
        "id": {
            "type": "string",
            "minLength": 1,
            "maxLength": 10
        },
        "user_id": {
            "type": "string",
            "minLength": 1,
            "maxLength": 10
        },
        "song_ids": {
            "type": "array",
            "items": {
                "type": "string",
                "minLength": 1,
                "maxLength": 10
            },
            "minItems": 1,
            "maxItems": 512,
            "uniqueItems": true,
            "default": []
        }
    },
    "additionalProperties": false,
    "required": [
        "id",
        "user_id",
        "song_ids"
    ]
}`

var InputSchema = `{
    "definitions": {
        "user": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
					"minLength": 1,
                    "maxLength": 10
                },
                "name": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 512
                }
            },
  			"additionalProperties": false,
            "required": [
                "id",
                "name"
            ]
        },
        "playlist": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
					"minLength": 1,
                    "maxLength": 10
                },
                "user_id": {
                    "type": "string",
					"minLength": 1,
                    "maxLength": 10
                },
                "song_ids": {
                    "type": "array",
                    "items": {
                        "type": "string",
						"minLength": 1,
                    	"maxLength": 10
                    },
 					"minItems": 1,
  					"maxItems": 512,
					"uniqueItems": true,
                    "default": []
                }
            },
			"additionalProperties": false,
            "required": [
                "id",
                "user_id",
                "song_ids"
            ]
        },
        "song": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
					"minLength": 1,
                    "maxLength": 10
                },
                "artist": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 512
                },
                "title": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 512
                }
            },
			"additionalProperties": false,
            "required": [
                "id",
                "artist",
                "title"
            ]
        }
    },
    "type": "object",
    "properties": {
        "users": {
            "type": "array",
            "items": {
                "$ref": "#/definitions/user"
            },
            "default": []
        },
        "playlists": {
            "type": "array",
            "items": {
                "$ref": "#/definitions/playlist"
            },
            "default": []
        },
        "songs": {
            "type": "array",
            "items": {
                "$ref": "#/definitions/song"
            },
            "default": []
        }
    },
    "additionalProperties": false,
    "required": [
        "users",
        "playlists",
        "songs"
    ]
}`
