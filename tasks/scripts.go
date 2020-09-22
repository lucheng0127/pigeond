package tasks

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path"

	log "pigeond/log"
)

const pigeondWorkDir = "/var/run/pigeond"

var scriptInventoryFile = path.Join(pigeondWorkDir, "script_inventory.csv")
var scriptDir = path.Join(pigeondWorkDir, "scripts")

var scriptInventory = make(map[string]*script)
var inventoryInited = false

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

func getInventoryData(file string) ([][]string, error) {
	var lines = [][]string{}
	f, err := os.OpenFile(scriptInventoryFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return lines, err
	}
	var fr io.Reader
	fr = f
	csvR := csv.NewReader(fr)

	lines, err = csvR.ReadAll()
	if err != nil {
		return lines, err
	}

	return lines, nil
}

func loadInventory(file string) error {

	if inventoryInited == true {
		// Only load once
		return nil
	}

	lines, err := getInventoryData(file)
	if err != nil {
		return err
	}
	for _, line := range lines {
		s := script{}
		s.Name = line[0]
		s.CreateTime = line[1]
		s.FileMD5 = line[2]
		s.FileName = line[3]
		scriptInventory[s.FileName] = &s
	}

	inventoryInited = true
	log.Log.Debug("Load script inventory finished")
	return nil
}

func init() {

	// Prepare pigeond work directory
	if _, err := os.Stat(scriptDir); os.IsNotExist(err) {
		err := os.MkdirAll(scriptDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	// Load script data into script inventory
	err := loadInventory(scriptInventoryFile)
	if err != nil {
		panic(err.Error())
	}
}

func listScripts(rstChan, errChan chan string) {

	sl := scriptList{}
	for _, v := range scriptInventory {
		sl.Scripts = append(sl.Scripts, v)
	}
	slByte, err := json.Marshal(sl)
	if err != nil {
		errChan <- err.Error()
	}

	rstChan <- string(slByte)
}
