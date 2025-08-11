package thumbnail

import (
	"github.com/mahdi-cpp/api-go-pkg/asset"
)

// ThumbnailStorage defines thumbnail persistence operations
type ThumbnailStorage interface {
	SaveThumbnail(assetID, width, height int, data []byte) error
	GetThumbnail(assetID, width, height int) ([]byte, error)
	GetAssetsWithoutThumbnails() ([]int, error)
	GetAsset(assetID int) (*asset.PHAsset, error)
	GetAssetContent(assetID int) ([]byte, error)
}
