package resources

import (
	"encoding/json"
	"errors"
	"fmt"
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
	return m.validateAndAddPlaylist(playlist)
}

// Remove a playlist from the storage model
func (m *MixTape) RemovePlayList(playlistID string) error {
	_, err := strconv.ParseUint(playlistID, 10, 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Playlist ID %v is invalid.", playlistID))
	}

	_, ok := m.playListMap[playlistID]
	if !ok {
		return errors.New(fmt.Sprintf("Playlist ID %v does not exist.", playlistID))
	}

	delete(m.playListMap, playlistID)

	return nil
}

// Add a song to a playlist in the storage model
func (m *MixTape) AddSongToPlayList(playlistID, songID string) error {
	_, err := strconv.ParseUint(playlistID, 10, 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Playlist ID %v is invalid.", playlistID))
	}

	_, ok := m.playListMap[playlistID]
	if !ok {
		return errors.New(fmt.Sprintf("Playlist ID %v does not exist.", playlistID))
	}

	_, err = strconv.ParseUint(songID, 10, 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Song ID %v is invalid.", songID))
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
	err := m.validateAndAddUsers()
	if err != nil {
		return err
	}

	err = m.validateAndAddSongs()
	if err != nil {
		return err
	}

	err = m.validateAndAddPlayLists()
	if err != nil {
		return err
	}

	return nil
}

func (m *MixTape) validateAndAddUsers() error {
	m.userMap = make(map[string]*User)
	for _, user := range m.Users {
		_, err := strconv.ParseUint(user.ID, 10, 32)
		if err != nil {
			return errors.New(fmt.Sprintf("User ID %v is invalid.", user.ID))
		}

		if _, ok := m.userMap[user.ID]; ok {
			return errors.New(fmt.Sprintf("Duplicate user ID %v.", user.ID))
		}

		m.userMap[user.ID] = user
	}

	return nil
}

func (m *MixTape) validateAndAddSongs() error {
	m.songsMap = make(map[string]*Song)
	for _, song := range m.Songs {
		_, err := strconv.ParseUint(song.ID, 10, 32)
		if err != nil {
			return errors.New(fmt.Sprintf("Song ID %v is invalid.", song.ID))
		}

		if _, ok := m.songsMap[song.ID]; ok {
			return errors.New(fmt.Sprintf("Duplicate song ID %v.", song.ID))
		}

		m.songsMap[song.ID] = song
	}

	return nil
}

func (m *MixTape) validateAndAddPlayLists() error {
	m.playListMap = make(map[string]*PlayList)
	for _, playlist := range m.PlayLists {
		m.validateAndAddPlaylist(playlist)
	}

	return nil
}

func (m *MixTape) validateAndAddPlaylist(playlist *PlayList) error {
	_, err := strconv.ParseUint(playlist.ID, 10, 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Playlist ID %v is invalid.", playlist.ID))
	}

	if _, ok := m.playListMap[playlist.ID]; ok {
		return errors.New(fmt.Sprintf("Duplicate playlist ID %v.", playlist.ID))
	}

	_, err = strconv.ParseUint(playlist.UserID, 10, 32)
	if err != nil {
		return errors.New(fmt.Sprintf("User ID %v is invalid.", playlist.UserID))
	}

	if _, ok := m.userMap[playlist.UserID]; !ok {
		return errors.New(fmt.Sprintf("The user ID %v does not exist.", playlist.UserID))
	}

	for _, songID := range playlist.SongIDs {
		_, err := strconv.ParseUint(songID, 10, 32)
		if err != nil {
			return errors.New(fmt.Sprintf("Song ID %v is invalid.", songID))
		}

		if _, ok := m.songsMap[songID]; !ok {
			return errors.New(fmt.Sprintf("The song ID %v does not exist.", songID))
		}
	}

	m.playListMap[playlist.ID] = playlist

	return nil
}
