package room

import (
	"codebooks/constants"
	"errors"
	"fmt"
	"os"
)

type RoomObject struct {
	ID         string `json:"id"`
	SourceCode string `json:"source_code"`
	Language   string `json:"language"`
}

func CheckRoomExistsObject(roomID string) bool {
	filePath := getFilePathFromRoomID(roomID)
	_, err := os.ReadFile(filePath)

	return err == nil
}

func ReadRoomObject(roomID string) (*RoomObject, error) {
	var room RoomObject
	fmt.Println("read room")

	workingDir, _ := os.Getwd()
	var filePath string = workingDir + constants.PathSeparatorStr + "data" + constants.PathSeparatorStr + roomID
	fmt.Println(filePath)

	// read the room from DB/file
	b, err := os.ReadFile(filePath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("file not found")
		return createDefaultRoom(roomID, filePath)
	} else if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		room = RoomObject{
			SourceCode: string(b),
			ID:         roomID,
			Language:   "c++",
		}
		return &room, nil
	}
}

func SaveRoomObject(room RoomObject) error {
	// save the room's code to DB/file
	workingDir, _ := os.Getwd()
	var filePath string = workingDir + constants.PathSeparatorStr + "data" + constants.PathSeparatorStr + room.ID
	err := os.WriteFile(filePath, []byte(room.SourceCode), 0644)
	if err != nil {
		return err
	}

	return nil
}

func createDefaultRoom(roomID, filePath string) (*RoomObject, error) {
	err := os.WriteFile(filePath, []byte(constants.DefaultSourceCode), 0644)
	if err != nil {
		return nil, err
	}

	return &RoomObject{
		SourceCode: constants.DefaultSourceCode,
		ID:         roomID,
		Language:   "c++",
	}, nil
}

func getFilePathFromRoomID(roomID string) string {
	workingDir, _ := os.Getwd()
	return workingDir + constants.PathSeparatorStr + "data" + constants.PathSeparatorStr + roomID
}
