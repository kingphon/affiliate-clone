package service

import (
	natsmodel "git.selly.red/Selly-Modules/natsio/model"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/file"
)

// FileResponseInterface ...
type FileResponseInterface interface {
	ConvertResponseFilePhoto(f *mgaffiliate.FilePhoto) *file.FilePhoto

	ConvertResponseFilePhotoNats(f *file.FilePhoto) *natsmodel.FilePhoto

	ConvertResponseListFilePhoto(f []*mgaffiliate.FilePhoto) []*file.FilePhoto
}

// fileResponseImplement ...
type fileResponseImplement struct{}

// ConvertResponseFilePhotoNats ...
func (s fileResponseImplement) ConvertResponseFilePhotoNats(f *file.FilePhoto) *natsmodel.FilePhoto {
	if f == nil {
		return nil
	}

	return &natsmodel.FilePhoto{
		ID:   f.ID,
		Name: f.Name,
		Dimensions: &natsmodel.FileDimensions{
			Small: &natsmodel.FileSize{
				Width:  f.Dimensions.Small.Width,
				Height: f.Dimensions.Small.Height,
			},
			Medium: &natsmodel.FileSize{
				Width:  f.Dimensions.Medium.Width,
				Height: f.Dimensions.Medium.Height,
			},
		},
	}
}

// FileResponse ...
func FileResponse() FileResponseInterface {
	return &fileResponseImplement{}
}

// ConvertResponseFilePhoto ...
func (fileResponseImplement) ConvertResponseFilePhoto(f *mgaffiliate.FilePhoto) *file.FilePhoto {
	if f == nil {
		return nil
	}

	return &file.FilePhoto{
		ID:   f.ID.Hex(),
		Name: f.Name,
		Dimensions: &file.FileDimensions{
			Small: &file.FileSize{
				Width:  f.Dimensions.Small.Width,
				Height: f.Dimensions.Small.Height,
			},
			Medium: &file.FileSize{
				Width:  f.Dimensions.Medium.Width,
				Height: f.Dimensions.Medium.Height,
			},
		},
	}
}

// ConvertResponseListFilePhoto ...
func (s fileResponseImplement) ConvertResponseListFilePhoto(f []*mgaffiliate.FilePhoto) []*file.FilePhoto {
	var result = make([]*file.FilePhoto, 0)
	if len(f) == 0 {
		return result
	}

	for _, photo := range f {
		fs := s.ConvertResponseFilePhoto(photo).GetResponseData()
		if fs != nil {
			result = append(result, fs)
		}
	}
	return result
}
