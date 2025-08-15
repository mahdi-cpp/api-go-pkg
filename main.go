package main

import (
	"fmt"
	"github.com/mahdi-cpp/api-go-pkg/collection_manager"
	"github.com/mahdi-cpp/api-go-pkg/metadata"
	"github.com/mahdi-cpp/api-go-pkg/plistcontrol"
	"github.com/mahdi-cpp/api-go-pkg/shared_model"
	"github.com/mahdi-cpp/api-go-pkg/test_model"
	"log"
)

func main() {
	//testCollection2()
	//textCollectionControl()
	//testAccountUserList()
}

func textCollectionControl() {

	MessageManager, err := collection_manager.NewCollectionManager[*test_model.Message]("/media/mahdi/Cloud/Happle/com.helium.messages/chats/7/messages", true)
	if err != nil {
		return
	}

	fmt.Println()

	messages, err := MessageManager.GetAll()
	if err != nil {
		fmt.Printf("Error getting all messages: %v\n", err)
		return
	}

	fmt.Println("All Messages:")
	for _, album := range messages {
		fmt.Printf("  ID: %d, Content: \"%s,              Type: %s \n", album.ID, album.Content, album.Type)
	}

	//album, err := MessageManager.Get(4)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(album.Title)
}

func testInfoPlist() {
	infoPlist := metadata.NewMetadataControl[shared_model.InfoPlist]("/media/mahdi/Cloud/Happle/com.helium.settings/Info.json")
	a, err := infoPlist.Read(true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CFBundleDevelopmentRegion: ", a.CFBundleDevelopmentRegion)
}

func testXml() {

	type InfoPlist struct {
		DevelopmentRegion string `plist:"CFBundleDevelopmentRegion"`
		Executable        string `plist:"CFBundleExecutable"`
		Identifier        string `plist:"CFBundleIdentifier"`
		Name              string `plist:"CFBundleName"`
		Version           string `plist:"CFBundleShortVersionString"`
		Build             string `plist:"CFBundleVersion"`
		CameraUsage       string `plist:"NSCameraUsageDescription"`
		RequiresIPhoneOS  bool   `plist:"LSRequiresIPhoneOS"`
	}

	plistControl := plistcontrol.NewPlistControl[InfoPlist]("/media/mahdi/Cloud/Happle/com.helium.settings/Info.plist")

	// Read existing PLIST
	data, err := plistControl.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("App: %s\nVersion: %s (%s)\n", data.Name, data.Version, data.Build)

	//// Update values
	//data.Build = "2.0"
	//data.Version = "1.1"
	//
	//// Write changes
	//if err := plistControl.Write(data); err != nil {
	//	log.Fatal(err)
	//}
}
