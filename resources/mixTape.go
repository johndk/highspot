package resources

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
)

type MixTapeApiModel struct {
	Users     []*User     `json:"users"`
	PlayLists []*PlayList `json:"playlists"`
	Songs     []*Song     `json:"songs"`
}

// The storage model is an in-memory key/value store
type MixTapeStorageModel struct {
	userMap        map[string]*User     `json:"-"'`
	songsMap       map[string]*Song     `json:"-"'`
	playListMap    map[string]*PlayList `json:"-"'`
	playListAutoID uint32               `json:"-"'` //0 - 4294967295 (MaxUint32)
}

type MixTape struct {
	MixTapeApiModel
	MixTapeStorageModel
}

//
// UnmarshalJSON is called when the mixtape input JSON file is unmarshalled (deserialized).
// The input data is validated and used to populate the key/value storage model.
//
func (m *MixTape) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &m.MixTapeApiModel)
	if err != nil {
		return err
	}

	err = m.populateStorageModel()
	if err != nil {
		return err
	}

	return nil
}

//
// MarshalJSON is called when the mixtape data is marshalled (serialized) to produce
// the output JSON file.
// The playlist storage model is written to the API model which is then marshalled.
// Only the playlist has changes.
//
func (m *MixTape) MarshalJSON() ([]byte, error) {
	m.PlayLists = make([]*PlayList, 0, len(m.PlayLists))

	for _, val := range m.playListMap {
		m.PlayLists = append(m.PlayLists, val)
	}

	return json.Marshal(&m.MixTapeApiModel)
}

// Add a playlist to the storage model
func (m *MixTape) AddPlayList(playlist *PlayList) error {
	_, ok := m.userMap[playlist.UserID]
	if !ok {
		return errors.New(fmt.Sprintf("User ID %v does not exist.", playlist.UserID))
	}

	for _, songID := range playlist.SongIDs {
		_, ok := m.songsMap[songID]
		if !ok {
			return errors.New(fmt.Sprintf("Song ID %v does not exist.", songID))
		}
	}

	// Auto increment a new ID for the playlist

	if m.playListAutoID < math.MaxUint32 {
		m.playListAutoID++
		playlist.ID = strconv.FormatUint(uint64(m.playListAutoID), 10)
		m.playListMap[playlist.ID] = playlist
	} else {
		return errors.New(fmt.Sprintf("Playlist ID %v exceeds maximum.", math.MaxUint32))
	}

	return nil
}

// Remove a playlist from the storage model
func (m *MixTape) RemovePlayList(playlistID string) error {
	_, ok := m.playListMap[playlistID]
	if !ok {
		return errors.New(fmt.Sprintf("Playlist ID %v does not exist.", playlistID))
	}

	delete(m.playListMap, playlistID)

	return nil
}

// Add a song to a playlist in the storage model
func (m *MixTape) AddSongToPlayList(playlistID, songID string) error {
	_, ok := m.playListMap[playlistID]
	if !ok {
		return errors.New(fmt.Sprintf("Playlist ID %v does not exist.", playlistID))
	}

	_, ok = m.songsMap[songID]
	if !ok {
		return errors.New(fmt.Sprintf("Song ID %v does not exist.", songID))
	}

	m.playListMap[playlistID].SongIDs = append(m.playListMap[playlistID].SongIDs, songID)

	return nil
}

// Validate the input data and populate the storage model
func (m *MixTape) populateStorageModel() error {
	err := m.validateAndPopulateUsers()
	if err != nil {
		return err
	}

	err = m.validateAndPopulateSongs()
	if err != nil {
		return err
	}

	err = m.validateAndPopulatePlayLists()
	if err != nil {
		return err
	}

	return nil
}

func (m *MixTape) validateAndPopulateUsers() error {
	m.userMap = make(map[string]*User)
	for _, user := range m.Users {
		if _, ok := m.userMap[user.ID]; ok {
			return errors.New(fmt.Sprintf("Duplicate user ID %v.", user.ID))
		}

		_, err := strconv.ParseUint(user.ID, 10, 32)
		if err != nil {
			return errors.New(fmt.Sprintf("User ID %v is invalid.", user.ID))
		}

		m.userMap[user.ID] = user
	}

	return nil
}

func (m *MixTape) validateAndPopulateSongs() error {
	m.songsMap = make(map[string]*Song)
	for _, song := range m.Songs {
		if _, ok := m.songsMap[song.ID]; ok {
			return errors.New(fmt.Sprintf("Duplicate song ID %v.", song.ID))
		}

		_, err := strconv.ParseUint(song.ID, 10, 32)
		if err != nil {
			return errors.New(fmt.Sprintf("Song ID %v is invalid.", song.ID))
		}

		m.songsMap[song.ID] = song
	}

	return nil
}

func (m *MixTape) validateAndPopulatePlayLists() error {
	var maxID uint64 = 0
	m.playListMap = make(map[string]*PlayList)
	for _, playlist := range m.PlayLists {
		if _, ok := m.playListMap[playlist.ID]; ok {
			return errors.New(fmt.Sprintf("Duplicate playlist ID %v.", playlist.ID))
		}

		if _, ok := m.userMap[playlist.UserID]; !ok {
			return errors.New(fmt.Sprintf("The user ID %v does not exist.", playlist.UserID))
		}

		for _, songID := range playlist.SongIDs {
			if _, ok := m.songsMap[songID]; !ok {
				return errors.New(fmt.Sprintf("The song ID %v does not exist.", songID))
			}
		}

		id, err := strconv.ParseUint(playlist.ID, 10, 32)
		if err != nil {
			return errors.New(fmt.Sprintf("Playlist ID %v is invalid.", playlist.ID))
		}

		if id > maxID {
			maxID = id
		}

		m.playListMap[playlist.ID] = playlist
	}

	m.playListAutoID = uint32(maxID)

	return nil
}
