package drive

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/vivym/igo/internal/session"
)

// Drive provides functions to access iCloud Drive
type Drive struct {
	baseURL string
	http    *resty.Client
	session *session.Session
}

// WebServiceName returns the service name of iCloud Drive
func WebServiceName() string {
	return "drivews"
}

// New creates a new Drive instance
func New(http *resty.Client, session *session.Session) *Drive {
	drive := Drive{
		baseURL: session.AccountInfo.GetWebServiceURL(WebServiceName()),
		http:    http,
		session: session,
	}

	return &drive
}

// Start starts the iCloud Drive service
func (d *Drive) Start() {

}

// Stop stops the iCloud Drive service
func (d *Drive) Stop() {

}

// Test does testing
func (d *Drive) Test(drivewsid string) {
	// d.retrieveItemDetailsInFolder(drivewsid)
	// d.downloadDocument(drivewsid)
	// d.createFolder("FOLDER::com.apple.CloudDocs::CEF6DB51-EF2D-48D6-9608-3EEB31EB83BC", "test_folder")
	item, _ := d.retrieveItemDetailsInFolder(drivewsid)
	item.Name = "haha"
	d.renameItem(item)
}

func (d *Drive) renameItem(item *ItemInfo) (*ItemInfo, error) {
	items, err := d.renameItems([]*ItemInfo{item})
	if items == nil {
		return nil, err
	}
	return items[0], err
}

func (d *Drive) renameItems(items []*ItemInfo) ([]*ItemInfo, error) {
	bodyItems := make([]map[string]string, 0, len(items))
	for _, item := range items {
		bodyItems = append(bodyItems, map[string]string{
			"drivewsid": item.DrivewsID,
			"etag":      item.Etag,
			"name":      item.Name,
			"extension": item.Extension,
		})
	}
	body := map[string][]map[string]string{}
	body["items"] = bodyItems

	rsp, err := d.http.R().
		SetPathParams(map[string]string{
			"clientBuildNumber":     d.session.ClientSetting.ClientBuildNumber,
			"clientMasteringNumber": d.session.ClientSetting.ClientMasteringNumber,
			"clientId":              d.session.ClientSetting.ClientID,
			"dsid":                  d.session.AccountInfo.DsInfo.DsID,
		}).
		SetBody(body).
		Post(d.baseURL + "/renameItems")
	if err != nil {
		return nil, err
	}

	type ResponseBody struct {
		Items []*ItemInfo `json:"items"`
	}
	rspBody := ResponseBody{}
	if err := json.Unmarshal(rsp.Body(), &rspBody); err != nil {
		return nil, err
	}

	return rspBody.Items, nil
}

func (d *Drive) deleteItem(item *ItemInfo) (*ItemInfo, error) {
	items, err := d.deleteItems([]*ItemInfo{item})
	if items == nil {
		return nil, err
	}
	return items[0], err
}

func (d *Drive) deleteItems(items []*ItemInfo) ([]*ItemInfo, error) {
	bodyItems := make([]map[string]string, 0, len(items))
	for _, item := range items {
		bodyItems = append(bodyItems, map[string]string{
			"drivewsid": item.DrivewsID,
			"etag":      item.Etag,
		})
	}
	body := map[string][]map[string]string{}
	body["items"] = bodyItems

	rsp, err := d.http.R().
		SetPathParams(map[string]string{
			"clientBuildNumber":     d.session.ClientSetting.ClientBuildNumber,
			"clientMasteringNumber": d.session.ClientSetting.ClientMasteringNumber,
			"clientId":              d.session.ClientSetting.ClientID,
			"dsid":                  d.session.AccountInfo.DsInfo.DsID,
		}).
		SetBody(body).
		Post(d.baseURL + "/deleteItems")
	if err != nil {
		return nil, err
	}

	type ResponseBody struct {
		Items []*ItemInfo `json:"items"`
	}
	rspBody := ResponseBody{}
	if err := json.Unmarshal(rsp.Body(), &rspBody); err != nil {
		return nil, err
	}

	return rspBody.Items, nil
}

func (d *Drive) createFolder(dstDrivewsID, folder string) (*ItemInfo, error) {
	items, err := d.createFolders(dstDrivewsID, []string{folder})
	if items == nil {
		return nil, err
	}
	return items[0], err
}

func (d *Drive) createFolders(dstDrivewsID string, folders []string) ([]*ItemInfo, error) {
	body := struct {
		DestinationDrivewsID string              `json:"destinationDrivewsId,omitempty"`
		Folders              []map[string]string `json:"folders,omitempty"`
	}{
		DestinationDrivewsID: dstDrivewsID,
	}
	for _, folder := range folders {
		body.Folders = append(body.Folders, map[string]string{
			"clientId": d.session.ClientSetting.ClientID,
			"name":     folder,
		})
	}

	rsp, err := d.http.R().
		SetPathParams(map[string]string{
			"clientBuildNumber":     d.session.ClientSetting.ClientBuildNumber,
			"clientMasteringNumber": d.session.ClientSetting.ClientMasteringNumber,
			"clientId":              d.session.ClientSetting.ClientID,
			"dsid":                  d.session.AccountInfo.DsInfo.DsID,
		}).
		SetBody(body).
		Post(d.baseURL + "/createFolders")
	if err != nil {
		return nil, err
	}

	type ResponseBody struct {
		DestinationDrivewsID string      `json:"destinationDrivewsId,omitempty"`
		Folders              []*ItemInfo `json:"folders"`
	}
	res := ResponseBody{}
	if err := json.Unmarshal(rsp.Body(), &res); err != nil {
		return nil, err
	}

	return res.Folders, nil
}

// TODO: move to docws service
func (d *Drive) downloadDocument(docwsid string) (*DocDownloadInfo, error) {
	baseURL := d.session.AccountInfo.GetWebServiceURL("docws")
	rsp, err := d.http.R().
		SetQueryParams(map[string]string{
			"document_id":           docwsid,
			"token":                 d.session.Auth.SessionToken,
			"clientBuildNumber":     d.session.ClientSetting.ClientBuildNumber,
			"clientMasteringNumber": d.session.ClientSetting.ClientMasteringNumber,
			"clientId":              d.session.ClientSetting.ClientID,
			"dsid":                  d.session.AccountInfo.DsInfo.DsID,
		}).
		Get(baseURL + "/ws/com.apple.CloudDocs/download/by_id")
	if err != nil {
		return nil, err
	}

	info := DocDownloadInfo{}
	if err := json.Unmarshal(rsp.Body(), &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (d *Drive) retrieveItemDetailsInFolder(drivewsid string) (*ItemInfo, error) {
	items, err := d.retrieveItemDetailsInFolders([]string{drivewsid})
	if items == nil {
		return nil, err
	}
	return items[0], err
}

func (d *Drive) retrieveItemDetailsInFolders(drivewsids []string) ([]*ItemInfo, error) {
	type RequestBody struct {
		Drivewsid   string `json:"drivewsid"`
		PartialData bool   `json:"partialData"`
	}

	body := make([]RequestBody, 0, len(drivewsids))
	for _, drivewsid := range drivewsids {
		body = append(body, RequestBody{
			Drivewsid:   drivewsid,
			PartialData: false,
		})
	}

	rsp, err := d.http.R().
		SetPathParams(map[string]string{
			"clientBuildNumber":     d.session.ClientSetting.ClientBuildNumber,
			"clientId":              d.session.ClientSetting.ClientID,
			"clientMasteringNumber": d.session.ClientSetting.ClientMasteringNumber,
			"dsid":                  d.session.AccountInfo.DsInfo.DsID,
		}).
		SetBody(body).
		Post(d.baseURL + "/retrieveItemDetailsInFolders")
	if err != nil {
		return nil, err
	}

	itemDetails := []*ItemInfo{}
	if err := json.Unmarshal(rsp.Body(), &itemDetails); err != nil {
		errRsp := ErrorResponse{}
		if err := json.Unmarshal(rsp.Body(), &errRsp); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("retrieve item details in folders error: %d %s", errRsp.ErrorCode, errRsp.ErrorReason)
	}

	return itemDetails, nil
}
