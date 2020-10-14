package tasks

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sync"
	"time"

	log "pigeond/log"
)

const pigeondWorkDir = "/var/run/pigeond"

var (
	scriptInventoryFile = path.Join(pigeondWorkDir, "script_inventory.csv")
	scriptDir           = path.Join(pigeondWorkDir, "scripts")
	scriptInventory     = make(map[string]*script)
	doOnce              sync.Once
)

type script struct {
	Name       string `json:"name"`
	CreateTime string `json:"create_time"`
	FileMD5    string `json:"-"`
	FileName   string `json:"-"`
}

type scriptList struct {
	Scripts []*script `json:"scripts"`
}

func getFileName(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func getFileMD5(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", scritpTaskError(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", scritpTaskError(err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Fix can't move file between different drive
func uploadFile(sourcePath, destPath string) error {

	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return uploadFileError(err.Error())
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		return uploadFileError(err.Error())
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return uploadFileError(err.Error())
	}
	return nil
}

func (s *script) addToInventory() error {
	// Add script into inventory file

	f, err := os.OpenFile(scriptInventoryFile, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	csvW := csv.NewWriter(f)
	err = csvW.Write([]string{s.Name, s.CreateTime, s.FileMD5, s.FileName})
	if err != nil {
		return err
	}
	csvW.Flush()
	return nil
}

func (s *script) removeFromInventory() error {

	// Get script inventory
	src, err := ioutil.ReadFile(scriptInventoryFile)
	if err != nil {
		return err
	}

	// Find script data and replace it
	r, err := regexp.Compile("," + s.FileName)
	if err != nil {
		return err
	}
	dst := r.ReplaceAll(src, []byte(""))

	// Write to file
	err = ioutil.WriteFile(scriptInventoryFile, dst, 0666)
	if err != nil {
		return err
	}
	return nil
}

func getInventoryData(file string) ([][]string, error) {
	// Get script inventory from file

	var lines = [][]string{}
	f, err := os.OpenFile(scriptInventoryFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return lines, err
	}
	defer f.Close()
	csvR := csv.NewReader(f)

	lines, err = csvR.ReadAll()
	if err != nil {
		return lines, err
	}

	return lines, nil
}

func loadInventory() {

	lines, err := getInventoryData(scriptInventoryFile)
	if err != nil {
		panic(err)
	}
	for _, line := range lines {
		s := script{}
		s.Name = line[0]
		s.CreateTime = line[1]
		s.FileMD5 = line[2]
		s.FileName = line[3]
		scriptInventory[s.FileName] = &s
	}

	log.Log.Debug("Load script inventory finished")
}

func init() {

	// Prepare pigeond work directory
	if _, err := os.Stat(scriptDir); os.IsNotExist(err) {
		err := os.MkdirAll(scriptDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	// Load script data into script inventory, only once
	doOnce.Do(loadInventory)
}

func listScripts(rstChan, errChan chan string) {

	sl := scriptList{}
	for _, v := range scriptInventory {
		sl.Scripts = append(sl.Scripts, v)
	}
	slByte, err := json.Marshal(sl)
	if err != nil {
		errChan <- err.Error()
		return
	}

	rstChan <- string(slByte)
}

func addScript(rstChan, errChan chan string, name, file string) {

	// Add script tar file into scripts dir
	_, err := os.Stat(file)
	if err != nil {
		errChan <- err.Error()
		return
	}

	// Check hash
	fileName := getFileName(name)
	if _, exist := scriptInventory[fileName]; exist {
		errChan <- "Script name not unique"
		return
	}
	fileMD5, err := getFileMD5(file)
	if err != nil {
		errChan <- err.Error()
		return
	}
	filePath := path.Join(scriptDir, fileName+".tar")
	err = uploadFile(file, filePath)
	if err != nil {
		errChan <- err.Error()
		return
	}
	s := &script{
		Name:       name,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		FileMD5:    fileMD5,
		FileName:   fileName,
	}

	// Add to script inventory
	err = s.addToInventory()
	if err != nil {
		// Remove file
		_ = os.Remove(filePath)
		errChan <- err.Error()
		return
	}
	scriptInventory[fileName] = s
	rstChan <- "Add script succeed"
}
