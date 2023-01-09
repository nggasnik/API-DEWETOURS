package middleware

import (
	"context"
	"encoding/json"
	"erlangga/dto"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// membuat middleware untuk menghandle upload file
func UpdateUserImage(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handling dan parsing data dari form data yang ada data file nya. Argumen 1024 pada method tersebut adalah maxMemory
		if err := r.ParseMultipartForm(1024); err != nil {
			panic(err.Error())
		}

		// mengambil id user dari url
		// id := mux.Vars(r)["id_user"]
		claims := r.Context().Value("userInfo").(jwt.MapClaims)
		id := strconv.Itoa(int(claims["id"].(float64)))

		// mengambil file dari form
		uploadedFile, handler, err := r.FormFile("image")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{
				Code:    http.StatusBadRequest,
				Message: "Please upload a JPG, JPEG or PNG image",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		defer uploadedFile.Close()

		// Apabila format file bukan .jpg, .jpeg atau .png, maka tampilkan error
		if filepath.Ext(handler.Filename) != ".jpg" && filepath.Ext(handler.Filename) != ".jpeg" && filepath.Ext(handler.Filename) != ".png" {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{
				Code:    http.StatusBadRequest,
				Message: "The provided file format is not allowed. Please upload a JPG, JPEG or PNG image",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		// mengambil direktori aktif
		dir, err := os.Getwd()
		if err != nil {
			panic(err.Error())
		}

		// menyusun string untuk digunakan sebagai nama file gambar
		random := strconv.FormatInt(time.Now().UnixNano(), 10)
		filenameStr := id + "-" + random

		// memberi nama pada file gambar
		filename := fmt.Sprintf("%s%s", filenameStr, filepath.Ext(handler.Filename))
		// fmt.Println(filename)

		// menentukan lokasi file
		fileLocation := filepath.Join(dir, "uploads/img", filename)
		// fmt.Println(fileLocation)

		// membuat file baru yang menjadi tempat untuk menampung hasil salinan file upload
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err.Error())
		}
		defer targetFile.Close()

		// Menyalin file hasil upload, ke file baru yang menjadi target
		if _, err := io.Copy(targetFile, uploadedFile); err != nil {
			panic(err.Error())
		}

		// membuat sebuah var untuk digunakan sebagai key pada context (untuk mengatasi warning should not use built-in type string as key)
		// var UploadFileID models.ContextKey = "FileName"

		// membuat sebuah context baru dengan menyisipkan value di dalamnya, valuenya adalah filename dan keynya adalah "NameOfUploadedFile"
		// ctx := context.WithValue(r.Context(), UploadFileID, filename)
		ctx := context.WithValue(r.Context(), "imageUploaded", filename)

		// mengirim nilai context ke object http.HandlerFunc yang menjadi parameter saat fungsi middleware ini dipanggil
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// membuat middleware untuk menghandle upload file
func UploadTripImage(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handling dan parsing data dari form data yang ada data file nya. Argumen 1024 pada method tersebut adalah maxMemory sebesar 1024byte, apabila file yang diupload lebih besar maka akan disimpan di file sementara
		if err := r.ParseMultipartForm(1024); err != nil {
			panic(err.Error())
		}

		// mengambil nama kerusakan dari req body
		titleTrip := r.FormValue("title")
		// fmt.Println(titleTrip)

		var arrImages []string

		files := r.MultipartForm.File["images"]
		for _, f := range files {
			// mengambil file dari form
			file, err := f.Open()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				response := dto.ErrorResult{
					Code:    http.StatusBadRequest,
					Message: "Please upload a JPG, JPEG or PNG image",
				}
				json.NewEncoder(w).Encode(response)
				return
			}
			defer file.Close()

			// fmt.Println(f.Filename)

			// Apabila format file bukan .jpg, .jpeg atau .png, maka tampilkan error
			if filepath.Ext(f.Filename) != ".jpg" && filepath.Ext(f.Filename) != ".jpeg" && filepath.Ext(f.Filename) != ".png" {
				w.WriteHeader(http.StatusBadRequest)
				response := dto.ErrorResult{
					Code:    http.StatusBadRequest,
					Message: "The provided file format is not allowed. Please upload a JPG, JPEG or PNG image",
				}
				json.NewEncoder(w).Encode(response)
				return
			}

			// mengambil direktori aktif
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
			}

			// menyusun string untuk digunakan sebagai nama file gambar
			filenameStr := titleTrip                       // mengambil title trip
			filenameStr = strings.ToLower(filenameStr)     // mem-format title trip menjadi huruf kecil
			filenameArr := strings.Split(filenameStr, " ") // memisahkan tiap kata pada title trip
			filenameStr = strings.Join(filenameArr, "-")   // menggabungkan title trip yang tadi dipisah (spasi berganti menjadi tanda '-')

			// menyusun string untuk digunakan sebagai nama file gambar
			random := strconv.FormatInt(time.Now().UnixNano(), 10)
			filenameStr = filenameStr + "-" + random

			// memberi nama pada file gambar
			filename := fmt.Sprintf("%s%s", filenameStr, filepath.Ext(f.Filename))
			// fmt.Println(filename)

			// menentukan lokasi file
			fileLocation := filepath.Join(dir, "uploads/img", filename)
			// fmt.Println(fileLocation)

			// membuat file baru yang menjadi tempat untuk menampung hasil salinan file upload
			targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				panic(err.Error())
			}
			defer targetFile.Close()

			// Menyalin file hasil upload, ke file baru yang menjadi target
			if _, err := io.Copy(targetFile, file); err != nil {
				panic(err.Error())
			}

			arrImages = append(arrImages, fileLocation)
		}

		// membuat sebuah var untuk digunakan sebagai key pada context (untuk mengatasi warning should not use built-in type string as key)
		// var UploadFileID models.ContextKey = "FileName"
		ctx := context.WithValue(r.Context(), "arrImages", arrImages)

		// mengirim nilai context ke object http.HandlerFunc yang menjadi parameter saat fungsi middleware ini dipanggil
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// membuat middleware untuk menghandle upload file transaction attachment
func UploadTransactionImage(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handling dan parsing data dari form data yang ada data file nya. Argumen 1024 pada method tersebut adalah maxMemory
		if err := r.ParseMultipartForm(1024); err != nil {
			panic(err.Error())
		}

		// mengambil id user dari middleware UserAuth
		claims := r.Context().Value("userInfo").(jwt.MapClaims)
		id := strconv.Itoa(int(claims["id"].(float64)))

		tripId := r.FormValue("tripId")

		// mengambil file dari form
		uploadedFile, handler, err := r.FormFile("attachment")
		if err != nil {
			ctx := context.WithValue(r.Context(), "attachment", "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		defer uploadedFile.Close()

		// Apabila format file bukan .jpg, .jpeg atau .png, maka tampilkan error
		if filepath.Ext(handler.Filename) != ".jpg" && filepath.Ext(handler.Filename) != ".jpeg" && filepath.Ext(handler.Filename) != ".png" {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{
				Code:    http.StatusBadRequest,
				Message: "The provided file format is not allowed. Please upload a JPG, JPEG or PNG image",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		// mengambil direktori aktif
		dir, err := os.Getwd()
		if err != nil {
			panic(err.Error())
		}

		// menyusun string untuk digunakan sebagai nama file gambar
		random := strconv.FormatInt(time.Now().UnixNano(), 10)
		filenameStr := "trx-" + id + "-" + tripId + "-" + random

		// memberi nama pada file gambar
		filename := fmt.Sprintf("%s%s", filenameStr, filepath.Ext(handler.Filename))
		// fmt.Println(filename)

		// menentukan lokasi file
		fileLocation := filepath.Join(dir, "uploads/img", filename)
		// fmt.Println(fileLocation)

		// membuat file baru yang menjadi tempat untuk menampung hasil salinan file upload
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err.Error())
		}
		defer targetFile.Close()

		// Menyalin file hasil upload, ke file baru yang menjadi target
		if _, err := io.Copy(targetFile, uploadedFile); err != nil {
			panic(err.Error())
		}

		// membuat sebuah var untuk digunakan sebagai key pada context (untuk mengatasi warning should not use built-in type string as key)
		// var UploadFileID models.ContextKey = "FileName"

		// membuat sebuah context baru dengan menyisipkan value di dalamnya, valuenya adalah filename dan keynya adalah "NameOfUploadedFile"
		// ctx := context.WithValue(r.Context(), UploadFileID, filename)
		ctx := context.WithValue(r.Context(), "attachment", filename)

		// mengirim nilai context ke object http.HandlerFunc yang menjadi parameter saat fungsi middleware ini dipanggil
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
