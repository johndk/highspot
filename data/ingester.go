package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"highspot/data/validation"
	"highspot/resources"
	"log"
	"regexp"
)

type Ingester struct {
	inputReader   Reader
	changesReader Reader
	outputWriter  Writer
}

func NewIngestor(inputReader Reader, changesReader Reader, outputWriter Writer) *Ingester {
	ingestor := Ingester{
		inputReader:   inputReader,
		changesReader: changesReader,
		outputWriter:  outputWriter,
	}
	return &ingestor
}

//
// For this exercise, you will write 3 functions for a command-line batch application.
// The three functions are ingestInput, ingestChanges, produceOutput
//
func (i *Ingester) Execute() error {
	//
	// Ingest an input JSON file which we will provide, mixtape.json.
	//
	mixtape, err := i.ingestInput()
	if err != nil {
		return errors.New(fmt.Sprintf("Ingest input failed. %v", err))
	}

	//
	// Ingest a changes file which you will create.
	//
	changes, err := i.ingestChanges()
	if err != nil {
		return errors.New(fmt.Sprintf("Ingest changes failed. %v", err))
	}

	//
	// Produce output.json which must have the same structure as the mixtape.json input.
	//
	err = i.produceOutput(mixtape, changes)
	if err != nil {
		return errors.New(fmt.Sprintf("Produce output failed. %v", err))
	}

	return nil
}

//
// Ingest and validate the input file.
//
func (i *Ingester) ingestInput() (*resources.MixTape, error) {
	//
	// Read and validate the input json document
	//

	data, err := i.readInput()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot read input file. %v.", err))
	}

	err = validation.Validate(validation.InputSchema, string(data))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid input file. %v", err))
	}

	//
	// Unmarshal (deserialize) input json document and validate
	//

	var mixtape resources.MixTape
	err = json.Unmarshal(data, &mixtape)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid input file. %v", err))
	}

	return &mixtape, nil
}

//
// Ingest and validate the changes file
//
func (i *Ingester) ingestChanges() ([]resources.Change, error) {
	//
	// Read and validate the input json document
	//

	data, err := i.readChanges()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot read changes file. %v", err))
	}

	err = validation.Validate(validation.PatchSchema, string(data))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid changes file. %v", err))
	}

	//
	// Unmarshal (deserialize) input json document and validate
	//

	var changes []resources.Change
	err = json.Unmarshal(data, &changes)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid changes file. %v", err))
	}

	return changes, nil
}

//
// Apply the changes and generate the output file
//
func (i *Ingester) produceOutput(mixtape *resources.MixTape, changes []resources.Change) error {
	//
	// Apply the changes
	//
	err := i.applyChanges(mixtape, changes)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot apply changes. %v", err))
	}

	//
	// Write the output file
	//
	data, err := json.Marshal(mixtape)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot write output file. %v", err))
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, data, "", "  ")
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot write output file. %v", err))
	}

	err = i.writeOutput(prettyJSON.Bytes())
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot write output file. %v", err))
	}

	return nil
}

func (i *Ingester) applyChanges(mixtape *resources.MixTape, changes []resources.Change) error {
	//
	// Loop over the changes and apply each change to the mixtape data model
	//
	for _, change := range changes {
		if change.Op == "add" {
			//
			// Add a new playlist; the playlist should contain at least one song.
			//

			ok, err := regexp.MatchString("/playlists/-", change.Path)
			if err != nil {
				return err
			}
			if ok {
				err = i.applyAddPlaylist(mixtape, &change)
				if err != nil {
					log.Printf("Skipping add playlist. %v", err)
				}
				continue
			}

			//
			// Add an existing song to an existing playlist
			//

			ok, err = regexp.MatchString("/playlists/[0-9]+/song_ids/-", change.Path)
			if err != nil {
				return err
			}
			if ok {
				err = i.applyAddSongToPlaylist(mixtape, &change)
				if err != nil {
					log.Printf("Skipping add song to playlist. %v", err)
				}
			}
		} else if change.Op == "remove" {
			//
			// Remove a playlist.
			//

			ok, err := regexp.MatchString("/playlists/[0-9]+", change.Path)
			if err != nil {
				return err
			}
			if ok {
				err = i.applyRemovePlaylist(mixtape, &change)
				if err != nil {
					log.Printf("Skipping remove playlist. %v", err)
				}
			}
		}
	}

	return nil
}

func (i *Ingester) applyAddPlaylist(mixtape *resources.MixTape, change *resources.Change) error {
	if change.Value == nil {
		return errors.New("Missing playlist value.")
	}

	playlistJSON, err := json.Marshal(change.Value)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid playlist value. %v", err))
	}

	err = validation.Validate(validation.PatchPlaylistSchema, string(playlistJSON))
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid playlist value. %v", err))
	}

	var playlist resources.PlayList
	err = json.Unmarshal(playlistJSON, &playlist)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid playlist value. %v", err))
	}

	return mixtape.AddPlayList(&playlist)
}

func (i *Ingester) applyRemovePlaylist(mixtape *resources.MixTape, change *resources.Change) error {
	re := regexp.MustCompile("/playlists/([0-9]+)")
	match := re.FindStringSubmatch(change.Path)
	return mixtape.RemovePlayList(match[1])
}

func (i *Ingester) applyAddSongToPlaylist(mixtape *resources.MixTape, change *resources.Change) error {
	re := regexp.MustCompile("/playlists/([0-9]+)/song_ids/-")
	if change.Value == nil {
		return errors.New("Missing song ID value.")
	}
	songID, ok := change.Value.(string)
	if !ok {
		return errors.New("Invalid song ID value.")
	}
	match := re.FindStringSubmatch(change.Path)
	return mixtape.AddSongToPlayList(match[1], songID)
}

func (i *Ingester) readInput() ([]byte, error) {
	data, err := i.inputReader.Read()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (i *Ingester) readChanges() ([]byte, error) {
	data, err := i.changesReader.Read()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (i *Ingester) writeOutput(data []byte) error {
	return i.outputWriter.Write(data)
}
