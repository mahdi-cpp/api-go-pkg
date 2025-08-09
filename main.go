package main

import (
	"fmt"
	"github.com/mahdi-cpp/api-go-pkg/collection"
	"github.com/mahdi-cpp/api-go-pkg/metadata"
	"github.com/mahdi-cpp/api-go-pkg/network"
	"github.com/mahdi-cpp/api-go-pkg/plistcontrol"
	"github.com/mahdi-cpp/api-go-pkg/shared_model"
	"log"
)

func main() {
	testCollection2()
}

func testCollection2() {
	users, err := collection.NewCollectionManager[*shared_model.User]("albums_test.json", false)
	if err != nil {
		fmt.Println("UserStorage:", err)
		return
	}

	item := &shared_model.User{FirstName: "Original"}
	create, err := users.Create(item)
	if err != nil {
		return
	}

	fmt.Println(create.FirstName)

	//fmt.Println(createdItem.ID))
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

func testNetwork() {

	userControl := network.NewNetworkControl[[]shared_model.User]("http://localhost:8080/api/v1/user/")

	// Make request (nil body if not needed)
	users, err := userControl.Read("list", nil)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Use the data
	for _, user := range *users {
		fmt.Printf("%d: %s (%s %s)\n",
			user.ID,
			user.Username,
			user.FirstName,
			user.LastName)
	}
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
