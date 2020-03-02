package drive

// ItemInfo holds item(file or folder) info
type ItemInfo struct {
	ItemInfoCommon

	ItemInfoFolder
	ItemInfoAppLibrary
	ItemInfoFile
}

// ItemInfoCommon holds item(file or folder) common info
type ItemInfoCommon struct {
	DrivewsID string `json:"drivewsid,omitempty"`
	DocwsID   string `json:"docwsid,omitempty"`
	Zone      string `json:"zone,omitempty"`
	Name      string `json:"name,omitempty"`
	ParentID  string `json:"parentId,omitempty"`
	Etag      string `json:"etag,omitempty"`
	Type      string `json:"type,omitempty"`
}

// ItemInfoFolder holds folder info
type ItemInfoFolder struct {
	AssetQuota          int         `json:"assetQuota,omitempty"`
	FileCount           int         `json:"fileCount,omitempty"`
	ShareCount          int         `josn:"shareCount,omitempty"`
	ShareAliasCount     int         `json:"shareAliasCount,omitempty"`
	DirectChildrenCount int         `json:"directChildrenCount,omitempty"`
	Items               []*ItemInfo `json:"items,omitempty"`
	NumberOfItems       int         `json:"numberOfItems,omitempty"`
}

// ItemInfoAppLibrary holds APP_LIBRARY info
type ItemInfoAppLibrary struct {
	DateCreated         string   `json:"dateCreated,omitempty"`
	MaxDepth            string   `json:"maxDepth,omitempty"`
	Icons               []icon   `json:"icons,omitempty"`
	SupportedExtensions []string `json:"supportedExtensions,omitempty"`
	SupportedTypes      []string `json:"supportedTypes,omitempty"`
}

type icon struct {
	URL  string `json:"url,omitempty"`
	Type string `json:"type,omitempty"`
	Size string `json:"size,omitempty"`
}

// ItemInfoFile holds file info
type ItemInfoFile struct {
	DateModified string `json:"dateModified,omitempty"`
	DateChanged  string `json:"dateChanged,omitempty"`
	Size         int64  `json:"size,omitempty"`
	Extension    string `json:"extension,omitempty"`
}

// ItemInfoUpate holds info for update result
type ItemInfoUpate struct {
	ClientID  string `json:"clientId,omitempty"`
	Status    string `json:"status,omitempty"`
	IsDeleted bool   `json:"isDeleted,omitempty"`
}

// ErrorResponse holds the error response
type ErrorResponse struct {
	ErrorCode   int    `json:"errorCode,omitempty"`
	ErrorReason string `json:"errorReason,omitempty"`
}

// DocDownloadInfo holds doc info for downlading
// TODO: move to docws service
type DocDownloadInfo struct {
	DocumentID string           `json:"document_id,omitempty"`
	DataToken  DocDownloadToken `json:"data_token,omitempty"`
	DoubleEtag string           `json:"double_etag,omitempty"`
}

// DocDownloadToken holds doc token for downlading
// TODO: move to docws service
type DocDownloadToken struct {
	URL                string `json:"url,omitempty"`
	Token              string `json:"token,omitempty"`
	Signature          string `json:"signature,omitempty"`
	WrappingKey        string `json:"wrapping_key,omitempty"`
	ReferenceSignature string `json:"reference_signature,omitempty"`
}
