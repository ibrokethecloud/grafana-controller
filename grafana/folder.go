package grafana

import (
	"github.com/grafana-tools/sdk"
	"github.com/ibrokethecloud/grafana-controller/api/v1alpha1"
	json "github.com/json-iterator/go"
)

func (a APIClient) GetFolderByName(folder string) (id int, uid string, found bool, err error) {

	folderList, err := a.Client.GetAllFolders(a.Context)
	if err != nil {
		return id, uid, found, err
	}

	for _, foundFolder := range folderList {
		if foundFolder.Title == folder {
			id = foundFolder.ID
			uid = foundFolder.UID
			found = true
		}
	}

	return id, uid, found, err
}

func (a APIClient) SetFolder(fSpec v1alpha1.FolderSpec) (folderStatus v1alpha1.FolderStatus, err error) {
	var f, createdFolder sdk.Folder

	if err := json.Unmarshal([]byte(fSpec.Content), &f); err != nil {
		return folderStatus, err
	}
	// check if folder already exists before we create it //
	id, uid, found, err := a.GetFolderByName(f.Title)

	if found {
		folderStatus.ID = id
		folderStatus.UID = uid
		folderStatus.Message = "Folder existed already"
	} else {
		createdFolder, err = a.CreateFolder(a.Context, f)
		if err != nil {
			return folderStatus, err
		}
		folderStatus.ID = createdFolder.ID
		folderStatus.UID = createdFolder.UID
		folderStatus.Message = ""
	}

	return folderStatus, err
}

func (a APIClient) DeleteFolder(uid string) (err error) {
	_, err = a.Client.DeleteFolderByUID(a.Context, uid)
	return err
}
