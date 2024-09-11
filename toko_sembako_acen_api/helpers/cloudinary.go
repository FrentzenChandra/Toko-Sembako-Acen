package helpers

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/spf13/viper"
)

func UploadToCloudinary(file *multipart.FileHeader) (string, error) {

	// Credentials
	cld, err := cloudinary.NewFromParams(viper.GetString("CLOUDINARY_CLOUD_NAME"), viper.GetString("CLOUDINARY_API_KEY"), viper.GetString("CLOUDINARY_API_SECRET"))

	if err != nil {
		log.Println("Error Uploading Image : " + err.Error())
		return "", err
	}

	cloudianryPathFolder := viper.GetString("CLOUDINARY_UPLOAD_FOLDER")

	// Upload the image on the cloud
	var ctx = context.Background()
	uploadResponse, err := cld.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{
			Folder:         cloudianryPathFolder,
			UniqueFilename: api.Bool(true),
			Overwrite:      api.Bool(false),
		},
	)

	if err != nil {
		log.Println("Error Uploading Image : " + err.Error())
		return "", err
	}

	// Return the image url
	return uploadResponse.SecureURL, nil
}

func DeleteAssetCloudinary(pictureUrl string) error {

	hasil := strings.Split(pictureUrl, "/")
	hasil = strings.Split(hasil[9], ".")

	cloudianryPathFolder := viper.GetString("CLOUDINARY_UPLOAD_FOLDER")
	fileUrl := cloudianryPathFolder + "/" + hasil[0]

	cld, err := cloudinary.NewFromParams(viper.GetString("CLOUDINARY_CLOUD_NAME"), viper.GetString("CLOUDINARY_API_KEY"), viper.GetString("CLOUDINARY_API_SECRET"))

	if err != nil {
		log.Println("Error Uploading Image : " + err.Error())
		return err
	}

	resp, err := cld.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: fileUrl})

	resultResp := resp.Result

	if resultResp != "ok" {
		err = errors.New("Delete Failed File : " + resultResp)
	}

	// if err != nil {
	// 	log.Println("ErrorDelete Image : " + err.Error())
	// 	return errors.New("Cloudinary Error Delete : " + resultResp)
	// }

	return nil

}
