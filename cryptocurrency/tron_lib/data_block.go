package tron_lib

type BlockData struct {
	BlockId      string             `json:"blockID"`
	BlockHeader  *BlockHeaderData   `json:"block_header"`
	Transactions []*TransactionData `json:"transactions"`
}

type BlockHeaderData struct {
	RawData          *BlockHeaderRawData `json:"raw_data"`
	WitnessSignature string              `json:"witness_signature"`
}

type BlockHeaderRawData struct {
	Number         int64  `json:"number"`
	TxTrieRoot     string `json:"txTrieRoot"`
	WitnessAddress string `json:"witness_address"`
	ParentHash     string `json:"parentHash"`
	Version        int64  `json:"version"`
	Timestamp      int64  `json:"timestamp"`
}

type TransactionData struct {
	Ret        []map[string]interface{} `json:"ret"`
	Signature  []string                 `json:"signature"`
	TxId       string                   `json:"txID"`
	RawDataHex string                   `json:"raw_data_hex"`
	RawData    *TransactionRawData      `json:"raw_data"`
}

type TransactionRawData struct {
	Contract      []*TransactionContractData `json:"contract"`
	RefBlockBytes string                     `json:"ref_block_bytes"`
	RefBlockHash  string                     `json:"ref_block_hash"`
	Expiration    int64                      `json:"expiration"`
	FeeLimit      int64                      `json:"fee_limit"`
	Timestamp     int64                      `json:"timestamp"`
}

type TransactionContractData struct {
	Parameter *ContractParameterData `json:"parameter"`
	Type      string                 `json:"type"`
}

type ContractParameterData struct {
	Value   *ContractParameterValueData `json:"value"`
	TypeUrl string                      `json:"type_url"`
}

type ContractParameterValueData struct {
	Amount          int64  `json:"amount"`
	AssetName       string `json:"asset_name"`
	Data            string `json:"data"`
	OwnerAddress    string `json:"owner_address"`
	ToAddress       string `json:"to_address"`
	ContractAddress string `json:"contract_address"`
	CallValue       int64  `json:"call_value"`
}
