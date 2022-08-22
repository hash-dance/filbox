package ipfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"github.com/sirupsen/logrus"
)

func Push(apiCtx *apicontext.APIContext, r io.Reader) (string, error) {
	add := apiCtx.Config.Ipfs.Address
	logrus.Debugf("ipfs api %s\n", add)
	sh := NewSH(apiCtx.Context, add)
	return sh.Add(r)
}

type fileInfo struct {
	path string
	info os.FileInfo
}

func walkDirs(dir string) []*fileInfo {
	files := make([]*fileInfo, 0)
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, &fileInfo{
					path: path,
					info: info,
				})
				// fmt.Println(path, info.Size())
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
	}
	return files
}

func write2log(fullname, hash string, size, ftype int64) {
	logrus.Infof("%s %s %d %d\n", fullname, hash, size, ftype)
}

// func insertOne(col *mongo.Collection, d interface{}) (*mongo.InsertOneResult, error) {
// 	return col.InsertOne(context.Background(), d)
// }
