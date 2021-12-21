package stream

import (
	"github.com/asumsi/api.inlive/internal/models/db"
)

func Get(id int64) (res Stream, err error) {
	var stream = Stream{}
	result := db.Connect().First(&stream, id)
	return stream, result.Error
}

func GetBySlugOrId(slugorid string) (res Stream, err error) {
	var stream = Stream{}
	result := db.Connect().Where("id = ? ", slugorid).Or("slug = ? ", slugorid).First(&stream)
	return stream, result.Error
}

func GetByUser(id uint) (res []*Stream, err error) {
	var streams = []*Stream{}
	result := db.Connect().Where("createdBy = ?", id).Find(&streams)
	return streams, result.Error
}

func Create(data Stream) (res Stream, err error) {
	result := db.Connect().Model(&Stream{}).Create(&data)
	return data, result.Error
}

func (data Stream) Update() (res Stream, err error) {
	result := db.Connect().Save(&data)
	return data, result.Error
}

func (data Stream) Delete() (res Stream, err error) {
	result := db.Connect().Delete(&data)
	return data, result.Error
}

func GetAll(params StreamParams) (res []Stream, err error){
	query := db.Connect()
	if params.Live {
		query = query.Where("start_date is not null").Where("end_date is null")
	}
	result := query.Find(&res)
	err = result.Error
	return
}
