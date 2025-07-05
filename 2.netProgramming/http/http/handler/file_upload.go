package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, req *http.Request) {
	// 上传单个文件
	// file, header, err := req.FormFile("files")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// err = saveFile(file, header.Filename)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// res := map[string]string{
	// 	"msg":       "上传文件成功",
	// 	"file_name": header.Filename,
	// }
	// bytes, _ := json.Marshal(res)
	// fmt.Fprintf(w, string(bytes))
	// return
	// 多个文件上传
	req.ParseMultipartForm(32 << 20)
	files := req.MultipartForm.File["files"]
	list := []string{}
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = saveFile(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		list = append(list, fileHeader.Filename)
	}
	res := map[string]interface{}{
		"msg":   "上传成功",
		"files": list,
	}
	bytes, _ := json.Marshal(res)
	fmt.Fprintf(w, string(bytes))
	return
}

func saveFile(file multipart.File, fileName string) error {
	filePath := "upload/" + fileName
	targetFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer targetFile.Close()
	_, err = io.Copy(targetFile, file)
	return err
}
