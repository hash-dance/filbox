package types

import (
	"time"
)

type PostDeal struct {
	Files       []string       `json:"files" validate:"gt=0,required"`
	MinerConfig []*MinerConfig `json:"minerConfig"  validate:"required,gt=0,dive"`
	Duration    int            `json:"duration" validate:"required,min=180,max=540"`
}

type MinerConfig struct {
	Miner string `json:"miner" validate:"required"`            // 矿工号
	Price string `json:"price" validate:"required"`            // 矿工价格
	Nums  int    `json:"nums" validate:"required,min=1,max=3"` // 存储份数
}

type DealOnline struct {
	Proposalcid struct {
		NAMING_FAILED string `json:"/"`
	} `json:"ProposalCid"`
	State      int    `json:"State"`
	Message    string `json:"Message"`
	Dealstages struct {
		Stages []struct {
			Name             string    `json:"Name"`
			Description      string    `json:"Description"`
			Expectedduration string    `json:"ExpectedDuration"`
			Createdtime      time.Time `json:"CreatedTime"`
			Updatedtime      time.Time `json:"UpdatedTime"`
			Logs             []struct {
				Log         string    `json:"Log"`
				Updatedtime time.Time `json:"UpdatedTime"`
			} `json:"Logs"`
		} `json:"Stages"`
	} `json:"DealStages"`
	Provider string `json:"Provider"`
	Dataref  struct {
		Transfertype string `json:"TransferType"`
		Root         struct {
			NAMING_FAILED string `json:"/"`
		} `json:"Root"`
		Piececid     interface{} `json:"PieceCid"`
		Piecesize    int         `json:"PieceSize"`
		Rawblocksize int         `json:"RawBlockSize"`
	} `json:"DataRef"`
	Piececid struct {
		NAMING_FAILED string `json:"/"`
	} `json:"PieceCID"`
	Size              int64     `json:"Size"`
	Priceperepoch     string    `json:"PricePerEpoch"`
	Duration          int       `json:"Duration"`
	Dealid            int       `json:"DealID"`
	Creationtime      time.Time `json:"CreationTime"`
	Verified          bool      `json:"Verified"`
	Transferchannelid struct {
		Initiator string `json:"Initiator"`
		Responder string `json:"Responder"`
		ID        int64  `json:"ID"`
	} `json:"TransferChannelID"`
	Datatransfer struct {
		Transferid int64 `json:"TransferID"`
		Status     int   `json:"Status"`
		Basecid    struct {
			NAMING_FAILED string `json:"/"`
		} `json:"BaseCID"`
		Isinitiator bool   `json:"IsInitiator"`
		Issender    bool   `json:"IsSender"`
		Voucher     string `json:"Voucher"`
		Message     string `json:"Message"`
		Otherpeer   string `json:"OtherPeer"`
		Transferred int64  `json:"Transferred"`
		Stages      struct {
			Stages []struct {
				Name        string      `json:"Name"`
				Description string      `json:"Description"`
				Createdtime time.Time   `json:"CreatedTime"`
				Updatedtime time.Time   `json:"UpdatedTime"`
				Logs        interface{} `json:"Logs"`
			} `json:"Stages"`
		} `json:"Stages"`
	} `json:"DataTransfer"`
}
