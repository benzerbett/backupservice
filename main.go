package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)
func main() {
	ctx := context.Background()
	// Load environment variables
	err := godotenv.Load()
	if err != nil{
		log.Fatalln(err)
	}
	minioUrl := os.Getenv("MINIO_ENDPOINT")
	backupsDir := os.Getenv("BACKUP_FILES_DIR")
	minioKey := os.Getenv("MINIO_ACCESSKEY")
	minioSecret := os.Getenv("MINIO_SECRET")
	minioSSL := false
	fmt.Println("URL: ",minioUrl)

	// Set up minio
	endpoint := minioUrl
	accessKeyID := minioKey
	secretAccessKey := minioSecret
	useSSL := minioSSL

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	//log.Printf("%#v\n", minioClient)

	// list buckets
	buckets, err := minioClient.ListBuckets(context.Background())
	if err != nil{
		log.Fatalln(err)
	}

	// get files
	dirFs, err := ioutil.ReadDir(backupsDir)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("\n\n")

	// loop through files
	for i, f := range dirFs {
		// for each file, find a bucket with a corresponding name
		fileName := strings.ToLower( strings.Split(f.Name(), ".sql")[0] )
		fileBucketIndex := slices.IndexFunc(buckets, func(b minio.BucketInfo) bool { return fileName == strings.ToLower(strings.Split(b.Name, "-backups")[0]) })
		if fileBucketIndex != -1 {
			fileBucket := buckets[fileBucketIndex]
			fmt.Println( fmt.Sprint(i+1)+") File: "+fmt.Sprint(f.Name())+", Bucket: "+fileBucket.Name)
			fopen, err := os.Open(backupsDir+"/"+f.Name())
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			if err != nil {
				log.Fatalln("Could not open file "+f.Name()+". ", err)
			}
			fileType, err := getFileContentType(fopen)
			if err != nil {
				log.Fatalln("Could not read file type.", err)
			}
			info, err := minioClient.FPutObject(ctx, fileBucket.Name, f.Name(), backupsDir+"/"+f.Name(), minio.PutObjectOptions{ContentType: fileType})
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("Successfully uploaded %s of size %d\n", f.Name(), info.Size)
		}else{
			fmt.Println( "File: "+fmt.Sprint(f.Name()) + ", Bucket not found")
		}
	}

}


func getFileContentType(ouput *os.File) (string, error) {
	buf := make([]byte, 512)
	_, err := ouput.Read(buf)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buf)
	return contentType, nil
}
