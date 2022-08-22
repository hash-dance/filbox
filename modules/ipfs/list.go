package ipfs

import (
	"path"

	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/sirupsen/logrus"
)

func List(apiCtx *apicontext.APIContext, cid string, recursive bool) ([]*shell.LsLink, error) {
	add := apiCtx.Config.Ipfs.Address
	logrus.Debugf("ipfs api %s\n", add)
	sh := NewSH(apiCtx.Context, add)
	return doList(sh, ".", cid, recursive)

}

func doList(sh *shell.Shell, base, cid string, recursive bool) ([]*shell.LsLink, error) {
	allLink := make([]*shell.LsLink, 0)
	links, err := sh.List(cid)
	if err != nil {
		logrus.Errorf("sh.List: %s", err.Error())
		return nil, err
	}
	if links == nil {
		return nil, nil
	}
	for _, l := range links {
		logrus.Infof("%s %s %d %d\n", path.Join(base, l.Name), l.Hash, l.Size, l.Type)
		allLink = append(allLink, l)
		if recursive && l.Type == 1 {
			res, err := doList(sh, path.Join(base, l.Name), l.Hash, recursive)
			if err != nil {
				continue
			}
			allLink = append(allLink, res...)
		}
	}
	return allLink, nil
}
