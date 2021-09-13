package users

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"listes_back/src/utils"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
)

var avatarsDir string

func InitAvatarsDir(avatarsDirPath string) error {
	avatarsDir = avatarsDirPath
	stats, err := os.Stat(avatarsDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Creating avatars directory ...")
			return os.MkdirAll(avatarsDirPath, 0755)
		}
		return err
	}
	if !stats.IsDir() {
		return errors.New("already exists but is a not a directory")
	}
	return nil
}

func getAvatarPath(id uint64) string {
	return path.Join(avatarsDir, strconv.FormatUint(id, 10)) + ".png"
}

func updateAvatar(userId uint64, r *http.Request) (error, int) {
	_ = r.ParseForm()
	err := r.ParseMultipartForm(5 << 20) // For 5MB max file size
	if err != nil {
		return err, http.StatusInternalServerError
	}

	avatarFile, header, err := r.FormFile("avatar")
	if err != nil {
		return err, http.StatusInternalServerError
	}
	defer avatarFile.Close()

	if header.Size > 5000000 { // 5000000 = 5MB
		return errors.New("the image is too big"), http.StatusBadRequest
	}

	img, format, err := image.Decode(avatarFile)
	if err != nil {
		return err, http.StatusBadRequest
	}
	if format != "png" {
		return fmt.Errorf("invalid image type: %s", format), http.StatusBadRequest
	}

	avatarPath := getAvatarPath(userId)
	targetFile, err := os.Create(avatarPath)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	defer targetFile.Close()

	err = png.Encode(targetFile, img)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func printAvatar(w http.ResponseWriter, id uint64) {
	file, err := os.Open(getAvatarPath(id))
	if err != nil {
		utils.Prettier(w, "no avatar found", nil, http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)

	_, err = io.Copy(w, file)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, "failed to send avatar", nil, http.StatusInternalServerError)
		return
	}
}

func deleteAvatar(id uint64) error {
	err := os.Remove(getAvatarPath(id))
	if err != nil {
		if err == os.ErrNotExist {
			return errors.New("no avatar found")
		}
		if runtime.GOOS == "windows" {
			return errors.New("error due to windows")
		}
		return err
	}
	return nil
}
