package helpers

import (
	"context"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2"
)

func UploadToCloudinary(file *multipart.FileHeader) (string, error) {
	// Remove from local
	defer func() {
		os.Remove("assets/uploads/" + file.Filename)
	}()

	cloudinary_url := "cloudinary://<api_key>:<api_secret>@<cloud_name>"
	cld, err := cloudinary.NewFromURL(cloudinary_url)

	// Upload the image on the cloud
	var ctx = context.Background()
	resp, err := cld.Upload.Upload(ctx, "assets/uploads/"+file.Filename, uploader.UploadParams{PublicID: "my_avatar" + "-" + file.Filename + "-" + GenerateUid()})

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// Return the image url
	return resp.SecureURL, nil
}
