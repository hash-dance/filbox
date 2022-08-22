package deal

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/models"
	"gitee.com/szxjyt/filbox-backend/modules/common/render"
	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/deal"
	"gitee.com/szxjyt/filbox-backend/modules/ipfs"
	"gitee.com/szxjyt/filbox-backend/types"
)

// Router handler for miner
func FileRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/upload", apicontext.Bind(handleUpload))
	r.Post("/deal", apicontext.Bind(makedeal, types.PostDeal{}))
	r.Get("/list", apicontext.Bind(fileList))
	// r.Route("/{miner}", func(r chi.Router) {
	// 	r.Get("/", apicontext.Bind(askMiner))
	// })
	return r
}

// upload logic
func handleUpload(ctx *apicontext.APIContext) {
	ctx.Req.ParseMultipartForm(32 << 20)
	file, handler, err := ctx.Req.FormFile("uploadfile")
	if err != nil {
		return
	}
	defer file.Close()
	cid, err := ipfs.Push(ctx, file)
	if err != nil {
		ctx.Error(render.ServerError, err, "ipfs")
		return
	}
	fileInfo := models.File{}
	if err := mysql.GetClient().Where("phone = ? AND filecid = ?", ctx.Principal.Phone, cid).Find(&fileInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			fileInfo = models.File{
				Phone:    ctx.Principal.Phone,
				Filename: handler.Filename,
				Filecid:  cid,
				Size:     handler.Size,
			}
			if err := models.CreateFile(mysql.GetClient(), &fileInfo); err != nil {
				ctx.Error(render.ServerError, err, "保存")
				return
			}
			goto OK
		}
		ctx.Error(render.ServerError, err, "ipfs")
		return
	}
OK:
	ctx.JSON(fileInfo)
}

func makedeal(ctx *apicontext.APIContext, postForm *types.PostDeal) {
	phone := ctx.Principal.Phone
	// epochs := abi.ChainEpoch(dur / (time.Duration(build.BlockDelaySecs) * time.Second))
	tx := mysql.GetClient().Begin()
	wallet, err := models.GetWallet(tx, phone)
	if err != nil {
		tx.Rollback()
		ctx.Error(render.ServerError, err, "获取钱包错误")
		return
	}
	for _, cid := range postForm.Files {
		for _, miner := range postForm.MinerConfig {
			for range make([]int, miner.Nums) {
				dealinfo := models.DealInfo{
					Phone:    phone,
					Filecid:  cid,
					Miner:    miner.Miner,
					Price:    miner.Price,
					Duration: postForm.Duration,
					Wallet:   wallet.Address,
					TotalPrice: "0",
				}
				if err := tx.Model(&models.DealInfo{}).Create(&dealinfo).Error; err != nil {
					tx.Rollback()
					ctx.Error(render.ServerError, err, "创建交易失败")
					return
				}
			}
		}
	}

	tx.Commit()
	ctx.JSON("success")
	return
}

func fileList(ctx *apicontext.APIContext) {
	filedeals, err := deal.ListFiles(ctx.QueryInfo, ctx.Principal)
	if err != nil {
		ctx.Error(render.ServerError, err, "list filedeals")
		return
	}
	ctx.JSONPagination(filedeals)
}
