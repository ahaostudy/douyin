package dao

import "main/model"

func SaveFileOne(video *model.Video) error {
	err := DB.Create(video).Error
	return err
}
