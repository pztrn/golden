package action

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	cmn "github.com/vit1251/golden/pkg/common"
	"github.com/vit1251/golden/pkg/tracker"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

type FileAreaUploadCompleteAction struct {
	Action
}

func NewFileAreaUploadCompleteAction() *FileAreaUploadCompleteAction {
	return new(FileAreaUploadCompleteAction)
}

func (self FileAreaUploadCompleteAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fileManager := self.restoreFileManager()
	configManager := self.restoreConfigManager()

	outb, _ := configManager.Get("main", "Outbound")
	passwd, _ := configManager.Get("main", "Password")
	from, _ := configManager.Get("main", "Address")

	//
	vars := mux.Vars(r)
	echoTag := vars["echoname"]
	log.Printf("echoTag = %v", echoTag)

	/* Get file area */
	area, err1 := fileManager.GetAreaByName(echoTag)
	if err1 != nil {
		panic(err1)
	}
	log.Printf("area = %+v", area)

	/* ... */
	var maxMemory int64 = 128 * 1024 * 1024
	r.ParseMultipartForm(maxMemory)

	// in your case file would be fileupload
	stream, header, err2 := r.FormFile("file")
	if err2 != nil {
		panic(err2)
	}
	defer stream.Close()

	/* Description */
	desc := r.PostForm.Get("desc")

	//
	log.Printf("FileAreaUploadCompleteAction: filename = %+v", header.Filename)

	// Copy the file data to my buffer

	tmpFile := path.Join(outb, header.Filename)
	writeStream, err3 := os.Create(tmpFile)
	if err3 != nil {
		panic(err3)
	}
	cacheWriter := bufio.NewWriter(writeStream)
	defer func () {
		cacheWriter.Flush()
		writeStream.Close()
	}()

	size, err4 := io.Copy(cacheWriter, stream)
	if err4 != nil {
		panic(err4)
	}

	/* Create TIC description */
	ticBuilder := tracker.NewTicBuilder()

	ticBuilder.SetArea(area.GetName())
	ticBuilder.SetOrigin(from)
	ticBuilder.SetFrom(from)
	ticBuilder.SetFile(header.Filename)
	ticBuilder.SetDesc(desc)
	ticBuilder.SetSize(size)
	ticBuilder.SetPw(passwd)

	/* Save TIC on disk */
	newName := cmn.MakeTickName()
	newPath := path.Join(outb, newName)

	newContent := ticBuilder.Build()
	writer, err5 := os.Create(newPath)
	if err5 != nil {
		panic(err5)
	}
	cacheWriter2 := bufio.NewWriter(writer)
	defer func() {
		cacheWriter2.Flush()
		writer.Close()
	}()
	cacheWriter2.WriteString(newContent)

	/* Redirect */
	newLocation := fmt.Sprintf("/file/%s", area.GetName())
	http.Redirect(w, r, newLocation, 303)

}
